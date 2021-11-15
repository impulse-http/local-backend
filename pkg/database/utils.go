package database

func boolToInt(val bool) int {
	if val {
		return 1
	}
	return 0
}

func intToBool(val int) bool {
	return val != 0
}
