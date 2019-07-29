package api

import (
	"order-backend/common"
	"order-backend/model/Form"
	"order-backend/service"
	"encoding/json"
)

func AcceptOrder(params string) string {
	var acceptOrderData Form.AcceptOrderData
	err := json.Unmarshal([]byte(params), &acceptOrderData)
	if err != nil || acceptOrderData.OrderId == "" || acceptOrderData.SellerId <= 0 {
		common.Log.Infof("[func=AcceptOrder] params==========%v \r\n", params)
		return common.RenderRpcJson("", "params error", common.ErrorFormData, false)
	}

	orderData, orderDetailData, errAccept := service.NewOrderService().AcceptOrder(acceptOrderData)
	if errAccept != nil {
		return common.RenderRpcJSONError(errAccept)
	}

	orderDataL := map[string]interface{}{
		"order" : orderData,
		"order_detail" : orderDetailData,
	}
	return common.RenderRpcJson(orderDataL, "success", 0, true)
}

func TakeGoodByCode(params string) string {
	var takeGoodsCodeData Form.TakeGoodByCodeData
	err := json.Unmarshal([]byte(params), &takeGoodsCodeData)
	if err != nil || takeGoodsCodeData.TakeGoodsCode <= 0 || takeGoodsCodeData.SellerId <= 0 {
		common.Log.Infof("[func=TakeGoodByCode] params==========%v \r\n", params)
		return common.RenderRpcJson("", "params error", common.ErrorFormData, false)
	}

	orderData, orderDetailData, errAccept := service.NewOrderService().TakeGoodByCode(takeGoodsCodeData)
	if errAccept != nil {
		return common.RenderRpcJSONError(errAccept)
	}

	orderDataL := map[string]interface{}{
		"order" : orderData,
		"order_detail" : orderDetailData,
	}
	return common.RenderRpcJson(orderDataL, "success", 0, true)
}

func ConfirmTakeGood(params string) string {
	var ConfirmTakeGoodData Form.ConfirmTakeGoodData
	err := json.Unmarshal([]byte(params), &ConfirmTakeGoodData)
	if err != nil || ConfirmTakeGoodData.OrderId =="" || ConfirmTakeGoodData.SellerId <= 0 {
		common.Log.Infof("[func=ConfirmTakeGood] params==========%v \r\n", params)
		return common.RenderRpcJson("", "params error", common.ErrorFormData, false)
	}

	err = service.NewOrderService().ConfirmTakeGoodData(ConfirmTakeGoodData)
	if err != nil {
		return common.RenderRpcJSONError(err)
	}

	return common.RenderRpcJson("", "success", 0, true)
}