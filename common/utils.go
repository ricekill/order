package common

import (
	"bytes"
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	AppContentJSON = "application/json; charset=utf-8"
	XForwardedFor  = "X-Forwarded-For"
	XRealIP        = "X-Real-IP"
)
type ErrorRes struct {
	Code int `json:"code"`
	Desc interface{}
}
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func GetMallId(c *gin.Context) int {
	mallId, err := strconv.Atoi(c.Params.ByName("mall_id"))
	if err != nil {
		mallId = 0
	}
	return mallId
}
func RenderJson(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, obj)
}

func RenderJSONWithError(c *gin.Context, err error,status ...int) {
	code := 0
	s := http.StatusInternalServerError
	cerr, ok :=err.(BadRequestError)
	if ok {
		code = cerr.Code()
		s = http.StatusBadRequest
	}
	if len(status) > 0 {
		s = status[0]
	}

	c.AbortWithStatusJSON(s, gin.H{
		"status" : s,
		"code"   : code,
		"message": err.Error(),
	})
}

func HttpPost(sUrl string, p interface{}) ([]byte, error) {
	bytesData, err :=json.Marshal(p)
	if err != nil {
		Log.Errorln(err)
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	resp, err := http.Post(sUrl,"application/json", reader)
	if err != nil {
		Log.Errorln(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewError(resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func Substr(str string, start int,end int) string {
	rs := []rune(str)
	length:=len(rs)
	if start <0 || start > length {
		panic("start is wrong")
	}
	if end<0 || end >length {
		panic("end is wrong")
	}
	return string(rs[start:end])
}

func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < min {
		RandInt64(min, max)
	}
	return i.Int64()
}

func GetRequestRealip(c *gin.Context) string {
	if c==nil {
		return "127.0.0.1"
	}

	var s_remoteAddr string
	if x_HeaderIP :=c.GetHeader(XForwardedFor);x_HeaderIP !="" {
		s_remoteAddr = x_HeaderIP
	} else if x_HeaderIP =c.GetHeader(XRealIP);x_HeaderIP !="" {
		s_remoteAddr = x_HeaderIP
	} else {
		s_remoteAddr, _, _ = net.SplitHostPort(c.Request.RemoteAddr)
		if s_remoteAddr == "::1" {
			s_remoteAddr = "127.0.0.1"
		}
	}
	//从ctx.Request.Header得到的 s_remoteAddr 可能出现ip,ip,ip
	ips := strings.Split(s_remoteAddr, ",")
	return ips[0]

}

func RenderJSON(c *gin.Context, obj interface{})  {
	c.JSON(http.StatusOK, obj)
}

func RenderRpcJson(data interface{}, message string, code int, status bool) string {
	result := make(map[string]interface{})
	result["status"] 	= status
	result["data"] 		= data
	result["message"] 	= message
	result["code"] 		= code

	j, err := json.Marshal(result)

	if err != nil {
		Log.Errorln("[JSON]","Json parsing error")
	}

	return string(j)
}

func RenderRpcJSONError(err error) string {
	err1 :=err.(BadRequestError)
	result := make(map[string]interface{})
	result["status"] 	= false
	result["data"] 		= ""
	result["message"] 	= err1.Error()
	result["code"] 		= err1.Code()

	j, e := json.Marshal(result)

	if e != nil {
		Log.Errorln("[JSON]","Json parsing error")
	}

	return string(j)
}
