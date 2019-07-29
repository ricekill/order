package model

type App struct {
	Id        int    `xorm:"not null pk autoincr INT(11)"`
	Code      string `xorm:"not null VARCHAR(45)"`
	AppKey    string `xorm:"not null VARCHAR(45)"`
	AppSecret string `xorm:"not null VARCHAR(45)"`
	Status    int    `xorm:"not null INT(11) default 1"`
}
