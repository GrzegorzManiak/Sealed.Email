package helpers

import "strings"

func FormatOrderString(order string) string {
	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		return "DESC"
	}
	return order
}
