package common

import (
	"order-backend/model"
	"order-backend/model/Form"
)
//定义一个OrderService接口类,用到project,或者其它平台接入的时候实现
type OrderService interface {
	OrderCreate(form Form.OrderCreateData) (string, error)
	PaySuccess(form Form.PaySuccessData) (model.Order, error)
	//OrderCancel(string, string) (int64,  error)
	OrderList(int,string,PageParam,map[string]interface{},map[string]interface{}) ([]model.Order, *Pagination)
	OrderCountById(int,string,map[string]interface{}) int64
	AcceptOrder(form Form.AcceptOrderData) (model.Order, []model.OrderDetail, error)
	TakeGoodByCode(form Form.TakeGoodByCodeData) (model.Order, []model.OrderDetail, error)
	ConfirmTakeGoodData(form Form.ConfirmTakeGoodData) error
}
