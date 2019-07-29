package service

import (
	"order-backend/common"
	"order-backend/common/itface/remote"
	"order-backend/rpcserver/project"
)

type LocalCartService struct {
	CartService remote.CartService
}

func NewCartService() remote.CartService{
	return new(LocalCartService).initService()
}
func (uo LocalCartService) initService() remote.CartService {
	if common.APP.Code == "project" {
		uo.CartService = new(project.RpcCartService)
	}
	return uo.CartService
}