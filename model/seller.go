package model

type Seller struct {
	SellerId        	int
	SellerName      	string
	ServiceRate     	float32
	PrepareGoodTime 	int
	Mobile          	string
	VendorType      	int //商家类型 1:外卖 2:预售
	AvailabilityFoodTim int //可取货时间
	StockingTime        int //备货用时时间
	StockingTimeType    int //备货用时类型 1:小时 2:天,
	Email               string //邮箱
	BookTimeType        int //可预订时间类型 1:固定时间范围 2:下单后N天
	BookTimeTalue       int //可预订时间值,
}

