package action

import (
	"order-backend/common"
	"order-backend/service"
	"errors"
	"encoding/json"
	"fmt"
	"github.com/hprose/hprose-golang/rpc"
	"reflect"
)

type Auth struct {
	AppId 		int 	`json:"app_id"`
	AppKey		string	`json:"app_key"`
	AppSecret	string	`json:"app_secret"`
}

func InvokeHandler(
	name string,
	args []reflect.Value,
	context rpc.Context,
	next rpc.NextInvokeHandler) (results []reflect.Value, err error) {
		//数据校验
		v := args[1].Elem().String()
		fmt.Printf("============== auth string ================ %v \r\n", v)

		var auth Auth
		json.Unmarshal([]byte(v), &auth)
		fmt.Println(auth)
		if auth.AppKey == "" || auth.AppSecret == "" {
			err = errors.New(
				common.RenderRpcJson("", "auth params error", common.ErrorUnauthorized, false))
			return args, err
		}

		if common.APP.AppKey != auth.AppKey {
			app, err1 := service.FindAppInfo(auth.AppKey, auth.AppSecret)
			common.APP = app
			fmt.Println("============== app auth ================")
			if err1 != nil || common.APP.AppKey != auth.AppKey {
				err = errors.New(
					common.RenderRpcJson("", "auth check failure", common.ErrorAuthSessionKeyFalse, false))
				return args, err
			}
		}

		//校验
		results, err = next(name, args, context)
		return
}