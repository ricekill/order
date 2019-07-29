package api

import (
	"order-backend/common"
	"order-backend/model/Form"
	"order-backend/service"
	"encoding/json"
)

func OrderCreate(params string) string {
	var orderData Form.OrderCreateData
	err := json.Unmarshal([]byte(params), &orderData)
	if err != nil || orderData.BuyerId <= 0 || orderData.SellerId <= 0 || int(orderData.TakeGoodsTime) <= 0 {
		common.Log.Infof("[func=OrderCreate] params==========%v \r\n", params)
		return common.RenderRpcJson("", "params error", common.ErrorFormData, false)
	}

	orderId, errCreate := service.NewOrderService().OrderCreate(orderData)
	if errCreate != nil {
		return common.RenderRpcJSONError(errCreate)
	}

	orderResult := make(map[string]interface{})
	orderResult["order_id"] = orderId
	return common.RenderRpcJson(orderResult, "success", 0, true)
}

func PaySuccess(params string) string {
	var payData Form.PaySuccessData
	err := json.Unmarshal([]byte(params), &payData)
	if err != nil || payData.OrderId == "" || payData.PayType <= 0 ||
		payData.ActualPayPrice <= 0 || payData.PaySerialNumber == "" {
		common.Log.Infof("[func=PaySuccess] params==========%v \r\n", params)
		return common.RenderRpcJson("", "params error", common.ErrorFormData, false)
	}

	orderData, errPay := service.NewOrderService().PaySuccess(payData)
	if errPay != nil {
		return common.RenderRpcJSONError(errPay)
	}

	return common.RenderRpcJson(orderData, "success", 0, true)
}
