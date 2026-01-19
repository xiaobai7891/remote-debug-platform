package protocols

import (
	"encoding/json"
	"errors"
)

// AutoJSProtocol Auto.js 协议适配器
type AutoJSProtocol struct{}

type AutoJSAuthMessage struct {
	Type   string          `json:"type"`
	Token  string          `json:"token"`
	Device AutoJSDeviceInfo `json:"device"`
}

type AutoJSDeviceInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Brand   string `json:"brand"`
	Model   string `json:"model"`
	OS      string `json:"os"`
	Version string `json:"version"`
	SDK     int    `json:"sdk"`
}

func (p *AutoJSProtocol) Name() string {
	return "autojs"
}

func (p *AutoJSProtocol) Authenticate(message []byte, token string) (*AuthResult, error) {
	var msg AutoJSAuthMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return nil, errors.New("消息格式错误")
	}

	// 检查消息类型
	if msg.Type != "register" && msg.Type != "auth" && msg.Type != "hello" {
		return nil, errors.New("无效的认证消息类型")
	}

	// 验证 token
	if msg.Token != token {
		return nil, errors.New("token 无效")
	}

	// 返回认证结果
	return &AuthResult{
		DeviceID:   msg.Device.ID,
		DeviceName: msg.Device.Name,
		Brand:      msg.Device.Brand,
		Model:      msg.Device.Model,
		OS:         msg.Device.OS,
		Version:    msg.Device.Version,
	}, nil
}

func (p *AutoJSProtocol) HandleMessage(message map[string]interface{}) {
	// Auto.js 特定的消息处理
}
