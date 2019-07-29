package model

type Order struct {
	Id                    int     `xorm:"not null pk autoincr INT(10)"`
	OrderId               string  `xorm:"not null comment('订单id') unique CHAR(25)"`
	PlatformId            int     `xorm:"not null comment('平台id') index(order_platform_user_status_index) INT(10)"`
	BuyerId               int     `xorm:"not null comment('买家id') index(order_platform_user_status_index) INT(10)"`
	BuyerName             string  `xorm:"comment('买家名字') VARCHAR(20)"`
	SellerId              int     `xorm:"not null comment('卖家id') INT(10)"`
	SellerName            string  `xorm:"comment('卖家名字') VARCHAR(20)"`
	CellphoneNumber       string  `xorm:"comment('买家手机号') VARCHAR(20)"`
	TradeType             int     `xorm:"not null comment('行业类型:1.餐饮 2.零售') TINYINT(3)"`
	SceneType             int     `xorm:"not null comment('场景类型:1.堂食 2.外卖') TINYINT(3)"`
	SellerType            int     `xorm:"not null comment('商店类型 1.外卖 2.预售') TINYINT(3)"`
	OrderStatus           int     `xorm:"not null comment('订单状态:1.待支付 2.待接单 3.备餐中 4.待取餐 5.已取餐 6.已完成 7.已取消 8.待备餐') index(order_platform_user_status_index) TINYINT(3)"`
	AdvanceOrderStatus    int     `xorm:"not null comment('预订单状态:0.新订单 1.预订单 2.待备餐') TINYINT(3)"`
	TakeGoodsCode         int     `xorm:"comment('取餐码') INT(10)"`
	TakeGoodsTime         int     `xorm:"comment('预约取餐时间') INT(10)"`
	PrepareGoodsTime      int     `xorm:"comment('配餐时间') INT(10)"`
	ServiceRate           float64 `xorm:"comment('服务费率(整数百分比如5%存5)') DOUBLE(4,2)"`
	OrderNumeral          int     `xorm:"not null default 0 comment('订单排号') MEDIUMINT(8)"`
	OrderCompleteTime     int     `xorm:"default 0 comment('订单完成时间') INT(10)"`
	OrderSource           int     `xorm:"not null comment('订单来源:1.H5 2.iOS 3.Android') TINYINT(3)"`
	OrderMsg              string  `xorm:"comment('订单描述') VARCHAR(255)"`
	OrderProductCounts    int     `xorm:"not null comment('订单商品总数量') SMALLINT(5)"`
	OrderGrandTotal       float64 `xorm:"not null comment('订单商品总金额(此订单所有商品的金额)') DECIMAL(15,2)"`
	OrderAmountTotal      float64 `xorm:"not null comment('订单总金额(订单要付款的价格)') DECIMAL(15,2)"`
	PackingFee            float64 `xorm:"not null default 0.00 comment('包装费') DECIMAL(15,2)"`
	ShippingFee           float64 `xorm:"not null default 0.00 comment('配送费') DECIMAL(15,2)"`
	CouponFee             float64 `xorm:"not null default 0.00 comment('优惠券优惠价格') DECIMAL(15,2)"`
	CouponUserCode        string  `xorm:"default '' comment('优惠券码') VARCHAR(32)"`
	CouponCode            string  `xorm:"default '' comment('优惠券规则id') VARCHAR(10)"`
	PayType               int     `xorm:"default 0 comment('支付类型:0.默认值没支付 1.信用卡 2.银联支付 3.微信支付 4.支付宝支付 5.paypal') TINYINT(3)"`
	ActualPayPrice        float64 `xorm:"comment('订单实际付款价格(同订单总金额，付款成功后写入，做校验使用)') DECIMAL(15,2)"`
	OrderPayTime          int     `xorm:"comment('支付时间') INT(10)"`
	PaySerialNumber       string  `xorm:"default ' ' comment('支付流水号') VARCHAR(25)"`
	AfterSalesStatus      int     `xorm:"not null default 0 comment('是否发起售后 0:不发起  1:发起售后 2.售后完成') TINYINT(3)"`
	PaymentsFinalStatus   int     `xorm:"not null default 0 comment('收支类型 0:订单收入 1:取消订单退款 2:拒单退款 3:商家售后退款 4:平台售后退款') TINYINT(3)"`
	Cext                  string  `xorm:"comment('扩展字段') VARCHAR(100)"`
	AvailabilityGoodsTime int     `xorm:"comment('取餐有效期时间') INT(10)"`
	CreatedAt             int     `xorm:"not null index(order_platform_user_status_index) INT(10)"`
	UpdatedAt             int     `xorm:"not null INT(10)"`
}
