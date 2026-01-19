package protocols

import (
	"encoding/json"
	"errors"
)

// LazymanProtocol 懒人精灵协议适配器
type LazymanProtocol struct{}

type LazymanAuthMessage struct {
	Type   string            `json:"type"`
	Token  string            `json:"token"`
	Device LazymanDeviceInfo `json:"device"`
}

type LazymanDeviceInfo struct {
	UDID       string `json:"udid"`
	DeviceName string `json:"device_name"`
	Model      string `json:"model"`
	OSVersion  string `json:"os_version"`
	AppVersion string `json:"app_version"`
}

func (p *LazymanProtocol) Name() string {
	return "lazyman"
}

func (p *LazymanProtocol) Authenticate(message []byte, token string) (*AuthResult, error) {
	var msg LazymanAuthMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return nil, errors.New("消息格式错误")
	}

	// 检查消息类型
	if msg.Type != "register" && msg.Type != "auth" && msg.Type != "connect" {
		return nil, errors.New("无效的认证消息类型")
	}

	// 验证 token
	if msg.Token != token {
		return nil, errors.New("token 无效")
	}

	// 返回认证结果
	return &AuthResult{
		DeviceID:   msg.Device.UDID,
		DeviceName: msg.Device.DeviceName,
		Brand:      "Apple", // 懒人精灵主要用于 iOS
		Model:      msg.Device.Model,
		OS:         "iOS " + msg.Device.OSVersion,
		Version:    msg.Device.AppVersion,
	}, nil
}

func (p *LazymanProtocol) HandleMessage(message map[string]interface{}) {
	// 懒人精灵特定的消息处理
}
