package main

import "reflect"

func count_instances_in_slice[T any](slice []T, element T) int {
	count := 0
	for _, value := range slice {
		if reflect.DeepEqual(value, element) {
			count++
		}
	}
	return count
}
