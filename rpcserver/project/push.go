package project

import (
	"order-backend/common"
	"order-backend/common/itface/remote"
	"github.com/hprose/hprose-golang/rpc"
)

type RpcPushService struct {
	RpcClient rpc.Client
}

func (reus *RpcPushService) getPushRpcClient() remote.RpcPushService {
	var pushService remote.RpcPushService
	reus.RpcClient = rpc.NewClient(common.Config.RpcServer.Project.Push)
	reus.RpcClient.UseService(&pushService)
	return pushService
}
func (reus *RpcPushService) PushMessage(
	orderId string,
	sellerId,
	messageType,
	newOrderType int) bool {
	//定义返回信息数据结构struct
	pushResult := reus.getPushRpcClient().Push_pushVendorMessage(orderId, sellerId, messageType, newOrderType)
	if !pushResult {
		common.Log.Infoln("[PUSH]","push error")
	}
	defer reus.RpcClient.Close()
	return pushResult
}
