package gocd

func GetSLice(values any) []string {
	vals := values.([]any)

	newValue := make([]string, 0, len(vals))

	for _, value := range values.([]any) {
		newValue = append(newValue, value.(string))
	}

	return newValue
}
