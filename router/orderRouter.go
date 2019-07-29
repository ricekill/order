package router

import (
	"order-backend/action"
	actionApi "order-backend/action/api"
	"fmt"
	"github.com/hprose/hprose-golang/rpc"
)

func TcpRouter(server *rpc.TCPServer) {
	fmt.Println("======== router start ===========")
	paths := []string{"orderCreate", "paySuccess", "acceptOrder", "takeGoodByCode", "confirmTakeGood"}
	methods := []interface{}{
		actionApi.OrderCreate,
		actionApi.PaySuccess,
		actionApi.AcceptOrder,
		actionApi.TakeGoodByCode,
		actionApi.ConfirmTakeGood,
	}
	server.AddFunctions(paths, methods).AddInvokeHandler(action.InvokeHandler).AddFilter(LogFilter{"server"})
}

type LogFilter struct {
	Prompt string `default:"server"`
}

func (lf LogFilter) InputFilter(data []byte, context rpc.Context) []byte {
	fmt.Printf("InputFilter == %v: %s\r\n", lf.Prompt, data)
	return data
}

func (lf LogFilter) OutputFilter(data []byte, context rpc.Context) []byte {
	fmt.Printf("OutputFilter == %v: %s\r\n", lf.Prompt, data)
	return data
}