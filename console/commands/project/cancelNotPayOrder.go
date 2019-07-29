package project

import (
	"order-backend/common"
	"order-backend/service/project"
	"sync"
)
////////实现单例/////////////////
var instanceOrder *project.ProjectOrder
var onceOrder sync.Once

func GetOrderInstance() *project.ProjectOrder {
	onceOrder.Do(func() {
		instanceOrder = new(project.ProjectOrder)
	})
	return instanceOrder
}
///////////////单列实现结束///////////////////////

//待支付的订单5分钟之内没付款则改为已取消
func CancelNotPayOrder(){
	//fmt.Printf("待支付的订单5分钟之内没付款则改为已取消%v\n",time.Now().Unix())
	errCON:=GetOrderInstance().OrderCancelNotPay()
	if errCON != nil {
		common.Log.Errorf("[CRON] CancelNotPayOrder Error:%v",errCON)
	}
	//else {
	//	common.Log.Infof("[CRON] CancelNotPayOrder Success at:%v",time.Now())
	//}
}
