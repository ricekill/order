package project

import (
	"order-backend/common"
	"order-backend/model"
	"order-backend/model/project"
	"order-backend/repositorie"
	Rpc "order-backend/rpcserver/project"
	"fmt"
	"github.com/zheng-ji/goSnowFlake"
	"strconv"
	"strings"
)

type ProjectOrder struct {
	AppInfo model.App
}

func (o *ProjectOrder) OrderCountById(id int, idType string, filter map[string]interface{}) (int64) {
	var filterW, idKey string
	var filterWArr []string
	if len(filter) > 0 {
		filterWArr = append(filterWArr, fmt.Sprintf("%s = %d", idKey, id))
		for fiK, fiV := range filter {
			if fiV != "" {
				filterWArr = append(filterWArr, fmt.Sprintf("%s = %s", fiK, fiV))
			}
		}
		filterW = strings.Join(filterWArr, " AND ")
	}

	if idType == "seller" {
		//idKey="seller_id"
		totalCount, _ := repositorie.NewDefaultOrderRepositories("").OrderSearchBySellerId(id, filterW)
		return totalCount
	} else {
		//idKey="buyer_id"
		totalCount, _ := repositorie.NewDefaultOrderRepositories("").OrderSearchByBuyerId(id, filterW)
		return totalCount
	}

}

func (o *ProjectOrder) OrderList(
	id int, idType string, pageParam common.PageParam,
	orderBy map[string]interface{}, filter map[string]interface{}) ([]model.Order, *common.Pagination) {
	//ordersList:=make([]model.Order,0)
	var filterW, orderbyW, idKey string
	var filterWArr, orderbyWArr []string
	if idType == "seller" {
		idKey = "seller_id"
	} else {
		idKey = "buyer_id"
	}
	filterWArr = append(filterWArr, fmt.Sprintf("%s = %d", idKey, id))
	for obK, obV := range orderBy {
		if obV != "" {
			orderbyWArr = append(orderbyWArr, fmt.Sprintf("%s = %s", obK, obV))
		}
	}
	for fiK, fiV := range filter {
		if fiV != "" {
			filterWArr = append(filterWArr, fmt.Sprintf("%s = %s", fiK, fiV))
		}
	}
	filterW = strings.Join(filterWArr, " AND ")
	orderbyW = strings.Join(orderbyWArr, ",")
	ordersList, pagination, _ := repositorie.NewDefaultOrderRepositories("").OrderSearch(filterW, orderbyW, pageParam)
	return ordersList, pagination
}
func (o *ProjectOrder) OrderCreateSuffixId(platformId int, buyerId int) (string, error) {
	var pSuffix = fmt.Sprintf("%02d", platformId)
	var bSuffix = fmt.Sprintf("%04d", buyerId)

	iw, err := goSnowFlake.NewIdWorker(1)
	if err != nil {
		//日志
		return "", err
	}
	suffixId, err := iw.NextId()
	if err != nil {
		return "", err
	}
	var suffix = strconv.FormatInt(suffixId, 10) + string(pSuffix) + string(bSuffix)
	return suffix, nil
}

func (o *ProjectOrder) getVendorInfo(SellerId int) (project.RpcSeller, error) {
	ErrBool = false
	vendorInfo, errVendor := new(Rpc.RpcSellerService).GetVendorInfo(SellerId)
	if ErrBool {
		return vendorInfo, common.NewError(common.ErrorVendorEmpty)
	}

	if errVendor != nil {
		return vendorInfo, errVendor
	}

	return vendorInfo, nil
}

//订单完成统计用户下单数已经商品下单数
func (o *ProjectOrder) UpdateDataByOrderComplete(orderId string, sellerId int, userId int) error {
	orderDetailData, err := repositorie.NewDefaultOrderRepositories(orderId).GetOrderDetailById(orderId)
	if err != nil {
		return err
	}

	if  len(orderDetailData) == 0 {
		common.NewError(common.ErrorOrderDoesNotExist)
	}

	productCountInfo := make(map[int]int)

	for _, orderDetail := range orderDetailData {
		productId := orderDetail.ProductId
		productCount := orderDetail.ProductCount
		productCountInfo[productId] += productCount
	}

	if len(productCountInfo) > 0 {
		err = new(Rpc.RpcProductsService).UpdateProductCount(productCountInfo, sellerId)
		if err != nil {
			return err
		}
	}

	if userId > 0 {
		err = new(Rpc.RpcBuyerService).DownOrderCount(userId)
		if err != nil {
			return err
		}
	}

	return nil
}
