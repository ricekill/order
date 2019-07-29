package project

import (
	"order-backend/common"
	"order-backend/common/itface/remote"
	"encoding/json"
	"github.com/hprose/hprose-golang/rpc"
	"order-backend/model/project"
)

type RpcSellerService struct {
	RpcClient rpc.Client
}

func (reus *RpcSellerService) getSellerRpcClient() remote.RpcSellerService {
	var sellerService remote.RpcSellerService
	reus.RpcClient = rpc.NewClient(common.Config.RpcServer.Project.Vendor)
	reus.RpcClient.UseService(&sellerService)
	return sellerService
}

func (reus *RpcSellerService) gV(vendorId int) string {
	vendorInfo := reus.getSellerRpcClient().Vendor_getInfo(vendorId)
	return vendorInfo
}

func (reus *RpcSellerService) GetVendorInfo(vendorId int) (project.RpcSeller, error) {
	//defer func() {
	//	if err := recover(); nil != err {
	//		common.Log.Errorln("RPC IP or port link failed")
	//		fmt.Println("RPC IP or port link failed")
	//		ErrBool = true
	//	}
	//}()

	//定义返回信息数据结构struct
	var rpcVendorInfo project.RpcSeller
	//从eaojoy获取购物车信息,返回json字符串
	vendorInfo := reus.gV(vendorId)
	err := json.Unmarshal([]byte(vendorInfo), &rpcVendorInfo)
	if rpcVendorInfo.SellerID <= 0 || err != nil {
		errRpc := common.NewError(common.ErrorVendorEmpty)
		return rpcVendorInfo, errRpc
	}
	defer reus.RpcClient.Close()
	return rpcVendorInfo, nil
}

func (reus *RpcSellerService) OperatingStatus(vendorId int, takeGoodTime int64) (project.RpcSellerOperatingStatus, error) {
	//从eaojoy获取购物车信息,返回json字符串
	vendorInfo := reus.getSellerRpcClient().Vendor_getOperatingStatus(vendorId, takeGoodTime)

	//定义返回信息数据结构struct
	var rpcVendorInfo project.RpcSellerOperatingStatus
	err := json.Unmarshal([]byte(vendorInfo), &rpcVendorInfo)
	if rpcVendorInfo.Status <= 0 || err != nil {
		errRpc := common.NewError(common.ErrorVendorEmpty)
		return rpcVendorInfo, errRpc
	}
	defer reus.RpcClient.Close()
	return rpcVendorInfo, nil
}

func (reus *RpcSellerService) GetMultiInfo(vendorIdS []int) (project.RpcSellers, error) {
	//从eaojoy获取购物车信息,返回json字符串
	vendorInfos := reus.getSellerRpcClient().Vendor_getMultiInfo(vendorIdS)
	//定义返回信息数据结构struct
	var rpcVendorInfos project.RpcSellers
	err := json.Unmarshal([]byte(vendorInfos), &rpcVendorInfos)
	if len(rpcVendorInfos) <= 0 || err != nil {
		errRpc := common.NewError(common.ErrorVendorEmpty)
		return nil, errRpc
	}
	defer reus.RpcClient.Close()
	return rpcVendorInfos, nil
}

//商家计算评分
func (reus *RpcSellerService) CalculateSellerScore(sellerId int) error {
	//从eaojoy获取购物车信息,返回json字符串
	vendorInfos := reus.getSellerRpcClient().Vendor_vendorScore(sellerId)
	//定义返回信息数据结构struct
	var rpcVendorReturn project.RpcReturn
	err := json.Unmarshal([]byte(vendorInfos), &rpcVendorReturn)
	if err != nil {
		return err
	}
	if !rpcVendorReturn.Status {
		return common.NewError(common.ErrorUpdateDataFail)
	}
	defer reus.RpcClient.Close()
	return nil
}

