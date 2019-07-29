package service

import (
	"order-backend/common"
	"order-backend/common/itface/remote"
	"order-backend/rpcserver/project"
)

type LocalBuyerService struct {
     BuyerService remote.BuyerService
}
func NewBuyerService() remote.BuyerService{
	return new(LocalBuyerService).InitService()
}
func (uo *LocalBuyerService) InitService() remote.BuyerService {
	if common.APP.Code == "project" {
		uo.BuyerService = new(project.RpcBuyerService)
	}
	return uo.BuyerService
}
