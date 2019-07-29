package middleware

import (
	"crypto/sha1"
	"order-backend/common"
	"order-backend/model"
	"order-backend/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const HEADER_CLIENTKEY  = "X-ClientKey"
const HEADER_AUTH  ="Authorization"

func Auth() gin.HandlerFunc  {
	return func(c *gin.Context) {
		verifySignature(c)
		c.Next()
	}
}

func getRequestClient(c *gin.Context)(model.App,error){
	appKey := c.Request.Header.Get(HEADER_CLIENTKEY)
	fmt.Printf("X-ClientKey============%v\n",appKey)
	appInfo, err :=service.FindAppInfoByKey(appKey)
	return appInfo,err
}

func verifySignature(c *gin.Context){
	appInfo,err := getRequestClient(c)
	fmt.Printf("appinfo============%v\n",appInfo)
	if err!=nil || appInfo.Id<=0  {
		common.RenderJSONWithError(c, common.NewError(common.ErrorUnauthorized), http.StatusUnauthorized)
	} else {
		secret:=appInfo.AppSecret
		if !verifyParamsSignature(c, secret) {
			common.RenderJSONWithError(c, common.NewError(common.ErrorUnauthorized), http.StatusUnauthorized)
		}
	}
	//将Appinfo信息存入全局变量
	common.APP.Id=appInfo.Id
	common.APP.AppSecret=appInfo.AppSecret
	common.APP.AppKey=appInfo.AppKey
	common.APP.Code=appInfo.Code
}
func verifyParamsSignature(c *gin.Context, secret string) bool {
	c.Request.ParseForm()
	formData:=c.Request.Form
	paramsMaps:=make(map[string]string)
	signature:=formData.Get("signature")
	//fmt.Printf("formData:%v",formData)
	for k, v := range formData {
		if k=="signature" {continue}
		if strings.Contains(k,"[") {
			start:=strings.IndexAny(k,"[")
			k=common.Substr(k,0,start)
			paramsMaps[k]=""
		} else {
			paramsMaps[k]=v[0]
		}
	}
	expectedSignature:=getParamsSignature(paramsMaps,secret)
	return signature==expectedSignature
}

func getParamsSignature(params map[string]string, secret string) string{
   keys :=make([]string, len(params))
   i := 0
   for k,_:=range params {
   	keys[i] = k
   	i++
   }
   sort.Strings(keys)
	var returnString string
	for _, k := range keys {
		kString:=url.QueryEscape(k)
		vString:=url.QueryEscape(params[k])
		kString=strings.Replace(kString,"+","%20",-1)
		vString=strings.Replace(vString,"+","%20",-1)
		returnString+=fmt.Sprintf("%s=%s",kString,vString)
	}
	returnString=fmt.Sprintf("%s.%s",returnString,secret)
	//fmt.Printf("=============%v",returnString)
	byteStr := []byte(returnString)
	returnSha1:=sha1.New()
	returnSha1.Write(byteStr)
	result :=fmt.Sprintf("%x",returnSha1.Sum(nil))
	return result

}