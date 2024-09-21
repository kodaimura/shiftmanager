package utils

import (
	"fmt"
	"reflect"
	"math/rand"
	"time"
	"strconv"
)


func AtoiSlice(sl []string) ([]int, error) {
	isl := make([]int, len(sl))
	for i, v := range sl {
		x , err := strconv.Atoi(v)
		isl[i] = x
		if err != nil {
			return []int{}, err
		}
	}
	return isl, nil
}


func ItoaSlice(sl []int) []string {
	asl := make([]string, len(sl))
	for i, v := range sl {
		asl[i] = strconv.Itoa(v)
	}

	return asl
} 


func RandomString(length int, options ...string) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if len(options) > 0 {
		charset = options[0]
	}
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}


func GetFieldValue(obj interface{}, fieldName string) (interface{}, error) {
	val := reflect.ValueOf(obj).Elem()
	fieldVal := val.FieldByName(fieldName)

	if !fieldVal.IsValid() {
		return nil, fmt.Errorf("No such field: %s", fieldName)
	}

	return fieldVal.Interface(), nil
}


func SetFieldValue(obj interface{}, fieldName string, newValue interface{}) error {
	val := reflect.ValueOf(obj).Elem()
	fieldVal := val.FieldByName(fieldName)

	if !fieldVal.IsValid() {
		return fmt.Errorf("No such field: %s", fieldName)
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Cannot set field: %s", fieldName)
	}

	newVal := reflect.ValueOf(newValue)
	if fieldVal.Type() != newVal.Type() {
		return fmt.Errorf("Provided value type didn't match obj field type")
	}

	fieldVal.Set(newVal)
	return nil
}


func IsZero(value interface{}) bool {
	if value == nil {
        return true
    }
    switch v := value.(type) {
    case string:
        return v == ""
    case int:
        return v == 0
    case float32:
        return v == 0.0
    case float64:
        return v == 0.0
    default:
        return false
    }
}