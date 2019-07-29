package repositorie

import (
	"order-backend/common"
	"order-backend/model"
)

func (o *orderRepository)GetOrderDetailById(orderId string)([]model.OrderDetail, error)  {
	var orderDetails []model.OrderDetail
	where := map[string]interface{}{"order_id" : orderId}
	session := common.DB.Table(o.orderDetailTableName)
	err := GetWhere(session, where).Find(&orderDetails)
	return orderDetails, err
}
