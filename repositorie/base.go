package repositorie

import (
	"github.com/go-xorm/xorm"
)

func GetWhere(session *xorm.Session, where map[string]interface{}) *xorm.Session {
	for key, val := range where {
		switch val.(type) { //多选语句switch
		case string:
			session.Where(key + " = ?", val)
		case int:
			session.Where(key + " = ?", val)
		case int64:
			session.Where(key + " = ?", val)
		case interface{}:
			if v2, ok := val.(map[string]interface{}); ok == true {
				switch v2["m"] {
				case "in":
					session.In(key, v2["v"])
				case "like":
					session.Where(key + " like '%?%'", v2["v"])
				case ">":
					session.Where(key + " > ?", v2["v"])
				case ">=":
					session.Where(key + " >= ?", v2["v"])
				case "<":
					session.Where(key + " < ?", v2["v"])
				case "<=":
					session.Where(key + " <= ?", v2["v"])
				}
			}
		}
	}

	return session
}
