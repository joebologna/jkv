package pkg

func BoolToString(isTrue bool) string {
	if isTrue {
		return "1"
	}
	return "0"
}

func StringToBool(isTrue string) bool {
	return isTrue == "1"
}
