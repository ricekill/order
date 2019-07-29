package project

import "encoding/json"

type RpcProduct struct {
	Status            int           `json:"status"`
	Price             json.Number   `json:"price"`
	LimitedStock      int           `json:"limited_stock"`
	Stock             int           `json:"stock"`
	IsWine            int           `json:"is_wine"`
	ID                int           `json:"id"`
	VendorID          int           `json:"vendor_id"`
	HasPackage        int           `json:"has_package"`
	Type              int           `json:"type"`
	Name              string        `json:"name"`
	Logo              string        `json:"logo"`
	OfferActivityInfo struct {
		OfferActivityID int    `json:"offer_activity_id"`
		OfferType       int    `json:"offer_type"`
		OfferValue      string `json:"offer_value"`
		LimitedNumber   int    `json:"limited_number"`
		DiscountPrice   json.Number    `json:"discount_price"`
	} `json:"offer_activity_info"`
	IsDeleted         int           `json:"is_deleted"`
	TypeName          string        `json:"type_name"`
	Num               string        `json:"num"`
	SaleTime          string        `json:"sale_time"`
	PackageInfo       []struct {
		ID          int    `json:"id"`
		FoodID      int    `json:"food_id"`
		PackageName string `json:"package_name"`
		PackageType int    `json:"package_type"`
		PackageTags []struct {
			TagName  		  string 		`json:"tag_name"`
			TagPrice 		  json.Number 	`json:"tag_price"`
			TagID    		  string 		`json:"tag_id"`
			FoodID   		  int 			`json:"food_id"`
			LimitedStock      int           `json:"limited_stock"`
			Stock             int           `json:"stock"`
		} `json:"package_tags"`
	} `json:"package_info"`
}
type RpcProducts []RpcProduct
