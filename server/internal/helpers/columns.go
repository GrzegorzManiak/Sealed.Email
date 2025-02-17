package helpers

import "strings"

func BuildColumnList(columns []string) string {
	columnList := ""
	for i, column := range columns {
		if i == 0 {
			columnList += column
		} else {
			columnList += "," + column
		}
	}
	return columnList
}

func FormatOrderString(order string) string {
	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		return "DESC"
	}
	return order
}
