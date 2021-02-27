package webhooks

func isStringInSlice(item string, list []string) bool {
	for _, listItem := range list {
		if listItem == item {
			return true
		}
	}

	return false
}
