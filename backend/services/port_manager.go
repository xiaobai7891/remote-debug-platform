package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"remote-debug-platform/config"
	"remote-debug-platform/models"
	"remote-debug-platform/protocols"
)

var (
	portManager     *PortManager
	portManagerOnce sync.Once
)

type PortManager struct {
	mu            sync.RWMutex
	usedPorts     map[int]bool           // 已使用的端口
	listeners     map[int]net.Listener   // 端口监听器
	connections   map[int]*PortConnections // 端口连接管理
	minPort       int
	maxPort       int
}

type PortConnections struct {
	mu       sync.RWMutex
	port     int
	portID   int64
	userID   int64
	protocol string
	devices  map[string]*DeviceConnection // deviceID -> connection
}

type DeviceConnection struct {
	DeviceID   string
	DeviceName string
	Conn       *websocket.Conn
	Protocol   protocols.Protocol
	ConnectedAt time.Time
}

func GetPortManager() *PortManager {
	portManagerOnce.Do(func() {
		cfg := config.Get()
		portManager = &PortManager{
			usedPorts:   make(map[int]bool),
			listeners:   make(map[int]net.Listener),
			connections: make(map[int]*PortConnections),
			minPort:     cfg.PortPool.Min,
			maxPort:     cfg.PortPool.Max,
		}
	})
	return portManager
}

// AllocatePort 分配一个空闲端口
func (pm *PortManager) AllocatePort() (int, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for port := pm.minPort; port <= pm.maxPort; port++ {
		if !pm.usedPorts[port] {
			// 检查端口是否真的可用
			listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				continue
			}
			listener.Close()

			// 检查数据库中是否已被占用
			used, _ := models.IsPortUsed(port)
			if used {
				continue
			}

			pm.usedPorts[port] = true
			return port, nil
		}
	}

	return 0, errors.New("no available port")
}

// ReleasePort 释放端口
func (pm *PortManager) ReleasePort(port int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.usedPorts, port)
}

// StartPortListener 启动端口监听
func (pm *PortManager) StartPortListener(port *models.Port) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 检查是否已经在监听
	if _, exists := pm.listeners[port.Port]; exists {
		return nil
	}

	// 创建 HTTP 服务器用于 WebSocket 升级
	mux := http.NewServeMux()

	// 创建连接管理器
	pc := &PortConnections{
		port:     port.Port,
		portID:   port.ID,
		userID:   port.UserID,
		protocol: port.Protocol,
		devices:  make(map[string]*DeviceConnection),
	}
	pm.connections[port.Port] = pc

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pm.handleConnection(w, r, port)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port.Port),
		Handler: mux,
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}

	pm.listeners[port.Port] = listener
	pm.usedPorts[port.Port] = true

	go func() {
		log.Printf("端口 %d 开始监听 (协议: %s)", port.Port, port.Protocol)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("端口 %d 监听错误: %v", port.Port, err)
		}
	}()

	return nil
}

