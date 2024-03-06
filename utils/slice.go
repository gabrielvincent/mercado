package utils

import "reflect"

func Contains(slice interface{}, target interface{}) bool {
	sliceValue := reflect.ValueOf(slice)

	if sliceValue.Kind() != reflect.Slice {
		panic("Contains: the first argument must be a slice")
	}

	for i := 0; i < sliceValue.Len(); i++ {
		element := sliceValue.Index(i).Interface()
		if reflect.DeepEqual(element, target) {
			return true
		}
	}

	return false
}
