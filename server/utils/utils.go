package utils

import "strconv"

func Interface_to_float(data interface{}) float64 {
	f, err := strconv.ParseFloat(data.(string), 64)
	if err != nil {
		panic(err)
	}
	return f
}
func String_to_bool(data interface{}) bool {
	f, err := strconv.ParseBool(data.(string))
	if err != nil {
		panic(err)
	}
	return f
}
