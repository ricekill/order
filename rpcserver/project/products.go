package project

import (
	"order-backend/common"
	"order-backend/common/itface/remote"
	"encoding/json"
	"fmt"
	"github.com/hprose/hprose-golang/rpc"
	"order-backend/model/project"
	"strconv"
	"strings"
	"time"
)

type RpcProductsService struct {
	RpcClient     rpc.Client
	VendorInfo    project.RpcSeller
	TakeGoodsTime int64               //取参时间
	BusinessTime  []map[string]string //营业时间
	Carts         project.RpcCart
}

func (reus *RpcProductsService) Init(attrs map[string]interface{}) {
	if attrs["VendorInfo"] != nil {
		reus.VendorInfo = attrs["VendorInfo"].(project.RpcSeller)
	}
	if attrs["TakeGoodsTime"] != nil {
		reus.TakeGoodsTime = attrs["TakeGoodsTime"].(int64)
	}
	if attrs["BusinessTime"] != nil {
		reus.BusinessTime = attrs["BusinessTime"].([]map[string]string)
	}
	if attrs["Carts"] != nil {
		reus.Carts = attrs["Carts"].(project.RpcCart)
	}
}

func (reus *RpcProductsService) getProductsRpcClient() remote.RpcProductsService {
	var productsService remote.RpcProductsService
	reus.RpcClient = rpc.NewClient(common.Config.RpcServer.Project.Products)
	reus.RpcClient.UseService(&productsService)
	return productsService
}

func (reus *RpcProductsService) GetProductsInfo() (map[string]interface{}, map[string]interface{}, error) {
	//fmt.Printf("========reus.VendorInfo===================%v",reus.VendorInfo)
	var pids []int
	//fmt.Printf("\n\n\n\n=========reus.Carts=%v======\n\n\n\n\n",reus.Carts)
	if len(reus.Carts) > 0 {
		for i := 0; i < len(reus.Carts); i++ {
			cart := reus.Carts[i]
			pid, _ := strconv.Atoi(cart.ProductID)
			pids = append(pids, pid)
		}
	} else {
		return nil, nil, common.NewError(common.ErrorCartEmpty)
	}
	//fmt.Printf("\n\n\n\n=========pids=%v======\n\n\n\n\n",pids)
	productsInfo := reus.getProductsRpcClient().Products_getDetails(pids, reus.VendorInfo.SellerID)
	//fmt.Printf("==========productsInfo===============%v",productsInfo)
	//定义返回信息数据结构struct
	var rpcProductsInfo []project.RpcProduct
	json.Unmarshal([]byte(productsInfo), &rpcProductsInfo)
	//fmt.Printf("\n\n\n\n====rpcGetProductsInfo=====productsInfo=%v======\n\n\n\n\n",productsInfo)
	//格式化产品获取的价格数组
	products, foodStock, activity, err := reus.formatProductItems(rpcProductsInfo)
	if err != nil {
		return nil, nil, err
	}
	//格式购物车的items最终的数组
	cartItems := reus.formatCartsItems()
	//fmt.Printf("\n\n\n\n====GetProductsInfo=====cartItems=%v======products=%v==\n\n\n\n\n",cartItems,products)
	if cartItems != nil && len(products) > 0 {
		result, orderFoodStock, err := reus.combineProduct(cartItems, products, foodStock, activity)
		//fmt.Printf("\n\n\n\n====GetProductsInfo-combineProduct=====result=%v===orderFoodStock=%v==err=%v==\n\n\n\n\n",result,orderFoodStock,err)
		////商品合计
		totalProductInfo, _ := reus.getTotalProductInfo(result)
		return totalProductInfo, orderFoodStock, err
	}
	defer reus.RpcClient.Close()
	return nil, nil, common.NewError(common.ErrorProductCheck)
}

func (reus *RpcProductsService) ModifyProductsStock(foodStocks map[int]int, action string, sellerId int) error {
	returnStr := reus.getProductsRpcClient().Products_modifyStock(foodStocks, action, sellerId)
	defer reus.RpcClient.Close()
	var productsStock project.RpcReturn
	json.Unmarshal([]byte(returnStr), &productsStock)
	//如果message不为空,返回错误信息
	if !productsStock.Status {
		fmt.Printf("\n=================modifyStock=========%v \r\n\r\n", returnStr)
		return common.NewError(common.ErrorFoodStockInsufficient)
	}
	return nil
}

