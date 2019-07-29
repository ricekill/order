package service

import (
	"order-backend/common"
	"order-backend/common/itface/remote"
	"order-backend/rpcserver/project"
)

type LocalSellerService struct {
	SellerService remote.SellerService
}
func NewVendorService() remote.SellerService{
	return new(LocalSellerService).initService()
}
func (uo LocalSellerService) initService() remote.SellerService {
	if common.APP.Code == "project" {
		uo.SellerService = new(project.RpcSellerService)
	}
	return uo.SellerService
}
