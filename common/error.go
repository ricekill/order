package common

const (
	ErrorUnauthorized     = 200001
	ErrorAuthLoginFalse   = 200002
	ErrorAuthUnionIdFalse = 200003
	ErrorAuthSendSmsFalse = 200004

	ErrorRedisErr		  = 201001

	ErrorAuthSessionKeyFalse        = 200005
	ErrorAuthPhoneNumberFalse       = 200007
	ErrorAuthOpenidFalse            = 200008
	ErrorAuthSignFalse              = 200009
	ErrorAuthPhoneNumberEmpty       = 200010

	ErrorFormData           	    = 200011
	ErrorCartEmpty          	    = 200012
	ErrorVendorEmpty			    = 200013
	ErrorOrderMsg           	    = 200014
	ErrorTakeFoodTimeF			    = 200015
	ErrorCartRpcServerConRefused    = 200016
	ErrorVendorRpcServerConRefused  = 200017
	ErrorVendorNotExist			    = 200018
	ErrorVendorNotOnSale		    = 200019
	ErrorTakeFoodTimeError		    = 200020
	ErrorProductCheck			    = 200021
	ErrorProductActivity		    = 200022
	ErrorProductInfo			    = 200023
	ErrorProductSoldOut			    = 200024
	ErrorProductStock			    = 200025
	ErrorBuyerInfo                  = 200026
	ErrorProductTotalPrice		    = 200027 //購買的商品總價值有誤！
	ErrorProductTotalPriceLt4Error  = 200028 //訂單金額不小於4元
	ErrorProductTotalCounts		    = 200029 //購買的商品總數量有誤！
	ErrorPTotalPGt20000Error	    = 200030 //訂單金額不小於4元
	ErrorUpdateDB				    = 200040
	ErrorVendorNotOnSaleTime	    = 200041
	ErrorFoodStockInsufficient	    = 200042
	ErrorDoUseCouponError		    = 200043
	ErrorCouponInvalid		  	    = 200044

	ErrorOrderDoesNotExist		    = 200101
	ErrorOrderPayPriceErr		    = 200102

	ErrorOrderStatusError		    = 200201

	ErrorAdvanceOrderTimeError	    = 200301

	ErrorTakeGoodsCode			    = 200401
	ErrorTakeGoodsCodeFailure1		= 200402
	ErrorTakeGoodsCodeFailure2		= 200403
	ErrorTakeGoodsCodeFailure3		= 200404


	ErrorUpdateDataFail				= 300001
)

var errorText = map[int]string{
	ErrorUnauthorized:              "认证失败",
	ErrorAuthLoginFalse:            "授权失败",
	ErrorAuthUnionIdFalse:          "获取不到unionId",
	ErrorAuthSendSmsFalse:          "短信发送失败",
	ErrorAuthSessionKeyFalse:       "session_key获取失败",
	ErrorAuthPhoneNumberFalse:      "解析手机号出错",
	ErrorAuthPhoneNumberEmpty:      "获取手机号失败",
	ErrorAuthOpenidFalse:           "openid获取失败",
	ErrorAuthSignFalse:             "校验失败",
	ErrorFormData:                  "请求参数错误",
	ErrorCartEmpty:                 "购物车为空",
	ErrorVendorEmpty:			    "商家信息为空",
	ErrorOrderMsg:				    "訂單描述長度超過限制",
	ErrorTakeFoodTimeF:             "購物車中的產品不在售賣時間段，請重新選購商品！",
	ErrorCartRpcServerConRefused:   "购物车服务尚未开启",
	ErrorVendorRpcServerConRefused: "商家服务尚未开启",
	ErrorVendorNotExist:		    "商家不存在",
	ErrorVendorNotOnSale:		    "商家已打烊",
	ErrorVendorNotOnSaleTime:	    "商家取货时间不运业",
	ErrorTakeFoodTimeError:		    "您填寫的取餐時間已過期，請重新填寫！",
	ErrorProductCheck:			    "商品信息有誤，無法下單，請重新選擇商品！",
	ErrorProductActivity:           "部分商品活动立减有误，無法下單！",
	ErrorProductInfo:               "部分商品信息有誤，無法下單！",
	ErrorProductSoldOut:	   	    "部分商品已售罄，无法下单！",
	ErrorProductStock:	   		    "部分商品庫存不足，無法下單！",
	ErrorBuyerInfo   :              "买家信息错误，無法下單！",
	ErrorProductTotalPrice:		    "購買的商品總價值有誤！",
	ErrorProductTotalPriceLt4Error: "訂單金額不小於4元",
	ErrorProductTotalCounts:	    "購買的商品總數量有誤！",
	ErrorPTotalPGt20000Error:       "訂單金額不与於20000元",
	ErrorUpdateDB:				    "更新数据库失败",
	ErrorFoodStockInsufficient:	    "部分商品库存不足",
	ErrorDoUseCouponError:		    "优惠价兑奖失败",
	ErrorCouponInvalid:			    "优惠券无效",
	ErrorOrderDoesNotExist:		    "订单不存在",
	ErrorOrderPayPriceErr:		    "支付金额不对",
	ErrorRedisErr:				    "redis请求失败",
	ErrorOrderStatusError:		    "当前订单状态不对",
	ErrorAdvanceOrderTimeError:	    "预订单时间不对",
	ErrorTakeGoodsCode:				"取餐码错误",
	ErrorTakeGoodsCodeFailure1:		"該取餐碼已失效",
	ErrorTakeGoodsCodeFailure2:		"該取餐碼已使用",
	ErrorTakeGoodsCodeFailure3:		"二維碼錯誤",
	ErrorUpdateDataFail:			"更新数据失败",
}

func ErrorText(code int) string {
	return errorText[code]
}

type BadRequestError struct {
	code int
}

func (e BadRequestError) Error() string {
	return ErrorText(e.code)
}

func (e BadRequestError) Code() int {
	return e.code
}

func NewError(code int) error {
	return BadRequestError{code: code}
}