//获取整合后的订单数据商品信息
func (reus *RpcProductsService) getTotalProductInfo(productItems []map[string]interface{}) (map[string]interface{}, error) {
	var productInfo map[string]interface{}
	var products []map[string]interface{}
	productInfo = make(map[string]interface{})
	if len(productItems) > 0 {
		total := float64(0.00)
		orderProductCounts := 0
		for i := 0; i < len(productItems); i++ {
			item := productItems[i]
			mm := make(map[string]interface{})
			mm["product_id"] = item["product_id"]
			mm["product_parent_id"] = item["product_parent_id"]
			mm["product_price"] = item["product_price"]
			mm["product_name"] = item["product_name"]
			mm["product_counts"] = item["product_counts"]
			mm["product_msg"] = item["product_msg"]
			mm["product_real_id"] = item["product_real_id"]
			mm["product_original_price"] = item["product_original_price"]
			productActionId, ok := item["product_action_id"].(int)
			if !ok {
				mm["product_action_id"] = 0
			} else {
				mm["product_action_id"] = productActionId
			}
			//////////////////////////////
			productPrice, _ := item["product_price"].(json.Number).Float64()
			total += float64(item["product_counts"].(int)) * productPrice
			products = append(products, mm)

			//////////orderProductCounts////////////////////////////
			if item["product_parent_id"] == "0" {
				orderProductCounts += item["product_counts"].(int)
			}
		}
		productInfo["items"] = products
		//订单商品总金额(此订单所有商品的金额)
		productInfo["order_amount_total"] = total
		productInfo["order_product_counts"] = orderProductCounts
	}
	//fmt.Printf("\n\n\n\n=====111====productInfo===%v============\n\n\n\n",productInfo)
	return productInfo, nil

}
func (reus *RpcProductsService) combineProduct(
	cartItems map[string]map[string]interface{},
	products map[string]map[string]interface{},
	foodStock map[string]map[string]interface{},
	activityInfo map[string]map[string]interface{}) ([]map[string]interface{}, map[string]interface{}, error) {

	var result []map[string]interface{}
	var offerActivityData map[string]interface{}
	var orderFoodStock map[string]interface{}
	var isActivity bool
	isActivity = false
	orderFoodStock = make(map[string]interface{})
	fmt.Printf("\n\n\n===zpm=ok====cartItems=========%v====\n\n\n\n\n\n\n", cartItems)
	for key, item := range cartItems {
		itemProductId := item["product_id"].(string)
		productsMap, ok := products[itemProductId]
		//fmt.Printf("\n\n\n===zpm=ok=%v==cartItems=%v=======products=%v====\n\n\n\n\n\n\n",ok,cartItems,products)
		if ok {
			//fmt.Printf("\n\n\n===zpm=productsMap=%v===============\n\n\n\n\n\n\n",productsMap)
			productItem := productsMap
			var product map[string]interface{}
			//组装购物车和产品的数据为最终数据
			//从购物车获取数据
			product = make(map[string]interface{})
			product["product_counts"] = item["product_num"]
			product["product_id"] = key                          //产品id为购物车id加产品id（购物车id_产品id或者购物车id_套餐id-选项id）
			product["product_real_id"] = item["product_real_id"] //产品真实的id
			product["product_msg"] = item["product_msg"]         //产品描述
			product["product_parent_id"] = "0"                   //父级默认为0不需要加购物车id
			//关联的父级id(购物车id_原来产品的父级id)
			if productItem["product_parent_id"].(string) != "0" {
				product["product_parent_id"] = fmt.Sprintf("%d_%s", item["id"].(int), productItem["product_parent_id"].(string))
			}
			//从产品获取数据
			product["product_price"] = productItem["product_price"]
			product["product_name"] = productItem["product_name"]
			product["product_original_price"] = productItem["product_price"]
			product["product_action_id"] = 0

			//判断有没有活动立减
			offerActivityId, err := strconv.Atoi(item["offer_activity_id"].(string))
			if err != nil {
				//fmt.Printf("\n\n\n===strconv.Atoi=offerActivityId=%v===============\n\n\n\n\n\n\n",err)
				offerActivityId = 0
			}
			if offerActivityId != 0 {
				if iPid := itemProductId; activityInfo[iPid] != nil &&
					strconv.Itoa(activityInfo[iPid]["offer_activity_id"].(int)) == item["offer_activity_id"].(string) {
					var tmp = map[string]interface{}{
						"product_num":       "0",
						"offer_activity_id": item["offer_activity_id"].(string),
					}
					oo, ok := offerActivityData[itemProductId].(map[string]interface{})
					if !ok {
						oo = tmp
					}
					productNum, err := strconv.Atoi(oo["product_num"].(string))
					if err != nil {
						productNum = 0
					}
					productItemNum := item["product_num"].(int)

					oo["product_num"] = productNum + productItemNum
					activity := activityInfo[item["product_id"].(string)]
					product["product_price"] = activity["discount_price"]
					product["product_action_id"] = activity["offer_activity_id"]

				} else {
					isActivity = true
					break
				}
			}
			//fmt.Printf("\n\n\n========combineProduct=========product=%v===============\n\n\n\n\n\n\n",product)
			result = append(result, product)
			//记录库存
			if item["product_real_id"].(string) != "0" {
				productItemNum := item["product_num"].(int)
				om, ok := orderFoodStock[item["product_real_id"].(string)].(interface{})
				if !ok {
					orderFoodStock[item["product_real_id"].(string)] = strconv.Itoa(productItemNum)
				} else {
					productRealId, err := strconv.Atoi(om.(string))
					if err != nil {
						productRealId = 0
					}
					orderFoodStock[item["product_real_id"].(string)] = strconv.Itoa(productRealId + productItemNum)
				}
				//fmt.Printf("\n\n\n\n=====item-product_real_id-=%v====orderFoodStock=%v====ok=%v====\n\n\n\n",item["product_real_id"],orderFoodStock,ok)
			}
		}
	}
	for productId, aInfo := range offerActivityData {
		activity, ok := activityInfo[productId]
		aInfo, ok2 := aInfo.(map[string]interface{})
		if ok && ok2 {
			productNum, err := strconv.Atoi(aInfo["product_num"].(string))
			if err != nil {
				productNum = 0
			}
			limitedNumber, err := strconv.Atoi(activity["limited_number"].(string))
			if err != nil {
				limitedNumber = 0
			}
			if activity["limited_number"].(string) != "0" && (productNum > limitedNumber) {
				isActivity = true
				break
			}
		}
	}

	if isActivity {
		return nil, nil, common.NewError(common.ErrorProductActivity)
	}

	if len(cartItems) != len(result) {
		//fmt.Printf("\n\n\n\n\n=combineProduct=======cartItems=%v\n=====result=%v=\n=products=%v=\n=====\n\n\n\n\n",cartItems,result,products)
		return nil, nil, common.NewError(common.ErrorProductInfo)
	}

	var isSoldOut, isStockInsufficient bool

	isSoldOut = false
	isStockInsufficient = false
	for foodId, foodCount := range orderFoodStock {
		ff, ok := foodStock[foodId]
		if ok {
			ffStock, err := strconv.Atoi(ff["stock"].(string))
			if err != nil {
				ffStock = 0
			}
			if ff["limited_stock"] == "1" {
				if ffStock == 0 {
					isSoldOut = true
				} else if foodCount.(int) > ffStock {
					isStockInsufficient = true
				}
			}
		}
	}
	if isSoldOut {
		return nil, nil, common.NewError(common.ErrorProductSoldOut)
	}

	if isStockInsufficient {
		return nil, nil, common.NewError(common.ErrorProductStock)
	}
	return result, orderFoodStock, nil
}
func (reus *RpcProductsService) formatCartsItems() map[string]map[string]interface{} {
	var data map[string]map[string]interface{}
	data = make(map[string]map[string]interface{})
	//fmt.Printf("\n\n\n\n\n============formatCartsItems===carts=%v===================\n\n\n\n\n",reus.Carts)
	if len(reus.Carts) > 0 {
		for i := 0; i < len(reus.Carts); i++ {
			cartItem := reus.Carts[i]
			//购物车id_产品id代表一条记录
			productKey := fmt.Sprintf("%d_%s", cartItem.ID, cartItem.ProductID)
			mm, ok := data[productKey]
			if !ok {
				mm = make(map[string]interface{})
				data[productKey] = mm
			}
			mm["id"] = cartItem.ID //购物车id
			mm["product_id"] = cartItem.ProductID
			mm["product_num"] = cartItem.ProductNum
			mm["product_msg"] = cartItem.ProductMsg
			mm["product_real_id"] = cartItem.ProductID
			mm["offer_activity_id"] = strconv.Itoa(cartItem.OfferActivityId)
			if cartItem.HasPackage {
				if len(cartItem.Package) > 0 {
					for i := 0; i < len(cartItem.Package); i++ {
						currentPackage := cartItem.Package[i]
						if len(currentPackage.SonPackage) > 0 {
							for j := 0; j < len(currentPackage.SonPackage); j++ {
								currentSonPItem := currentPackage.SonPackage[j]
								productPackageKey := fmt.Sprintf("%d_%d-%s", cartItem.ID, currentPackage.ID, currentSonPItem.Sid)
								//这里的产品id为套餐id-选项id 键值用购物车id_套餐id-选项id 区分
								//（有可能产品id相同的情况下是一个套餐的）
								mm2, ok := data[productPackageKey]
								if !ok {
									mm2 = make(map[string]interface{})
									data[productPackageKey] = mm2
								}
								mm2["id"] = cartItem.ID //购物车id
								mm2["product_id"] = fmt.Sprintf("%d-%s", currentPackage.ID, currentSonPItem.Sid)
								mm2["product_num"] = cartItem.ProductNum
								mm2["product_msg"] = ""
								mm2["product_real_id"] = strconv.Itoa(currentSonPItem.FoodID)
								mm2["offer_activity_id"] = ""
							}
						}
					}
				}
			}
		}
		return data
	} else {
		return nil
	}
}

