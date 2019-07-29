package model

type OrderIndex struct {
	Id                 int    `xorm:"not null pk autoincr INT(10)"`
	OrderId            string `xorm:"not null comment('订单id') unique CHAR(25)"`
	PlatformId         int    `xorm:"not null comment('平台id') INT(10)"`
	BuyerId            int    `xorm:"not null comment('买家id') INT(10)"`
	SellerId           int    `xorm:"not null comment('卖家id') INT(10)"`
	SellerType         int    `xorm:"not null comment('商店类型 1.外卖 2.预售') TINYINT(3)"`
	CellphoneNumber    string `xorm:"comment('买家手机号') VARCHAR(20)"`
	OrderStatus        int    `xorm:"not null comment('订单状态:1.待支付 2.待接单 3.备餐中 4.待取餐 5.已取餐 6.已完成 7.已取消 8.待备餐') TINYINT(3)"`
	AdvanceOrderStatus int    `xorm:"not null comment('预订单状态:0.新订单 1.预订单 2.待备餐') TINYINT(3)"`
	TakeGoodsCode      int    `xorm:"comment('取货码') INT(10)"`
	TakeGoodsTime      int    `xorm:"comment('预约取货时间') INT(10)"`
	OrderNumeral       int    `xorm:"not null default 0 comment('订单排号') MEDIUMINT(8)"`
	OrderSource        int    `xorm:"not null comment('订单来源:1.H5 2.iOS 3.Android') TINYINT(3)"`
	CreatedAt          int    `xorm:"INT(10)"`
	UpdatedAt          int    `xorm:"INT(10)"`
}
