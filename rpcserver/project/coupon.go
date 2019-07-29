package project

import (
	"order-backend/common"
	"order-backend/common/itface/remote"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hprose/hprose-golang/rpc"
	"order-backend/model/project"
)

type RpcCouponService struct {
	RpcClient rpc.Client
}

func (reus *RpcCouponService) getCouponRpcClient() remote.RpcCouponService {
	var couponService remote.RpcCouponService
	reus.RpcClient = rpc.NewClient(common.Config.RpcServer.Project.Coupon)
	reus.RpcClient.UseService(&couponService)
	return couponService
}

func (reus *RpcCouponService) CheckIsOnCoupon(sellerId int, couponUserCode string, orderAmountTotal float64, buyerId int) (
	cFee float64,
	cCode string,
	cUserCode string,
	err error) {
	couponMap := map[string]interface{}{
		"status":  "no",
		"data":    "",
		"message": "",
	}
	//从eaojoy获取购物车信息,返回json字符串
	couponInfo := reus.getCouponRpcClient().Coupon_checkIsOnCoupon(sellerId, couponUserCode, orderAmountTotal, buyerId)
	fmt.Println("\n\n==============couponInfo======rpc=======", couponInfo)
	defer reus.RpcClient.Close()
	json.Unmarshal([]byte(couponInfo), &couponMap)
	fmt.Printf("\n\n\n\n=================couponInfo=========%f\n\n\n\n", couponMap)
	//如果message不为空,返回错误信息
	if couponMap["status"] == "no" && couponMap["message"] != "" {
		return 0, "", "", common.NewError(common.ErrorCouponInvalid)
	} else if couponMap["status"] == "ok" && couponMap["message"] == "" {
		couponData := couponMap["data"].(map[string]interface{})
		couponsFee := couponData["coupons_fee"].(float64)
		couponCode := couponData["coupon_code"].(string)
		couponUserCode := couponData["coupon_user_code"].(string)
		return couponsFee, couponCode, couponUserCode, nil
	}
	return 0, "", "", nil
}

func (reus *RpcCouponService) ReturnCouponCode(buyerId int, couponUserCode string) error {
	returnJsonReturn := reus.getCouponRpcClient().Coupon_returnCouponCode(buyerId, couponUserCode)
	defer reus.RpcClient.Close()
	returnJsonMap := map[string]interface{}{
		"status":  "no",
		"data":    "",
		"message": "",
	}
	errJsonParse := json.Unmarshal([]byte(returnJsonReturn), &returnJsonMap)
	if errJsonParse != nil {
		return errJsonParse
	}
	if returnJsonMap["status"] == "no" && returnJsonMap["message"] != "" {
		return errors.New(returnJsonMap["message"].(string))
	} else if returnJsonMap["status"] == "ok" && returnJsonMap["message"] == "" {
		return nil
	}
	return nil
}

func (reus *RpcCouponService) DoUseCoupon(buyerId int, couponUserCode string) error {
	returnJsonReturn := reus.getCouponRpcClient().Coupon_doUseCoupon(couponUserCode, buyerId)
	defer reus.RpcClient.Close()
	var returnJsonMap project.RpcReturn
	errJsonParse := json.Unmarshal([]byte(returnJsonReturn), &returnJsonMap)
	if errJsonParse != nil || !returnJsonMap.Status {
		return common.NewError(common.ErrorDoUseCouponError)
	}
	return nil
}
