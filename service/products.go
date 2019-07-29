package service

import (
	"order-backend/common"
	"order-backend/common/itface/remote"
	"order-backend/rpcserver/project"
)

type LocalProductsService struct {
	ProductsService remote.ProductsService
}

func NewProductsService() remote.ProductsService{
	return new(LocalProductsService).initService()
}
func (uo LocalProductsService) initService() remote.ProductsService {
	if common.APP.Code == "project" {
		uo.ProductsService = new(project.RpcProductsService)
	}
	return uo.ProductsService
}