// StopPortListener 停止端口监听
func (pm *PortManager) StopPortListener(port int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if listener, exists := pm.listeners[port]; exists {
		listener.Close()
		delete(pm.listeners, port)
		log.Printf("端口 %d 停止监听", port)
	}

	// 断开所有连接
	if pc, exists := pm.connections[port]; exists {
		pc.mu.Lock()
		for _, dc := range pc.devices {
			if dc.Conn != nil {
				dc.Conn.Close()
			}
		}
		pc.devices = make(map[string]*DeviceConnection)
		pc.mu.Unlock()
		delete(pm.connections, port)
	}

	// 设置该端口所有设备离线
	if portModel, err := models.GetPortByNumber(port); err == nil {
		models.SetPortDevicesOffline(portModel.ID)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (pm *PortManager) handleConnection(w http.ResponseWriter, r *http.Request, port *models.Port) {
	// 升级为 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		return
	}

	// 获取协议适配器
	protocol := protocols.GetProtocol(port.Protocol)

	// 等待认证消息
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Printf("读取认证消息失败: %v", err)
		conn.Close()
		return
	}
	conn.SetReadDeadline(time.Time{})

	// 解析认证消息
	authResult, err := protocol.Authenticate(message, port.Token)
	if err != nil {
		log.Printf("认证失败: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"auth_failed","message":"`+err.Error()+`"}`))
		conn.Close()
		return
	}

	// 发送认证成功
	conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"auth_success"}`))

	// 注册设备
	deviceInfo := &models.DeviceInfo{
		Brand:   authResult.Brand,
		Model:   authResult.Model,
		OS:      authResult.OS,
		Version: authResult.Version,
	}
	device, err := models.CreateOrUpdateDevice(port.UserID, port.ID, authResult.DeviceID, authResult.DeviceName, deviceInfo)
	if err != nil {
		log.Printf("注册设备失败: %v", err)
		conn.Close()
		return
	}

	log.Printf("设备连接: %s (%s) 端口: %d", device.DeviceName, device.DeviceID, port.Port)

	// 添加到连接管理
	pm.mu.RLock()
	pc, exists := pm.connections[port.Port]
	pm.mu.RUnlock()

	if !exists {
		conn.Close()
		return
	}

	dc := &DeviceConnection{
		DeviceID:    authResult.DeviceID,
		DeviceName:  authResult.DeviceName,
		Conn:        conn,
		Protocol:    protocol,
		ConnectedAt: time.Now(),
	}

	pc.mu.Lock()
	pc.devices[authResult.DeviceID] = dc
	pc.mu.Unlock()

	// 记录日志
	models.AddDeviceLog(device.ID, "info", "设备上线")

	// 处理消息
	pm.handleMessages(conn, port, device, protocol)

	// 断开连接处理
	pc.mu.Lock()
	delete(pc.devices, authResult.DeviceID)
	pc.mu.Unlock()

	device.SetOffline()
	models.AddDeviceLog(device.ID, "info", "设备下线")
	log.Printf("设备断开: %s (%s)", device.DeviceName, device.DeviceID)
}

func (pm *PortManager) handleMessages(conn *websocket.Conn, port *models.Port, device *models.Device, protocol protocols.Protocol) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// 解析消息
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		msgType, _ := msg["type"].(string)

		switch msgType {
		case "ping":
			// 心跳
			conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"pong"}`))
		case "log":
			// 日志
			level, _ := msg["level"].(string)
			content, _ := msg["message"].(string)
			if level == "" {
				level = "info"
			}
			models.AddDeviceLog(device.ID, level, content)
		case "result":
			// 执行结果
			result, _ := msg["data"].(string)
			models.AddDeviceLog(device.ID, "info", "执行结果: "+result)
		case "screenshot":
			// 截图结果
			data, _ := msg["data"].(string)
			models.AddDeviceLog(device.ID, "info", "截图完成: "+data[:min(len(data), 50)]+"...")
		default:
			// 其他消息由协议处理
			protocol.HandleMessage(msg)
		}
	}
}

// SendToDevice 向设备发送消息
func (pm *PortManager) SendToDevice(port int, deviceID string, message interface{}) error {
	pm.mu.RLock()
	pc, exists := pm.connections[port]
	pm.mu.RUnlock()

	if !exists {
		return errors.New("port not found")
	}

	pc.mu.RLock()
	dc, exists := pc.devices[deviceID]
	pc.mu.RUnlock()

	if !exists {
		return errors.New("device not found")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return dc.Conn.WriteMessage(websocket.TextMessage, data)
}

// BroadcastToPort 向端口所有设备广播消息
func (pm *PortManager) BroadcastToPort(port int, message interface{}) error {
	pm.mu.RLock()
	pc, exists := pm.connections[port]
	pm.mu.RUnlock()

	if !exists {
		return errors.New("port not found")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	pc.mu.RLock()
	defer pc.mu.RUnlock()

	for _, dc := range pc.devices {
		dc.Conn.WriteMessage(websocket.TextMessage, data)
	}

	return nil
}

// GetOnlineDevices 获取端口在线设备
func (pm *PortManager) GetOnlineDevices(port int) []string {
	pm.mu.RLock()
	pc, exists := pm.connections[port]
	pm.mu.RUnlock()

	if !exists {
		return nil
	}

	pc.mu.RLock()
	defer pc.mu.RUnlock()

	devices := make([]string, 0, len(pc.devices))
	for id := range pc.devices {
		devices = append(devices, id)
	}

	return devices
}

// InitFromDB 从数据库恢复活跃端口
func (pm *PortManager) InitFromDB() error {
	ports, err := models.GetActivePorts()
	if err != nil {
		return err
	}

	for _, port := range ports {
		// 检查是否过期
		if port.ExpireAt.Before(time.Now()) {
			port.Status = "expired"
			port.Update()
			models.SetPortDevicesOffline(port.ID)
			continue
		}

		if err := pm.StartPortListener(port); err != nil {
			log.Printf("启动端口 %d 失败: %v", port.Port, err)
		}
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
