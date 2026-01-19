package protocols

// AuthResult 认证结果
type AuthResult struct {
	DeviceID   string
	DeviceName string
	Brand      string
	Model      string
	OS         string
	Version    string
}

// Protocol 协议接口
type Protocol interface {
	// Authenticate 认证
	Authenticate(message []byte, token string) (*AuthResult, error)
	// HandleMessage 处理消息
	HandleMessage(message map[string]interface{})
	// Name 协议名称
	Name() string
}

// GetProtocol 获取协议适配器
func GetProtocol(name string) Protocol {
	switch name {
	case "autojs":
		return &AutoJSProtocol{}
	case "lazyman":
		return &LazymanProtocol{}
	case "easyclick":
		return &EasyClickProtocol{}
	default:
		return &GenericProtocol{}
	}
}
