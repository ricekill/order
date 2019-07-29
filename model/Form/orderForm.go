package Form

type OrderCreateData struct {
	BuyerId				int  	`json:"user_id"     			binding:"required"` //用户id
	SellerId			int  	`json:"vendor_id"   			binding:"required"` //商家id
	OrderProductCounts	int 	`json:"order_product_counts"	binding:"required"` //订单商品总数量
	OrderGrandTotal		float64 `json:"order_grand_total"		binding:"required"` //订单总金额(订单要付款的价格)
	PackingFee			float64 `json:"packing_fee"`								//包装费(没有默认0)
	CouponsFee			float64 `json:"coupons_fee"`								//优惠券优惠价格(没有默认0)
	CouponUserCode		string  `json:"coupon_user_code"`							//优惠券码
	OrderMsg			string  `json:"order_msg"`        							//订单描述
	TakeGoodsTime		int64  	`json:"take_food_time"			binding:"required"` //预约取货时间（传送格式 2018-03-26 23:50:00）
}

type PaySuccessData struct {
	OrderId         string  `json:"order_id"`
	PayType         int     `json:"pay_type"`
	ActualPayPrice  float64 `json:"actual_pay_price"`
	PaySerialNumber string  `json:"pay_serial_number"`
}

type AcceptOrderData struct {
	OrderId      string `json:"order_id"`
	NewOrderType int    `json:"new_order_type"`
	SellerId     int    `json:"vendor_id"`
	IsBackData   int    `json:"is_back_data"`
}

type TakeGoodByCodeData struct {
	TakeGoodsCode int `json:"take_food_code"`
	SellerId      int `json:"vendor_id"`
}

type ConfirmTakeGoodData struct {
	OrderId  string `json:"order_id"`
	SellerId int 	`json:"vendor_id"`
}
