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

type RpcCartService struct {
	RpcClient rpc.Client
}
func (reus *RpcCartService) getCartRpcClient() remote.RpcCartService {
	var  cartService remote.RpcCartService
	reus.RpcClient=rpc.NewClient(common.Config.RpcServer.Project.Cart)
	reus.RpcClient.UseService(&cartService)
	return cartService
}

func (reus *RpcCartService) GetCartInfo(userid int ,vendorId int) (project.RpcCart,error) {
	//从eaojoy获取购物车信息,返回json字符串
	cartInfo:=reus.getCartRpcClient().Cart_getCart(userid, vendorId)
	defer reus.RpcClient.Close()
	//定义返回信息数据结构struct
	var rpcCartInfo project.RpcCart
	//将返回数据绑定到该数据结构上---这种用法很牛,不用手动一个个遍历
	//json生成struct 有一个小工具　https://mholt.github.io/json-to-go/
	//fmt.Printf("cartInfo==================%v",cartInfo)
	err := json.Unmarshal([]byte(cartInfo), &rpcCartInfo)
	//fmt.Printf("rpcCartInfo==================%v",len(rpcCartInfo))
	if len(rpcCartInfo)<=0 || err!=nil {
		fmt.Printf("rpcCartInfo==================%v",err)
		errRpc :=common.NewError(common.ErrorCartEmpty)
		return nil,errRpc
	}
	return rpcCartInfo,nil
}

func (reus *RpcCartService) DeleteCart(buyerId int ,sellerId int) (error) {
	jsonReturn:=reus.getCartRpcClient().Cart_delCart(buyerId,sellerId)
	defer reus.RpcClient.Close()
	returnJsonMap := map[string]interface{}{
		"status":"no",
		"code":0,
		"data":"",
		"message":"",
	}
	errJsonParse := json.Unmarshal([]byte(jsonReturn), &returnJsonMap)
	if errJsonParse != nil {
		return errJsonParse
	}
	fmt.Printf("\n\n\n\n=================couponInfo=========%f\n\n\n\n",returnJsonMap)
	//如果message不为空,返回错误信息
	if returnJsonMap["status"]=="no" && returnJsonMap["message"] != "" {
		return errors.New(returnJsonMap["message"].(string))
	} else if returnJsonMap["status"]=="ok" && returnJsonMap["message"] == ""{
		return nil
	}
	return nil
}
