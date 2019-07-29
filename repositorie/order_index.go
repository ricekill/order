package repositorie

import (
	"order-backend/common"
	"order-backend/model"
	"fmt"
)

func (o *orderRepository) GetOrderIndex(where map[string]interface{}, fields string) ([]model.OrderIndex, error) {
	var orderIndexDataS []model.OrderIndex
	//使用find查询多条,使用get查询1条
	session := common.DB.Table(o.orderIndexTableName).Cols(fields)
	errSelect := GetWhere(session, where).Find(&orderIndexDataS)
	return orderIndexDataS, errSelect
}

func (o *orderRepository)InsertOrderIndex(orderIndex model.OrderIndex)(int64,error) {
	rId,err:=common.DB.Table(o.orderIndexTableName).InsertOne(orderIndex)
	return rId,err
}

func (o *orderRepository)DeleteOrderIndex(wheresql string)(int64,error){
	if wheresql == "" {
		return 0,nil
	}
	rid,err:=common.DB.Table(o.orderIndexTableName).Where(wheresql).Delete(model.OrderIndex{})
	return rid,err
}

func (o *orderRepository)UpdateOrderIndexByOid(orderIndex model.OrderIndex)(int64, error){
	var whereSql string
	whereSql=fmt.Sprintf("order_id=%s ",orderIndex.OrderId)
	affected,err := common.DB.Table(o.orderIndexTableName).Where(whereSql).Update(orderIndex)
	return affected,err
}