package project

type RpcBuyer struct {
	BuyerId         int      `json:"user_id"`
	Nickname        string   `json:"user_nickname"`
	Phone           string   `json:"cellphone_number"`
	OpenUserId      int      `json:"open_user_id"`
}
