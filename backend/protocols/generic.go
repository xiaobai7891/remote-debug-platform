package protocols

import (
	"encoding/json"
	"errors"
)

// GenericProtocol 通用协议适配器
type GenericProtocol struct{}

type GenericAuthMessage struct {
	Type   string            `json:"type"`
	Token  string            `json:"token"`
	Device GenericDeviceInfo `json:"device"`
}

type GenericDeviceInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Brand   string `json:"brand"`
	Model   string `json:"model"`
	OS      string `json:"os"`
	Version string `json:"version"`
}

func (p *GenericProtocol) Name() string {
	return "generic"
}

func (p *GenericProtocol) Authenticate(message []byte, token string) (*AuthResult, error) {
	var msg GenericAuthMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return nil, errors.New("消息格式错误")
	}

	// 检查消息类型
	if msg.Type != "register" && msg.Type != "auth" && msg.Type != "connect" && msg.Type != "hello" {
		return nil, errors.New("无效的认证消息类型")
	}

	// 验证 token
	if msg.Token != token {
		return nil, errors.New("token 无效")
	}

	// 返回认证结果
	deviceID := msg.Device.ID
	if deviceID == "" {
		deviceID = "unknown-" + randomString(8)
	}

	deviceName := msg.Device.Name
	if deviceName == "" {
		deviceName = "Unknown Device"
	}

	return &AuthResult{
		DeviceID:   deviceID,
		DeviceName: deviceName,
		Brand:      msg.Device.Brand,
		Model:      msg.Device.Model,
		OS:         msg.Device.OS,
		Version:    msg.Device.Version,
	}, nil
}

func (p *GenericProtocol) HandleMessage(message map[string]interface{}) {
	// 通用消息处理
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
