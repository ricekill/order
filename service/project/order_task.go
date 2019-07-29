package project

import (
	"order-backend/common"
	Config "order-backend/config/project"
	Rpc "order-backend/rpcserver/project"
	"order-backend/model"
	"order-backend/repositorie"
	"fmt"
	"strings"
	"time"
)

type ProjectOrderTask struct {
	AppInfo model.App
}

func (ot *ProjectOrderTask) OrderCancelNotPay() error {
	now := time.Now().Unix()
	cancelWithNotPayOnTime := common.App.Project.SCHEDULE_CANCEL_TIME
	timeCondition := int(now) - cancelWithNotPayOnTime
	where := map[string]interface{}{
		"order_status" : Config.CCOrderStatus_UNPAID,
		"created_at" : map[string]interface{}{
			"m" : "<",
			"v" : timeCondition,
		}}
	orderIndexDatas, errS := repositorie.NewDefaultOrderRepositories("").GetOrderIndex(where, "order_id")
	if errS != nil {
		return errS
	}
	//修改订单状态操作
	if len(orderIndexDatas) > 0 {
		//逐个修改状态
		for _, orderIndexData := range orderIndexDatas {
			var orderModel model.Order
			orderModel.OrderStatus = Config.CCOrderStatus_CANCEL
			orderModel.UpdatedAt = int(now)
			orderModel.OrderId = orderIndexData.OrderId
			_, errOr := repositorie.NewDefaultOrderRepositories(orderIndexData.OrderId).OrderUpdate(orderModel)
			if errOr != nil {
				return errOr
			}
			var orderIndexModel model.OrderIndex
			orderIndexModel.OrderStatus = Config.CCOrderStatus_CANCEL
			orderIndexModel.OrderId = orderIndexData.OrderId
			orderIndexModel.UpdatedAt = int(now)
			_, errOir := repositorie.NewDefaultOrderRepositories("").UpdateOrderIndexByOid(orderIndexModel)
			if errOir != nil {
				return errOir
			}

			ot.OrderCancelReturnData(orderIndexData.OrderId)
			ot.OrderCancelReturnStock(orderIndexData.OrderId)
		}
	}
	return nil

}
func (ot *ProjectOrderTask) OrderCancelReturnData(orderId string) error {
	orderModel, err := repositorie.NewDefaultOrderRepositories(orderId).OrderGetOneByOid(orderId)
	if err != nil {
		return err
	}
	if orderModel.CouponUserCode != "" {
		//返还优惠券
		errorRC := new(Rpc.RpcCouponService).ReturnCouponCode(orderModel.BuyerId, orderModel.CouponUserCode)
		if errorRC != nil {
			return errorRC
		}
	}
	//fmt.Printf("==========%v==========",orderModel)
	return nil
}

//返还商品库存
func (ot *ProjectOrderTask) OrderCancelReturnStock(orderId string) error {
	detailList, errGet := repositorie.NewDefaultOrderRepositories(orderId).GetOrderDetailById(orderId)
	if errGet != nil {
		return errGet
	}
	if len(detailList) > 0 {
		fmt.Printf("========detailList=%v========", detailList)
		var foodStack map[int]int
		foodStack = make(map[int]int)
		var action = "add"
		var sellerId = 0
		for _, dItem := range detailList {
			if dItem.ProductId > 0 {
				sellerId = dItem.SellerId
				mm, ok := foodStack[dItem.ProductId]
				if !ok {
					foodStack[dItem.ProductId] = dItem.ProductPriceCount
				} else {
					foodStack[dItem.ProductId] = mm + dItem.ProductPriceCount
				}
			}
		}
		return new(Rpc.RpcProductsService).ModifyProductsStock(foodStack, action, sellerId)
		fmt.Printf("========detailList=%v========", foodStack)
	}
	return nil
}

func (ot *ProjectOrderTask) OrderCancel(orderId string, reason string) (int64, error) {
	var order model.Order
	order.OrderId = orderId
	return 1, nil

}
func (ot *ProjectOrderTask) OrderDetailByOid(orderId string) (model.Order, error) {
	return repositorie.NewDefaultOrderRepositories(orderId).OrderGetOneByOid(orderId)
}
func (ot *ProjectOrderTask) OrderCountById(id int, idType string, filter map[string]interface{}) (int64) {
	var filterW, idKey string
	var filterWArr []string
	if len(filter) > 0 {
		filterWArr = append(filterWArr, fmt.Sprintf("%s = %d", idKey, id))
		for fiK, fiV := range filter {
			if fiV != "" {
				filterWArr = append(filterWArr, fmt.Sprintf("%s = %s", fiK, fiV))
			}
		}
		filterW = strings.Join(filterWArr, " AND ")
	}

	if idType == "seller" {
		//idKey="seller_id"
		totalCount, _ := repositorie.NewDefaultOrderRepositories("").OrderSearchBySellerId(id, filterW)
		return totalCount
	} else {
		//idKey="buyer_id"
		totalCount, _ := repositorie.NewDefaultOrderRepositories("").OrderSearchByBuyerId(id, filterW)
		return totalCount
	}

}
