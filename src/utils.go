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

// Get top N entries from a map
func getTopNFromMap(src map[string]int, n int) []KeyValue {
	var pairs []KeyValue
	for k, v := range src {
		pairs = append(pairs, KeyValue{k, v})
	}
	
	// Simple bubble sort for top N (sufficient for small datasets)
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i].Value < pairs[j].Value {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
	
	if len(pairs) < n {
		n = len(pairs)
	}
	
	return pairs[:n]
}
