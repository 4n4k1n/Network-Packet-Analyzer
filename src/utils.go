package main

func getMaxOfMap(src map[string]int) string {

	var max_value int = 0
	var index_value string

	for k, v := range src {
		if v > max_value {
			max_value = v
			index_value = k
		}
	}
	return index_value
}
