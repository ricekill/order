package project

import (
	"order-backend/common"
	Config "order-backend/config/project"
	"order-backend/model"
	"order-backend/model/Form"
	"order-backend/model/project"
	"order-backend/repositorie"
	Rpc "order-backend/rpcserver/project"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const RedisSellerOrderNumeral = "S-ON:"
var ErrBool = false

/**
 * 创建订单
 */
func (o *ProjectOrder) OrderCreate(orderData Form.OrderCreateData) (oid string, e error) {
	//defer func(){
	//	if err:=recover();err!=nil {
	//		oid = ""
	//		e   = fmt.Errorf("order create error: %v", err)
	//		return
	//	}
	//}()
	//o.OrderCancelReturnStock("293157951718297600010002")
	//return "",nil
	//生成
	platformId := common.APP.Id

	if utf8.RuneCountInString(orderData.OrderMsg) > 200 {
		return "", common.NewError(common.ErrorOrderMsg)
	}

	//获取商家信息
	vendorInfo, errVendor := o.getVendorInfo(orderData.SellerId)
	if errVendor != nil {
		return "", errVendor
	}
	TakeGoodsTime := orderData.TakeGoodsTime

	nowTime := int(time.Now().Unix())
	//校验商品是否在售卖时间 取餐时间>=当前时间
	if int(TakeGoodsTime) < (nowTime + vendorInfo.PrepareGoodsTime*60) {
		return "", common.NewError(common.ErrorTakeFoodTimeF)
	}

	businessTime, err := o.checkIsOnSales(vendorInfo, TakeGoodsTime)

	if err != nil {
		return "", err
	}

	//====================================校验當前商家是否在营业 end

	//预售店检查时间 start ==============================================
	availabilityFoodTime := 0
	if vendorInfo.SellerType == Config.SELLER_TYPE_PRESELL {
		checkResult, errT := o.checkReserveTakeFoodTime(vendorInfo, TakeGoodsTime)
		if errT != nil {
			return "", common.NewError(common.ErrorTakeFoodTimeError)
		}
		if !checkResult {
			return "", common.NewError(common.ErrorTakeFoodTimeError)
		}
		//预售取货有效期
		availabilityFoodTime = int(TakeGoodsTime) + vendorInfo.AvailabilityGoodsTime
	}
	//fmt.Printf("===========availabilityFoodTime============%v \r\n", availabilityFoodTime)
	// ==================================================预售店检查时间 end
	//获取购物车
	NewRpcCartService := new(Rpc.RpcCartService)
	cartsInfo, errCart := NewRpcCartService.GetCartInfo(orderData.BuyerId, orderData.SellerId)
	if errCart != nil {
		fmt.Printf("===========cartsInfo Empty============%v \r\n", errCart)
		return "", errCart //购物车为空
	}
	//fmt.Printf("===========cartsInfo============%v",cartsInfo)
	var initMap = map[string]interface{}{
		"VendorInfo":    vendorInfo,
		"TakeGoodsTime": TakeGoodsTime,
		"BusinessTime":  businessTime,
		"Carts":         cartsInfo,
	}
	productsService := new(Rpc.RpcProductsService)
	productsService.Init(initMap)
	productInfo, orderFoodStock, errGetProducts := productsService.GetProductsInfo()
	if errGetProducts != nil {
		//fmt.Printf("=======GetProductsInfo=105===%v============%v\n\n\n", productInfo, orderFoodStock)
		return "", errGetProducts //购物车为空
	}
	fmt.Printf("=======GetProductsInfo=108===%v===========orderFoodStock=%v\n\n\n", productInfo, orderFoodStock)
	//优惠券相关start
	//======================================================
	//....................................
	//校验商品的价格和数量是否正确
	orderGrandTotal := orderData.OrderGrandTotal
	orderAmountPrice := productInfo["order_amount_total"].(float64)
	var couponFee float64 = 0
	couponCode := ""
	couponUserCode := ""
	if orderData.CouponUserCode != "" {
		cFee, cCode, cUserCode, err := new(Rpc.RpcCouponService).CheckIsOnCoupon(
			vendorInfo.SellerID,
			orderData.CouponUserCode,
			orderAmountPrice,
			orderData.BuyerId)
		if err != nil {
			return "", err
		}
		couponCode = cCode
		couponUserCode = cUserCode
		couponFee = cFee
		orderAmountPrice = orderAmountPrice - couponFee
	}

	if orderAmountPrice != orderGrandTotal {
		fmt.Printf("=======GetProductsInfo=140==caculatePrice=%v===========orderGrandTotal=%v\n\n\n", orderAmountPrice, orderGrandTotal)
		return "", common.NewError(common.ErrorProductTotalPrice) //購買的商品總價值有誤！
	}
	if orderAmountPrice < 4 {
		return "", common.NewError(common.ErrorProductTotalPriceLt4Error) //訂單金額不小於4元
	}
	if orderAmountPrice > 20000 {
		return "", common.NewError(common.ErrorPTotalPGt20000Error) //訂單金額不大于於20000元
	}
	if productInfo["order_product_counts"] != orderData.OrderProductCounts {
		fmt.Printf("\n\n\n=======GetProductsInfo=150==order_product_counts=%v===========OrderProductCounts=%v\n\n\n", productInfo["order_product_counts"], orderData.OrderProductCounts)
		return "", common.NewError(common.ErrorProductTotalCounts) //購買的商品總數量有誤！
	}
	//............优惠券等待处理............
	//======================================================
	//优惠券相关end

	//获取用户信息
	buyerInfo, errBuyerInfo := new(Rpc.RpcBuyerService).GetBuyerInfo(orderData.BuyerId)
	if errBuyerInfo != nil {
		fmt.Printf("===========buyerInfo Empty============%v\n\n\n", errBuyerInfo)
		return "", errBuyerInfo //买家信息为空
	}

	prepareFoodTime := int(TakeGoodsTime) - vendorInfo.PrepareGoodsTime*60 + 30

	//减掉商品库存
	var orderFoodStockF map[int]int
	orderFoodStockF = make(map[int]int)
	for k, v := range orderFoodStock {
		intK, _ := strconv.Atoi(k)
		intV, _ := strconv.Atoi(v.(string))
		orderFoodStockF[intK] = intV
	}
	err1 := productsService.ModifyProductsStock(orderFoodStockF, "minus", vendorInfo.SellerID) //action minus或者add
	if err1 != nil {
		return "", err1
	}

	//获取订单号
	orderId, _ := o.OrderCreateSuffixId(platformId, orderData.BuyerId)
	//生成orderModel数据
	orderProductCounts := productInfo["order_product_counts"].(int)
	orderAmountTotal := productInfo["order_amount_total"].(float64)

	ServiceRateFloat64, _ := strconv.ParseFloat(vendorInfo.ServiceRate, 64)

	var OrderModel model.Order
	nowTimeInt := int(time.Now().Unix())
	OrderModel.OrderId 								= orderId
	OrderModel.PlatformId 							= platformId
	OrderModel.BuyerId 								= buyerInfo.BuyerId
	OrderModel.BuyerName 							= buyerInfo.Nickname
	OrderModel.SellerId 							= vendorInfo.SellerID
	OrderModel.SellerName 							= vendorInfo.SellerName
	OrderModel.CellphoneNumber 						= buyerInfo.Phone
	OrderModel.TradeType							= 1
	OrderModel.SceneType							= 2
	OrderModel.SellerType							= vendorInfo.SellerType
	OrderModel.OrderStatus 							= Config.CCOrderStatus_UNPAID //ORDER_UNPAID
	OrderModel.TakeGoodsTime						= int(TakeGoodsTime)
	OrderModel.PrepareGoodsTime						= prepareFoodTime
	OrderModel.ServiceRate 							= ServiceRateFloat64
	OrderModel.OrderProductCounts 					= orderProductCounts //订单商品总数量（不包括配菜）
	OrderModel.OrderGrandTotal 						= orderGrandTotal
	OrderModel.OrderAmountTotal 					= orderAmountTotal             //订单商品总金额(此订单所有商品的金额)
	OrderModel.CreatedAt 							= nowTimeInt
	OrderModel.UpdatedAt 							= nowTimeInt
	OrderModel.CouponFee 							= couponFee
	OrderModel.CouponUserCode 						= couponUserCode
	OrderModel.CouponCode 							= couponCode
	OrderModel.ActualPayPrice 						= 0
	OrderModel.SellerType 							= vendorInfo.SellerType
	OrderModel.AvailabilityGoodsTime 				= availabilityFoodTime
	OrderModel.PackingFee 							= 0
	OrderModel.ShippingFee 							= 0


	productInfoItems, ok := productInfo["items"].([]map[string]interface{})
	//定义orderDetail数据结构
	var OrderDetailModel model.OrderDetail
	var OrderDetailModels []model.OrderDetail
	if ok && len(productInfoItems) > 0 {
		for _, productInfoItem := range productInfoItems {
			OrderDetailModel.OrderId 				= orderId
			OrderDetailModel.SellerId 				= vendorInfo.SellerID
			OrderDetailModel.ProductComment 		= productInfoItem["product_msg"].(string)
			OrderDetailModel.Sku 					= productInfoItem["product_id"].(string)
			OrderDetailModel.ParentSku 				= productInfoItem["product_parent_id"].(string)
			OrderDetailModel.ProductCount 			= productInfoItem["product_counts"].(int)
			productRealId, _ := strconv.Atoi(productInfoItem["product_real_id"].(string))
			OrderDetailModel.ProductId 				= productRealId
			OrderDetailModel.ProductName 			= productInfoItem["product_name"].(string)
			OrderDetailModel.ProductPriceCount		= productInfoItem["product_counts"].(int) //商品计价数量

			productItemPrice, _ := productInfoItem["product_price"].(json.Number).Float64()
			productUnitPrice1 := fmt.Sprintf("%.2f", productItemPrice)
			OrderDetailModel.ProductUnitPrice = productUnitPrice1 //商品销售单价

			productPrice1 := fmt.Sprintf("%.2f", float64(productInfoItem["product_counts"].(int))*productItemPrice)
			OrderDetailModel.ProductPrice = productPrice1 //商品总销售价

			productItemOriginalPrice, _ := productInfoItem["product_original_price"].(json.Number).Float64()
			productOriginalPrice1 := fmt.Sprintf("%.2f", productItemOriginalPrice)
			OrderDetailModel.ProductUnitOriginalPrice = productOriginalPrice1 //商品原始单价

			ProductOriginalPrice1 := fmt.Sprintf("%.2f", float64(productInfoItem["product_counts"].(int))*productItemOriginalPrice)
			OrderDetailModel.ProductOriginalPrice = ProductOriginalPrice1

			OrderDetailModel.ProductType = 1 //1实物 2虚拟
			OrderDetailModels = append(OrderDetailModels, OrderDetailModel)
		}
	}

	//写入订单索引表------
	var OrderIndexModel model.OrderIndex
	OrderIndexModel.OrderId 			= orderId
	OrderIndexModel.SellerId 			= vendorInfo.SellerID
	OrderIndexModel.SellerType 			= OrderModel.SellerType
	OrderIndexModel.BuyerId 			= buyerInfo.BuyerId
	OrderIndexModel.PlatformId 			= platformId
	OrderIndexModel.OrderStatus 		= Config.CCOrderStatus_UNPAID
	OrderIndexModel.CellphoneNumber 	= buyerInfo.Phone
	OrderIndexModel.TakeGoodsTime 		= int(TakeGoodsTime)
	OrderIndexModel.TakeGoodsCode 		= 0
	OrderIndexModel.OrderNumeral 		= 0
	OrderIndexModel.UpdatedAt 			= nowTimeInt
	OrderIndexModel.CreatedAt 			= nowTimeInt

	////orderSys清加事务///////////////////////////////////////////////////////////////////////////////////
	////写入订单表
	err = repositorie.NewDefaultOrderRepositories(OrderModel.OrderId).CreateOrder(OrderModel, OrderDetailModels, OrderIndexModel)
	if err != nil {
		//回滚库存
		productsService.ModifyProductsStock(orderFoodStockF, "add", vendorInfo.SellerID) //action minus或者add
		return "", err
	}

	//删除购物车数据
	errCartDel := NewRpcCartService.DeleteCart(buyerInfo.BuyerId, vendorInfo.SellerID)
	if errCartDel != nil {
		return "", errCartDel
	}
	////orderSys清加事务 结束/////////////////////////////////////////////
	return orderId, nil
}

// 检查商家运营情况
func (o *ProjectOrder) checkIsOnSales(vendorInfo project.RpcSeller, takeGoodTime int64) ([]map[string]string, error) {
	vendorStatusInfo, errVendor := new(Rpc.RpcSellerService).OperatingStatus(vendorInfo.SellerID, 0)

	if errVendor != nil {
		return nil, common.NewError(common.ErrorVendorNotExist)
	}
	if vendorStatusInfo.Status == 3 {
		return nil, common.NewError(common.ErrorVendorNotExist)
	}
	if vendorStatusInfo.Status != 1 {
		return nil, common.NewError(common.ErrorVendorNotOnSale) // 'merchants_isnot_on_sale' => '商家已打烊！',
	}
	businessTimeO := o.getBusinessTime(vendorStatusInfo.BusinessTime, vendorInfo.SellerType)
	if vendorInfo.SellerType == 2 {
		vendorStatusInfo2, errVendor2 := new(Rpc.RpcSellerService).OperatingStatus(vendorInfo.SellerID, takeGoodTime)
		if errVendor2 != nil {
			return nil, common.NewError(common.ErrorTakeFoodTimeF)
		}
		if vendorStatusInfo2.Status != 1 {
			return nil, common.NewError(common.ErrorVendorNotOnSaleTime)
		}
		businessTimeO = o.getBusinessTime(vendorStatusInfo2.BusinessTime, vendorInfo.SellerType)
	}
	return businessTimeO, nil
}

/**
校验预售店用户的取餐时间选择是否正确
 */
func (o *ProjectOrder) checkReserveTakeFoodTime(vendorInfo project.RpcSeller, takeGoodsTimeStamp int64) (bool, error) {

	if vendorInfo.BookTimeType == 0 {
		return false, nil
	}
	//可预订时间类型 1:固定时间范围 2:下单后N天
	if vendorInfo.BookTimeType == 1 { //检验可取貨時間
		//固定时间端
		tmpTimeArr := strings.Split(vendorInfo.BookTimeValue, "~")
		if len(tmpTimeArr) != 2 {
			return false, nil
		}
		startTime, errT := time.Parse("2006-01-02", tmpTimeArr[0])
		endTime, errT2 := time.Parse("2006-01-02", tmpTimeArr[1])
		if errT != nil || errT2 != nil {
			return false, errT
		}
		if startTime.Unix() < endTime.Unix() || takeGoodsTimeStamp > endTime.Unix() {
			return false, nil
		}
	}
	return true, nil
}

/**
     * 时间对比 返回状态
     *
     * @param int    status 营业状态
     * @param string businessTime 营业时间  //00:00~12:00,12:00~23:59 或者　00:00~12:00
     * @param int dateTime 指定某个时间的状态和营业时间
     */
func (o *ProjectOrder) compareTime(status int, businessTime string, dateTime int64) (int, string, error) {
	resultStatus := 2
	resultBusinessTime := businessTime

	timesArr := strings.Split(businessTime, ",")
	//fmt.Printf("\n\n\n\n=====compareTime===================timesArr=%v\n\n\n\n",timesArr[0])
	var now string
	if dateTime != 0 {
		goodTimeS := time.Unix(dateTime, 0) //dateTime格式为2006-01-02 15:04:05或者为空
		if goodTimeS.Hour() < 10 {
			if goodTimeS.Minute() < 10 {
				now = fmt.Sprintf("0%d:0%d", goodTimeS.Hour(), goodTimeS.Minute())
			} else {
				now = fmt.Sprintf("0%d:%d", goodTimeS.Hour(), goodTimeS.Minute())
			}

		} else {
			if goodTimeS.Minute() < 10 {
				now = fmt.Sprintf("%d:0%d", goodTimeS.Hour(), goodTimeS.Minute())
			} else {
				now = fmt.Sprintf("%d:%d", goodTimeS.Hour(), goodTimeS.Minute())
			}

		}

	} else {
		nowTime := time.Now()
		if nowTime.Hour() < 10 {
			if nowTime.Minute() < 10 {
				now = fmt.Sprintf("0%d:0%d", nowTime.Hour(), nowTime.Minute())
			} else {
				now = fmt.Sprintf("0%d:%d", nowTime.Hour(), nowTime.Minute())
			}
		} else {
			if nowTime.Minute() < 10 {
				now = fmt.Sprintf("%d:0%d", nowTime.Hour(), nowTime.Minute())
			} else {
				now = fmt.Sprintf("%d:%d", nowTime.Hour(), nowTime.Minute())
			}
		}

	}
	for i := 0; i < len(timesArr); i++ {
		v := timesArr[i]
		temp := strings.Split(v, "~")
		//fmt.Printf("\n\n\n\n=====compareTime===================temp=%v\n\n\n\n",temp)
		st := temp[0]
		end := temp[1]
		if end > st {
			//fmt.Printf("\n\n\n\n=====compareTime===================end > st===now=%v===vendor=%v\n\n\n\n",now,end)
			if now >= st && now <= end {
				return status, resultBusinessTime, nil
			}
		} else {
			if st >= now {
				if end >= now {
					return status, resultBusinessTime, nil
				}
			} else {
				// 开始比现在早
				// 15:00 ~ 13:00
				// now  14:00
				tmpEnd := strings.Split(end, ":")
				tmpEnd0, _ := strconv.Atoi(tmpEnd[0])
				tmpEnd[0] = strconv.Itoa(tmpEnd0 + 24)
				end = strings.Join(tmpEnd, ":")
				if now <= end {
					return status, resultBusinessTime, nil
				}
			}
		}
	}
	return resultStatus, resultBusinessTime, nil
}

/**
商家类型 1:外卖 2:预售
 */
func (o *ProjectOrder) getBusinessTime(businessTime string, sellerType int) []map[string]string {
	period := strings.Split(businessTime, ",")
	var data []map[string]string
	for i := 0; i < len(period); i++ {
		item := period[i]
		splitTime := strings.Split(item, "~")
		//fmt.Printf("======getBusinessTime===%v=====================",splitTime)
		var tmp map[string]string
		tmp = make(map[string]string)
		if splitTime[0] != "" {
			tmp["start_time"] = splitTime[0]
		} else {
			tmp["start_time"] = "0"
		}
		if splitTime[1] != "" {
			tmp["end_time"] = splitTime[1]
		} else {
			tmp["end_time"] = "0"
		}
		if sellerType == 1 {
			if len(data) > i {
				tmpStartTime, _ := time.Parse("15:04", tmp["start_time"])
				tmpEndTime, _ := time.Parse("15:04", data[i-1]["end_time"])
				//e肚仔要求40分钟
				if tmpStartTime.Unix()-tmpEndTime.Unix() < 2400 {
					data[i-1] = map[string]string{
						"start_time": data[i-1]["start_time"],
						"end_time":   tmp["end_time"],
					}
				}
			} else {
				data = append(data, tmp)
			}
		} else {
			data = append(data, tmp)
		}

	}
	//fmt.Printf("\n\n\n\n\n======getBusinessTime==data=%v======sellerType=%v==============\n\n\n\n\n",data,sellerType)
	return data
}

func (o *ProjectOrder) PaySuccess(params Form.PaySuccessData) (model.Order, error) {
	orderRepository := repositorie.NewDefaultOrderRepositories(params.OrderId)
	orderData, err := orderRepository.OrderGetOneByOid(params.OrderId)
	if err != nil {
		return orderData, common.NewError(common.ErrorOrderDoesNotExist)
	}

	if orderData.OrderGrandTotal != params.ActualPayPrice {
		return orderData, common.NewError(common.ErrorOrderPayPriceErr)
	}

	orderData.OrderStatus = Config.CCOrderStatus_NEED_ACCEPT
	if orderData.SellerType == Config.SELLER_TYPE_TAKE_OUT {
		orderData.OrderStatus = Config.CCOrderStatus_PREPARE_GOODS
	}

	orderNumeral, err := o.generateOrderNumeral(orderData.SellerId, orderData.SellerType)
	if err != nil {
		return orderData, err
	}

	now := time.Now().Unix()
	orderData.PaySerialNumber	= params.PaySerialNumber
	orderData.ActualPayPrice 	= params.ActualPayPrice
	orderData.PayType			= params.PayType
	orderData.UpdatedAt			= int(now)
	orderData.OrderPayTime		= int(now)
	orderData.OrderNumeral		= orderNumeral
	orderData.TakeGoodsCode		= o.generateTakeFoodCode()
	orderData.OrderStatus 		= Config.CCOrderStatus_NEED_ACCEPT
	if orderData.SellerType == Config.SELLER_TYPE_PRESELL {
		orderData.OrderStatus = Config.CCOrderStatus_PREPARE_GOODS
	}
	orderColumns := []string{
		"PayType",
		"ActualPayPrice",
		"OrderPayTime",
		"PaySerialNumber",
		"TakeGoodsCode",
		"OrderStatus",
		"OrderNumeral",
		"UpdatedAt"}

	var orderIndex model.OrderIndex
	orderIndex.OrderStatus 		= orderData.OrderStatus
	orderIndex.UpdatedAt		= orderData.UpdatedAt
	orderIndex.OrderNumeral		= orderData.OrderNumeral
	orderIndex.TakeGoodsCode	= orderData.TakeGoodsCode
	orderIndexColums := []string{"OrderStatus", "UpdatedAt", "OrderNumeral", "TakeGoodsCode"}

	where := map[string]interface{}{
		"order_id" : params.OrderId,
		"order_status" : Config.CCOrderStatus_UNPAID,
	}
	errUp := orderRepository.UpdateOrder(where, orderData, orderColumns, orderIndex, orderIndexColums)
	if errUp != nil {
		return orderData, errUp
	}

	//推送
	newOrderType := Config.ORDER_ACCEPT_NEW
	messageType := 0
	if orderData.SellerType == Config.SELLER_TYPE_TAKE_OUT {
		prepareFoodTime := orderData.PrepareGoodsTime -
			Config.DEFAULT_ADVANCE_ORDER_UNIT_TIME*60 +
			Config.DEFAULT_ADVANCE_ORDER_DELAY_TIME
		if orderData.OrderPayTime < prepareFoodTime {
			newOrderType = Config.ORDER_ACCEPT_ADVANCE
		}

		messageType = Config.SELLER_PUSH_NEWORDER
	} else {
		messageType = Config.SELLER_PUSH_PRINT
	}
	new(Rpc.RpcPushService).
		PushMessage(orderData.OrderId, orderData.SellerId, messageType, newOrderType)

	return orderData, nil
}

//生成排单号
func (o *ProjectOrder) generateOrderNumeral(sellerId, sellerType int) (int, error) {
	timeStr := ""
	if sellerType == Config.SELLER_TYPE_TAKE_OUT {
		timeStr = strings.Join([]string{":", time.Now().Format("2006-01-02")}, "")
	}

	key := strings.Join(
		[]string{
			RedisSellerOrderNumeral,
			strconv.Itoa(sellerType),
			":",
			strconv.Itoa(sellerId),
			timeStr},
		"")
	fmt.Println(key)
	redisConn := common.RedisPool.Get()
	orderNumeral, err := redisConn.Do("incr", key)
	if err != nil {
		fmt.Println("===", err)
		return 0, common.NewError(common.ErrorRedisErr)
	}

	if sellerType == Config.SELLER_TYPE_TAKE_OUT {
		_, err := redisConn.Do("expire", key, 86400)
		if err != nil {
			return 0, common.NewError(common.ErrorRedisErr)
		}
	}

	return int(orderNumeral.(int64)), nil
}

//生产二维码
func (o *ProjectOrder) generateTakeFoodCode() int {
	rand.Seed(time.Now().UnixNano())
	takeFoodCodeStr := fmt.Sprintf("%08d", rand.Intn(99999999))
	takeFoodCode, _ := strconv.Atoi(takeFoodCodeStr)
	return takeFoodCode
}

