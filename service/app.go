package service

import (
	"order-backend/model"
	"order-backend/repositorie"
)

func FindAppInfo(appKey string, appSecret string)(model.App, error){
	app, err :=repositorie.AppGetInfo(appKey, appSecret)
	return app, err
}

func FindAppInfoByKey(appKey string)(model.App, error) {
	if appKey == "" {
		return model.App{},nil
	}
	app,err:=repositorie.AppGetInfoByKey(appKey)
	return app, err
}

