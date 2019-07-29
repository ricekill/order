package service

import (
	"order-backend/common"
	"order-backend/repositorie"
	"order-backend/service/project"
	"fmt"
	"time"
)

func NewOrderService() common.OrderService {
	var os common.OrderService
	switch common.APP.Code {
		case "project":
			os = new(project.ProjectOrder)
			break
	}
	return os
}

//清理数据库order_index 7天前的数据
func ClearOrderIndexData() (int64, error) {
	orderIndexRepostorie := repositorie.NewDefaultOrderRepositories("")
	dayBefore := time.Now().Unix() - 7*86400
	whereSql := fmt.Sprintf("take_goods_time < %d", dayBefore)
	rid, err := orderIndexRepostorie.DeleteOrderIndex(whereSql)
	return rid, err
}
