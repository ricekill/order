package remote

type RpcProductsService struct {
	Products_getDetails func([]int, int) string//按用户id,商家id获取购物车
	Products_modifyStock func(map[int]int,string,int) string//减扣库存
	Products_updateCounts func(map[int]int, int) string //完成订单 统计商品总数
}

type ProductsService interface {
	Init(map[string]interface{})
	GetProductsInfo()  (map[string]interface{},map[string]interface{},error) //按用户id,商家id获取购物车
	ModifyProductsStock(map[int]int,string,int)error	//减扣库存
	UpdateProductCount(map[int]int, int) error //完成订单 统计商品总数
}
