package project

import "encoding/json"

//用这个工具生成　https://mholt.github.io/json-to-go/
type RpcCart []struct {
	ID              int         `json:"id"`                //购物车自曾ID
	ProductID       string      `json:"product_id"`        //商品ID
	Price           json.Number `json:"price"`             //商品价格
	Name            string      `json:"name"`              //商品名称
	ProductNum      int         `json:"product_num"`       //加入购物车商品数量
	HasPackage      bool        `json:"has_package"`       //是否套餐
	LimitedNumber   int         `json:"limited_number"`    //优惠限购数量 0不限购
	LimitedStock    int         `json:"limited_stock"`     //0和1,是否限制库存数
	Stock           int         `json:"stock"`             //库存数
	OfferActivityId int         `json:"offer_activity_id"` //优惠活动id
	SaleTime        string      `json:"sale_time"`         //时间段　11:30~14:30
	Package         []struct {
		//套餐选项
		ID         int    `json:"id"`
		Name       string `json:"name"` //套餐名字 eg,可选饮料 必选青菜,购物车界面不显示
		SonPackage []struct {
			//选项
			FoodID       int         `json:"food_id"` //套餐选项FoodID  为关联商品时foodId大于0,否则为0
			Sid          string      `json:"sid"`     //套餐选项ID
			Name         string      `json:"name"`    //套餐选项　003（米兰假日）
			Price        json.Number `json:"price"`
			LimitedStock int         `json:"limited_stock"` //0和1,是否限制库存数
			Stock        int         `json:"stock"`         //当前剩余库存
		} `json:"son_package"`
	} `json:"package"`
	ProductMsg string      `json:"product_msg"`
	OfferType  int         `json:"offer_type"`
	OfferValue json.Number `json:"offer_value"`
}
