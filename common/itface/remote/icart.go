package remote

import "order-backend/model/project"

type RpcCartService struct {
	Cart_getCart func(int, int) string//按用户id,商家id获取购物车
	Cart_delCart func(int, int) string//按用户id,商家id删除购物车
}

type CartService interface {
	GetCartInfo(int,int)  (project.RpcCart,error) //按用户id,商家id获取购物车
	DeleteCart(int, int)  error
}