package project

//订单状态
const (
	CCOrderStatus_UNPAID 			 =	1 //待支付
	CCOrderStatus_NEED_ACCEPT		 =	2 //待接單
	CCOrderStatus_PREPARE_GOODS		 =	3 //備餐中
	CCOrderStatus_WAIT_TAKE_GOODS	 =	4 //待取餐
	CCOrderStatus_TAKE_GOODS		 =	5 //已取餐
	CCOrderStatus_DONE				 =	6 //已完成
	CCOrderStatus_CANCEL 			 = 	7 //已取消
	CCOrderStatus_WAIT_PREPARE_GOODS =	8 //待備餐
)

//商家类型
const (
	SELLER_TYPE_TAKE_OUT 	= 1 //外卖
	SELLER_TYPE_PRESELL 	= 2 //预售
)

//商家推送类型
const (
	SELLER_PUSH_NEWORDER 			= 3 //新订单
	SELLER_PUSH_PLATFROM_AFTERSALES = 4 //平台通知商家售后
	SELLER_PUSH_CANCEL 				= 5 //取消订单
	SELLER_PUSH_ACCEPT 				= 6 //超时未接单 转 售后订单
	SELLER_PUSH_ADVANCE 			= 7 //预订单到时间 转 待备餐
	SELLER_PUSH_TAKEFOOD 			= 8 //待备餐 到時 轉 待取餐
	SELLER_PUSH_PRINT 				= 9 //只打印
)

//用户端推送
const (
	BUYER_PUSH_CANCEL 		= 11 //商家超时未接单，拒单，取消
	BUYER_PUSH_BEFORE 		= 12 //在用户取餐时间前五分钟推送取餐提醒
	BUYER_PUSH_READY 		= 13 //订单超过预定时间15分钟后还未确认收货
	BUYER_PUSH_AFTERSALES 	= 14 //用戶申請售後成功/失敗，推送消息
	BUYER_PUSH_REMIND 		= 15 //晚上8点发送，推送消息
)

const (
	DEFAULT_ADVANCE_ORDER_UNIT_TIME 	= 40 //判断是否预订单加40分钟
	DEFAULT_ADVANCE_ORDER_DELAY_TIME 	= 30 //预订单到时候延迟30秒
)

const (
	ORDER_ACCEPT_NEW 		= 1 //新订单
	ORDER_ACCEPT_ADVANCE 	= 2 //预订单
	ORDER_ACCEPT_PREPARE 	= 3 //待备餐
)