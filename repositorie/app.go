package repositorie

import (
	"order-backend/common"
	"order-backend/model"
)

func AppGetInfoByKey(appKey string) (model.App, error) {
	app := model.App{}
	_, err := common.DB.Where("app_key = ?", appKey).And("status=1").Get(&app)
	return app, err
}
func AppGetInfo(appKey string, appSecret string) (model.App, error) {
	app := model.App{}
	_, err := common.DB.Where("app_key = ? AND app_secret = ?", appKey, appSecret).Get(&app)
	return app, err
}
