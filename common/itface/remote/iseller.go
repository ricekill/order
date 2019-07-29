package remote

import "order-backend/model/project"

type RpcSellerService struct {
	Vendor_getInfo func(int) string //按商家id获取
	Vendor_getMultiInfo func([]int) string//按商家ids获取多个商家信息
	Vendor_getOperatingStatus func(int, int64) string//按商家ids获取多个商家信息
	Vendor_vendorScore func(int) string //计算商家分数
}

type SellerService interface {
	GetVendorInfo(int)  (project.RpcSeller,error)
	GetMultiInfo([]int) (project.RpcSellers,error)
	CalculateSellerScore(int) error
}