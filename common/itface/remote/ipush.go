package remote

type RpcPushService struct {
	Push_pushVendorMessage func(string, int, int, int) bool //按买家id获取
}

type PushService interface {
	PushMessage(string, int, int, int)(bool,error)
}
