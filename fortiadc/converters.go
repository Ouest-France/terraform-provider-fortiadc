package fortiadc

// enableToBool convert "enable"/"disable" strings to equivalent boolean
func enableToBool(status string) bool {
	return status == "enable"
}

// boolToEnable convert boolean to equivalent "enable"/"disable" strings
func boolToEnable(status bool) string {
	if status {
		return "enable"
	}
	return "disable"
}
