package helpers

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
