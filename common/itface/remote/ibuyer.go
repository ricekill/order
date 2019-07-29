package remote

import "order-backend/model/project"

type RpcBuyerService struct {
	Buyer_getInfo func(int) string //按买家id获取
	Buyer_updateDownOrderCount func(int) string //更新用户下单数量
}

type BuyerService interface {
	GetBuyerInfo(int)(project.RpcBuyer,error)
	DownOrderCount(int) error
}
