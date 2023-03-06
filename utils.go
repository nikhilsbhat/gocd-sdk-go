package gocd

func GetSLice(values interface{}) []string {
	newValue := make([]string, 0)
	for _, value := range values.([]interface{}) {
		newValue = append(newValue, value.(string))
	}

	return newValue
}