/**
改变三维数组的产品信息为二维数组
 */
func (reus *RpcProductsService) formatProductItems(productsInfo []project.RpcProduct) (
	map[string]map[string]interface{},
	map[string]map[string]interface{},
	map[string]map[string]interface{},
	error) {
	var data, foodStock, activity map[string]map[string]interface{}
	data = make(map[string]map[string]interface{})
	foodStock = make(map[string]map[string]interface{})
	activity = make(map[string]map[string]interface{})
	if len(productsInfo) > 0 {
		for i := 0; i < len(productsInfo); i++ {
			pItem := productsInfo[i]
			_, err := reus.getProductSaleTime(pItem)
			if err != nil {
				return nil, nil, nil, err
			}
			pItemID := strconv.Itoa(pItem.ID)
			mm, ok := data[pItemID]
			if !ok {
				mm = make(map[string]interface{})
				data[pItemID] = mm
			}
			mm["product_id"] = pItemID
			mm["product_parent_id"] = "0"
			mm["product_price"] = pItem.Price
			mm["product_name"] = pItem.Name

			ff, ok := foodStock[pItemID]
			if !ok {
				ff = make(map[string]interface{})
				foodStock[pItemID] = ff
			}
			ff["limited_stock"] = pItem.LimitedStock
			ff["stock"] = strconv.Itoa(pItem.Stock)

			//立减活动信息
			if pItem.OfferActivityInfo.OfferActivityID > 0 {
				aa, ok := activity[pItemID]
				if !ok {
					aa = make(map[string]interface{})
					activity[pItemID] = aa
				}
				aa["offer_activity_id"] = pItem.OfferActivityInfo.OfferActivityID
				aa["limited_number"] = pItem.OfferActivityInfo.LimitedNumber
				aa["discount_price"] = pItem.OfferActivityInfo.DiscountPrice
			}

			if pItem.HasPackage > 0 {
				if len(pItem.PackageInfo) > 0 {
					for i := 0; i < len(pItem.PackageInfo); i++ {
						currentPackage := pItem.PackageInfo[i]
						if len(currentPackage.PackageTags) > 0 {
							for j := 0; j < len(currentPackage.PackageTags); j++ {
								currentTag := currentPackage.PackageTags[j]
								currentPackageTagId := fmt.Sprintf("%d-%s", currentPackage.ID, currentTag.TagID)
								tt, ok := data[currentPackageTagId]
								if !ok {
									tt = make(map[string]interface{})
									data[currentPackageTagId] = tt
								}
								tt["product_id"] = currentPackageTagId
								tt["product_parent_id"] = pItemID
								tt["product_price"] = currentTag.TagPrice
								tt["product_name"] = currentTag.TagName

								if currentTag.FoodID > 0 {
									tagFoodId := strconv.Itoa(currentTag.FoodID)
									ff, ok := foodStock[tagFoodId]
									if !ok {
										ff = make(map[string]interface{})
										foodStock[tagFoodId] = ff
									}
									ff["limited_stock"] = currentTag.LimitedStock
									ff["stock"] = strconv.Itoa(currentTag.Stock)
								}

							}
						}
					}
				}
			}
		}
	}
	//fmt.Printf("\n\n\n\n\n==products=formatProductItems=370=data=%v====foodStock=%v===activity=%v====\n\n\n\n\n",data,foodStock,activity)
	return data, foodStock, activity, nil
}

