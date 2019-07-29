package main

import (
	"github.com/jasonlvhit/gocron"
	"order-backend/common"
	"order-backend/console/commands/project"
)

func main()  {
	//初始化日志
	common.CheckErr(common.LoadConfig())
	common.CheckErr(common.SetupLogger())
	common.CheckErr(common.OpenDb())
	//定义计划任务
	//
	//待支付的订单5分钟之内没付款则改为已取消
	gocron.Every(1).Second().Do(project.CancelNotPayOrder)
	//
	////待接单的下单时间大于等于10分钟没接单更改为已取消
	//gocron.Every(1).Second().Do(project.CancelWithPayOrder)//
	//
	////待接单，下单时间大于等于3分钟没接单邮件提醒
	//gocron.Every(1).Second().Do(project.CancelOrderWithPayEmailPush)//
	//
	////备餐中的订单,且当前时间大于等于取餐时间的更改为待取餐
	//gocron.Every(1).Second().Do(project.ReadyTakeFood)//
	//
	////待取餐订单超过15分钟没确认收货则推送消息
	//gocron.Every(1).Second().Do(project.ReadyPushWithOutConfirm)//
	//
	////取餐時間前五分鐘推送提醒,若距取餐時間不足5分钟则不提醒
	//gocron.Every(1).Second().Do(project.BeforeTakeFoodPush)//
	//
	////待备餐订单,时间到则更改状态
	//gocron.Every(1).Second().Do(project.AdvanceOrderToNewOrder)//

	//清楚订单临时表的数据 每天凌晨5点
	//gocron.Every(1).Day().At("5:00").Do(project.ClearOrderIndex)

	gocron.RunAll()//所有的任务先执行一次
	//启动计划任务
	<- gocron.Start()

}
