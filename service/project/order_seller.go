package project

import (
	"order-backend/common"
	Config "order-backend/config/project"
	"order-backend/model"
	"order-backend/model/Form"
	"order-backend/repositorie"
	Rpc "order-backend/rpcserver/project"
	"strconv"
	"time"
)

// 卖家接单
func (o *ProjectOrder) AcceptOrder(params Form.AcceptOrderData) (model.Order, []model.OrderDetail, error) {
	orderStatus := Config.CCOrderStatus_NEED_ACCEPT
	orderRepository := repositorie.NewDefaultOrderRepositories(params.OrderId)
	var orderDetailList []model.OrderDetail
	orderData, errOrder := orderRepository.OrderGetOneByOid(params.OrderId)
	if errOrder != nil {
		return orderData, orderDetailList, errOrder
	}

	var orderIndexData model.OrderIndex
	orderColumns := []string{"OrderStatus", "UpdatedAt"}

	now := int(time.Now().Unix())
	if params.NewOrderType == Config.ORDER_ACCEPT_NEW ||
		params.NewOrderType == Config.ORDER_ACCEPT_PREPARE {
		//新订单和预订单接单
		if params.NewOrderType == Config.ORDER_ACCEPT_PREPARE {
			orderStatus = Config.CCOrderStatus_WAIT_PREPARE_GOODS
		}
		if orderData.OrderStatus != orderStatus {
			return orderData, orderDetailList, common.NewError(common.ErrorOrderStatusError)
		}

		orderData.OrderStatus = Config.CCOrderStatus_PREPARE_GOODS
		orderData.UpdatedAt = now
		orderIndexData.OrderStatus = orderData.OrderStatus
		orderIndexData.UpdatedAt = orderData.UpdatedAt

		where := map[string]interface{}{
			"order_id":     params.OrderId,
			"seller_id":    params.SellerId,
			"order_status": orderStatus,
		}
		errUp := orderRepository.UpdateOrder(where, orderData, orderColumns, orderIndexData, orderColumns)
		if errUp != nil {
			return orderData, orderDetailList, errUp
		}
		if params.IsBackData == 1 {
			orderDetailS, errD := orderRepository.GetOrderDetailById(params.OrderId)
			if errD != nil {
				return orderData, orderDetailList, common.NewError(common.ErrorUpdateDB)
			}
			orderDetailList = orderDetailS
		}
	} else if params.NewOrderType == Config.ORDER_ACCEPT_ADVANCE {
		//待取餐
		if orderData.OrderStatus != orderStatus {
			return orderData, orderDetailList, common.NewError(common.ErrorOrderStatusError)
		}
		if orderData.PrepareGoodsTime < int(now) {
			return orderData, orderDetailList, common.NewError(common.ErrorAdvanceOrderTimeError)
		}

		orderData.OrderStatus = Config.CCOrderStatus_WAIT_PREPARE_GOODS
		orderData.UpdatedAt = now
		orderIndexData.OrderStatus = orderData.OrderStatus
		orderIndexData.UpdatedAt = orderData.UpdatedAt

		where := map[string]interface{}{
			"order_id":     params.OrderId,
			"seller_id":    params.SellerId,
			"order_status": orderStatus,
		}
		errUp := orderRepository.UpdateOrder(where, orderData, orderColumns, orderIndexData, orderColumns)
		if errUp != nil {
			return orderData, orderDetailList, errUp
		}
	} else {
		return orderData, orderDetailList, common.NewError(common.ErrorFormData)
	}

	return orderData, orderDetailList, nil
}

// 卖家扫码
func (o *ProjectOrder) TakeGoodByCode(params Form.TakeGoodByCodeData) (model.Order, []model.OrderDetail, error) {
	var orderData model.Order
	var err error
	var orderDetailData []model.OrderDetail
	vendorInfo, errV := o.getVendorInfo(params.SellerId)
	if errV != nil {
		return orderData, orderDetailData, errV
	}

	where := make(map[string]interface{})
	takeGoodsCodeString := strconv.Itoa(params.TakeGoodsCode)
	if len([]byte(takeGoodsCodeString)) <= 3 {
		where["order_numeral"] = params.TakeGoodsCode
		if vendorInfo.SellerType == Config.SELLER_TYPE_TAKE_OUT {
			where["created_at"] = map[string]interface{}{
				"m": ">",
				"v": common.GetDayZeroTime(),
			}
		}
	} else {
		where["take_goods_code"] = params.TakeGoodsCode
	}

	repository := repositorie.NewDefaultOrderRepositories("")
	orderIndexData, errI := repository.GetOrderIndex(where, "order_id, order_status")
	if errI != nil {
		return orderData, orderDetailData, common.NewError(common.ErrorTakeGoodsCode)
	}

	orderId := orderIndexData[0].OrderId
	orderStatus := orderIndexData[0].OrderStatus

	orderData, err = repository.OrderGetOneByOid(orderId)
	if err != nil {
		return orderData, orderDetailData, err
	}
	orderDetailData, err = repository.GetOrderDetailById(orderId)
	if err != nil {
		return orderData, orderDetailData, err
	}

	switch orderStatus {
	case Config.CCOrderStatus_DONE:
		err = common.NewError(common.ErrorTakeGoodsCodeFailure3)
	case Config.CCOrderStatus_WAIT_TAKE_GOODS:
		err = common.NewError(common.ErrorTakeGoodsCodeFailure2)
	case Config.CCOrderStatus_PREPARE_GOODS:
		err = common.NewError(common.ErrorTakeGoodsCodeFailure1)
	}

	return orderData, orderDetailData, err
}

func (o *ProjectOrder) ConfirmTakeGoodData(params Form.ConfirmTakeGoodData) error {
	where := map[string]interface{}{
		"order_id":  params.OrderId,
		"seller_id": params.SellerId,
		"order_status": map[string]interface{}{
			"m": "in",
			"v": []int{
				Config.CCOrderStatus_WAIT_TAKE_GOODS,
				Config.CCOrderStatus_WAIT_PREPARE_GOODS,
			},
		},
	}

	repository := repositorie.NewDefaultOrderRepositories(params.OrderId)
	orderData, err := repository.OrderGetOneByOid(params.OrderId)
	if err != nil {
		return err
	}

	t := int(time.Now().Unix())
	orderData.OrderStatus 		= Config.CCOrderStatus_DONE
	orderData.UpdatedAt 		= t
	orderData.OrderCompleteTime = t

	var orderIndexData model.OrderIndex
	orderIndexData.OrderStatus 	= orderData.OrderStatus
	orderIndexData.UpdatedAt 	= orderData.UpdatedAt

	err = repositorie.NewDefaultOrderRepositories(params.OrderId).UpdateOrder(
		where,
		orderData,
		[]string{"OrderStatus", "UpdatedAt", "OrderCompleteTime"},
		orderIndexData,
		[]string{"OrderStatus", "UpdatedAt"})
	if err != nil {
		return err
	}

	err = new(Rpc.RpcSellerService).CalculateSellerScore(orderData.SellerId)
	if err != nil {
		return err
	}

	err = o.UpdateDataByOrderComplete(orderData.OrderId, orderData.SellerId, orderData.BuyerId)
	if err != nil {
		return err
	}

	return nil
}