/**
获取产品的供应时间
 */
func (reus *RpcProductsService) getProductSaleTime(productsInfo project.RpcProduct) (map[string]string, error) {
	if reus.VendorInfo.SellerType == 1 { //商家类型 1:外卖 2:预售
		takeGoodsTime := reus.TakeGoodsTime //string　日期格式
		prepareFoodTime := reus.VendorInfo.PrepareGoodsTime
		businessTime := reus.BusinessTime
		//转换成时间绰
		var timeT = int(takeGoodsTime)
		//取餐时间>=产品供应时间+备餐时间
		prepareStartTime := prepareFoodTime * 60
		//取餐时间<=产品供应结束时间+备餐时间+20分钟 默认的取餐时间要加上十分钟
		prepareEndTime := (prepareFoodTime + 20) * 60
		saleTime := productsInfo.SaleTime
		//fmt.Printf("businessTime=%v--takeTime=%v---prepareEndTIme=%v----prepareStartTime=%v---saleTime=%v---len(businessTime)=%v",businessTime,timeT,prepareEndTime,prepareStartTime,saleTime,len(businessTime))
		if saleTime == "" {
			for i := 0; i < len(businessTime); i++ {
				splitTime := businessTime[i]
				startTimeStamp, _ := formatBusinessTimeToStamp(splitTime["start_time"])
				endTimeStamp, _ := formatBusinessTimeToStamp(splitTime["end_time"])
				//fmt.Printf("\n\n\n\n\nsplitTimestart_time=%v=====splitTimeend_time=%v\n\n\n\n\n",startTimeStamp,endTimeStamp)
				if timeT >= (int(startTimeStamp)+int(prepareStartTime)) &&
					timeT <= (int(endTimeStamp)+int(prepareEndTime)) {
					return splitTime, nil
				}

			}
		} else {
			var returnData map[string]string
			var startTime, endTime string
			returnData = make(map[string]string)
			splitTime := strings.Split(saleTime, "~")
			if splitTime[0] != "" {
				startTime = splitTime[0]
			} else {
				startTime = "0"
			}
			if splitTime[1] != "" {
				endTime = splitTime[1]
			} else {
				endTime = "0"
			}
			returnData["start_time"] = startTime
			returnData["end_time"] = endTime
			startTimeStamp, _ := formatBusinessTimeToStamp(returnData["start_time"])
			endTimeStamp, _ := formatBusinessTimeToStamp(returnData["end_time"])
			//fmt.Printf("\n\n\n\n\nstartTimeStamp=%v===endTimeStamp=%v=prepareStartTime=%v===prepareEndTime=%v========\n\n\n\n\n",startTimeStamp,endTimeStamp,prepareStartTime,prepareEndTime)
			if (timeT >= (int(startTimeStamp) + int(prepareStartTime))) &&
				timeT <= (int(endTimeStamp)+int(prepareEndTime)) {
				return returnData, nil
			}
		}
		return nil, common.NewError(common.ErrorTakeFoodTimeF)
	}
	return nil, nil
}
func formatBusinessTimeToStamp(timeStr string) (int64, error) {
	now := time.Now().Format("2006-01-02")
	timeP, _ := time.Parse("15:04", timeStr)
	//fmt.Printf("\ntimeStr=%v==timeP.Hour=%v====timeP.Minute=%v==\n",timeStr,timeP.Hour(),timeP.Minute())
	var lastTime string
	if timeP.Hour() < 10 {
		if timeP.Minute() < 10 {
			lastTime = fmt.Sprintf("%s 0%d:0%d", now, timeP.Hour(), timeP.Minute())
		} else {
			lastTime = fmt.Sprintf("%s 0%d:%d", now, timeP.Hour(), timeP.Minute())
		}

	} else {
		if timeP.Minute() < 10 {
			lastTime = fmt.Sprintf("%s %d:0%d", now, timeP.Hour(), timeP.Minute())
		} else {
			lastTime = fmt.Sprintf("%s %d:%d", now, timeP.Hour(), timeP.Minute())
		}

	}

	//fmt.Printf("\nlastTime=%v\n",lastTime)
	timeR, err := time.Parse("2006-01-02 15:04", lastTime)
	if err != nil {
		//fmt.Printf("\n err timeStr=%v\n",timeStr)
		return 0, err
	}
	//fmt.Printf("\nlastTimeStamp=%v\n",timeR.Unix())
	return timeR.Unix(), err
}

func (reus *RpcProductsService) UpdateProductCount(productsInfo map[int]int, sellerId int) error {
	returnStr := reus.getProductsRpcClient().Products_updateCounts(productsInfo, sellerId)
	defer reus.RpcClient.Close()
	var productsStock project.RpcReturn
	json.Unmarshal([]byte(returnStr), &productsStock)
	//如果message不为空,返回错误信息
	if !productsStock.Status {
		fmt.Printf("\n=================modifyStock=========%v \r\n\r\n", returnStr)
		return common.NewError(common.ErrorVendorNotExist)
	}
	defer reus.RpcClient.Close()
	return nil
}