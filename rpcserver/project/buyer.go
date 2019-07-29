package project

import (
	"order-backend/common"
	"order-backend/common/itface/remote"
	"encoding/json"
	"github.com/hprose/hprose-golang/rpc"
	"order-backend/model/project"
)

type RpcBuyerService struct {
	   RpcClient rpc.Client
}

func (reus *RpcBuyerService) getBuyerRpcClient() remote.RpcBuyerService {
	var  buyerService remote.RpcBuyerService
	reus.RpcClient=rpc.NewClient(common.Config.RpcServer.Project.User)
	reus.RpcClient.UseService(&buyerService)
	return buyerService
}
func (reus *RpcBuyerService) GetBuyerInfo(bid int)(project.RpcBuyer,error) {
			//定义返回信息数据结构struct
		  var rpcBuyerInfo project.RpcBuyer
		  buyerInfo:=reus.getBuyerRpcClient().Buyer_getInfo(bid)
		  err := json.Unmarshal([]byte(buyerInfo), &rpcBuyerInfo)
		  if  rpcBuyerInfo.BuyerId <= 0 || err != nil {
			  errRpc :=common.NewError(common.ErrorBuyerInfo)
			  return rpcBuyerInfo,errRpc
		  }
		  defer reus.RpcClient.Close()
		  return rpcBuyerInfo,nil
}

func (reus *RpcBuyerService) DownOrderCount(userId int) error {
	buyerResult := reus.getBuyerRpcClient().Buyer_updateDownOrderCount(userId)
	var buyerData project.RpcReturn
	err := json.Unmarshal([]byte(buyerResult), &buyerData)
	if err != nil {
		return err
	}
	if !buyerData.Status {
		return common.NewError(common.ErrorUpdateDataFail)
	}
	defer reus.RpcClient.Close()
	return nil
}


