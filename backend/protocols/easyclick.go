package protocols

import (
	"encoding/json"
	"errors"
)

// EasyClickProtocol EasyClick 协议适配器
type EasyClickProtocol struct{}

type EasyClickAuthMessage struct {
	Type   string              `json:"type"`
	Token  string              `json:"token"`
	Device EasyClickDeviceInfo `json:"device"`
}

type EasyClickDeviceInfo struct {
	DeviceID   string `json:"deviceId"`
	DeviceName string `json:"deviceName"`
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	OSVersion  string `json:"osVersion"`
	SDKVersion int    `json:"sdkVersion"`
}

func (p *EasyClickProtocol) Name() string {
	return "easyclick"
}

func (p *EasyClickProtocol) Authenticate(message []byte, token string) (*AuthResult, error) {
	var msg EasyClickAuthMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return nil, errors.New("消息格式错误")
	}

	// 检查消息类型
	if msg.Type != "register" && msg.Type != "auth" && msg.Type != "init" {
		return nil, errors.New("无效的认证消息类型")
	}

	// 验证 token
	if msg.Token != token {
		return nil, errors.New("token 无效")
	}

	// 返回认证结果
	return &AuthResult{
		DeviceID:   msg.Device.DeviceID,
		DeviceName: msg.Device.DeviceName,
		Brand:      msg.Device.Brand,
		Model:      msg.Device.Model,
		OS:         "Android " + msg.Device.OSVersion,
		Version:    "",
	}, nil
}

func (p *EasyClickProtocol) HandleMessage(message map[string]interface{}) {
	// EasyClick 特定的消息处理
}
