package project

type RpcSeller struct {
	SellerID              int      `json:"seller_id"`
	SellerName            string   `json:"seller_name"`
	ServiceRate           string   `json:"service_rate"`
	PrepareGoodsTime      int      `json:"prepare_goods_time"`
	Mobile                string   `json:"mobile"`
	SellerType            int      `json:"seller_type"`
	AvailabilityGoodsTime int      `json:"availability_goods_time"`
	StockingTime          int      `json:"stocking_time"`
	StockingTimeType      int      `json:"stocking_time_type"`
	Email                 []string `json:"email"`
	BookTimeType          int      `json:"book_time_type"`//可预订时间类型 1:固定时间范围 2:下单后N天
	BookTimeValue         string   `json:"book_time_value"`//可预订时间值 2018-11-12~2019-12-22格式或者　整型表示N天
	Status          	  int      `json:"status"`
	BusinessTime		  string   `json:"business_time"`//00:00~12:00,12:00~23:59 或者　00:00~12:00
	InBusiness			  int 	   `json:"in_business"`
	UseWeek			 	  int 	   `json:"use_week"`
	WeekSet			 	  string   `json:"week_set"`
}
type RpcSellers []RpcSeller

type RpcSellerOperatingStatus struct {
	Status				int			`json:"status"`
	BusinessTime		string		`json:"business_time"`
}

