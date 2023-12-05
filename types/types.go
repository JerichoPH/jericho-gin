package types

import (
	"reflect"
	"strings"
)

type MapStringToAny = map[string]any
type MapStringToString = map[string]string
type MapStringToFloat32 = map[string]float32
type MapStringToFloat64 = map[string]float64
type MapStringToUint8 = map[string]uint8
type MapStringToUint32 = map[string]uint32
type MapStringToUint64 = map[string]uint64
type MapStringToUint = map[string]uint
type MapStringToInt8 = map[string]int8
type MapStringToInt32 = map[string]int32
type MapStringToInt64 = map[string]int64
type MapStringToInt = map[string]int
type MapStringToBool = map[string]bool
type MapIntToBool = map[int]bool
type MapInt8ToBool = map[int8]bool
type MapInt16ToBool = map[int16]bool
type MapInt32ToBool = map[int32]bool
type MapInt64ToBool = map[int64]bool
type MapUintToBool = map[uint]bool
type MapUint8ToBool = map[uint8]bool
type MapUint16ToBool = map[uint16]bool
type MapUint32ToBool = map[uint32]bool
type MapUint64ToBool = map[uint64]bool
type ListString = []string
type ListInt = []int
type ListInt8 = []int8
type ListInt16 = []int16
type ListInt32 = []int32
type ListInt64 = []int64
type ListUint = []uint
type ListUint8 = []uint8
type ListUint16 = []uint16
type ListUint32 = []uint32
type ListUint64 = []uint64
type ListFloat32 = []float32
type ListFloat64 = []float64
type ListAny = []any

// IsString 判断是否是字符串类型
func IsString(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.String
}

// IsNumber 判断是否是数字类型
func IsNumber(value any) bool {
	if IsInt(value) || IsInt8(value) || IsInt16(value) || IsInt32(value) || IsInt64(value) || IsUint(value) || IsUint8(value) || IsUint16(value) || IsUint32(value) || IsUint64(value) {
		return true
	}
	return false
}

// IsInt 判断是否是整数类型
func IsInt(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Int
}

// IsInt8 判断是否是int8类型
func IsInt8(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Int8
}

// IsInt16 判断是否是int8类型
func IsInt16(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Int16
}

// IsInt32 判断是否是int32类型
func IsInt32(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Int32
}

// IsInt64 判断是否是int64类型
func IsInt64(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Int64
}

// IsUint 判断是否是uint类型
func IsUint(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Int
}

// IsUint8 判断是否是uint8类型
func IsUint8(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Uint8
}

func IsUint16(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Uint16
}

// IsUint32 判断是否是uint32类型
func IsUint32(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Uint32
}

// IsUint64 判断是否是uint64类型
func IsUint64(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Uint64
}

// IsFloat 判断是否是浮点类型
func IsFloat(value any) bool {
	// if IsFloat32(value) {
	// 	isFloat = true
	// }
	// if IsFloat64(value) {
	// 	isFloat = true
	// }
	// return

	if IsFloat32(value) || IsFloat64(value) {
		return true
	}
	return false
}

// IsFloat32 判断是否是float32
func IsFloat32(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Float32
}

// IsFloat64 判断是否是float64
func IsFloat64(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Float64
}

// IsEmptyStruct 判断是否为空结构
func IsEmptyStruct(value, typeMode any) any {
	return reflect.Zero(reflect.TypeOf(value)).Interface()
}

// MapString 字符串字典处理
type MapString struct{}

// GetKeys 获取字符串字典所有key
func (receiver *MapString) GetKeys(value MapStringToAny) ListString {
	var keys = make(ListString, 0)

	for key := range value {
		keys = append(keys, key)
	}

	return keys
}

// IsTrue 判断字符串是否代表true
func IsTrue(value string) bool {
	switch strings.ToLower(value) {
	case "1":
	case "true":
	case "on":
	case "yes":
		return true
	}

	return false
}

// IsFalse 判断字符串是否代表false
func IsFalse(value string) bool {
	switch strings.ToLower(value) {
	case "0":
	case "false":
	case "off":
	case "no":
		return true
	}
	return false
}
