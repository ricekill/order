package model

type OrderDetail struct {
	Id                   		int    		`xorm:"not null pk autoincr INT(10)"`
	SellerId					int 		`xorm:"comment('商家ID') INT(10)"`
	OrderId              		string 		`xorm:"not null comment('订单id') index CHAR(25)"`
	ParentSku            		string 		`xorm:"comment('父级SKU') VARCHAR(25)"`
	Sku                  		string 		`xorm:"comment('SKU') VARCHAR(25)"`
	ProductId            		int    		`xorm:"comment('商品ID') INT(10)"`
	ProductName          		string 		`xorm:"comment('商品名称') VARCHAR(25)"`
	ProductCount         		int    		`xorm:"comment('商品数量') INT(10)"`
	ProductPrice         		string 		`xorm:"comment('商品总售价') DECIMAL(15,2)"`
	ProductOriginalPrice 		string 		`xorm:"comment('商品总原价') DECIMAL(15,2)"`
	ProductType          		int    		`xorm:"comment('商品分类：1实物 2虚拟') TINYINT(3)"`
	ProductUnitPrice     		string 		`xorm:"comment('商品单价') DECIMAL(15,2)"`
	ProductUnitOriginalPrice	string 		`xorm:"comment('商品原始单价') DECIMAL(15,2)"`
	ProductPriceCount    		int    		`xorm:"comment('商品计价数量') TINYINT(3)"`
	ProductComment       		string 		`xorm:"comment('商品备注') VARCHAR(200)"`
}

