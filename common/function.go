package common

import "time"

func DataType(i interface{}) string {  //函数t有一个参数i
	switch i.(type) { //多选语句switch
		case string:
			return "string"
		case int:
			return "int"
		case int8:
			return "int8"
		case []string:
			return "[]string"
		case []int:
			return "[]int"
		default:
			return ""
	}
}

//获取当天凌晨的时间戳
func GetDayZeroTime() int {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	timeNumber := t.Unix()
	return int(timeNumber)
}
