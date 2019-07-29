package api

import (
	"order-backend/common"
	"order-backend/service"
	"github.com/gin-gonic/gin"
)

func HomePage(c *gin.Context) error {
	appKey :=c.Params.ByName("app_key")//c.Query("app_key")
	appSecret :=c.Params.ByName("app_secret")

	appInfo, err := service.FindAppInfo(appKey,appSecret)
	if err != nil {
		return err
	}
	//返回结果
	common.RenderJSON(c, gin.H{
		"id":appInfo.Id,
		"app_key":appInfo.AppKey,
		"app_secret":appInfo.AppSecret,
	})
	return nil
}
