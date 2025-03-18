package cache

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

const (
	B = 1 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
	EB
	DefaultMemSize    = 100 << 20
	DefaultMemSizeStr = "100MB"
)

func ParseSize(size string) (parseByteSize int64, parseByteSizeStr string) {
	defer func() {
		fmt.Println(B, KB, MB, GB, TB, GB, EB, DefaultMemSize, DefaultMemSizeStr)
		fmt.Println(parseByteSize, parseByteSizeStr)
		if parseByteSize == 0 || len(parseByteSizeStr) == 0 {
			log.Println("解析size失败，返回默认值")
			parseByteSize = DefaultMemSize
			parseByteSizeStr = DefaultMemSizeStr
		}
	}()

	length := len(size)
	if length < 2 || unicode.ToUpper(rune(size[length-1])) != 'B' {
		return
	}

	var unit string
	// B 只有一个字母位，单独做下处理
	if unicode.IsDigit(rune(size[length-2])) {
		unit = "B"
	} else {
		unit = strings.ToUpper(size[length-2:])
	}

	byteNum, err := strconv.Atoi(size[:length-len(unit)])
	if err != nil {
		return
	}
	parseByteSize = int64(byteNum)
	size = strings.ToUpper(size)

	switch unit {
	case "B":
		return parseByteSize * B, size
	case "KB":
		return parseByteSize * KB, size
	case "MB":
		return parseByteSize * MB, size
	case "GB":
		return parseByteSize * GB, size
	case "TB":
		return parseByteSize * TB, size
	case "PB":
		return parseByteSize * PB, size
	case "EB":
		return parseByteSize * EB, size
	default:
	}

	return
}

func GetValueSize(val any) int64 {
	size := calculateSize(reflect.ValueOf(val))
	return size
}

func calculateSize(v reflect.Value) int64 {
	switch v.Kind() {
	case reflect.Bool:
		return int64(unsafe.Sizeof(false))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(unsafe.Sizeof(int(0)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(unsafe.Sizeof(uint(0)))
	case reflect.Float32, reflect.Float64:
		return int64(unsafe.Sizeof(float64(0)))
	case reflect.Complex64, reflect.Complex128:
		return int64(unsafe.Sizeof(complex128(0)))
	case reflect.String:
		return int64(len(v.String()))
	case reflect.Array, reflect.Slice:
		var size int64
		for i := 0; i < v.Len(); i++ {
			size += calculateSize(v.Index(i))
		}
		return size
	case reflect.Map:
		var size int64
		for _, key := range v.MapKeys() {
			size += calculateSize(key) + calculateSize(v.MapIndex(key))
		}
		return size
	case reflect.Struct:
		var size int64
		for i := 0; i < v.NumField(); i++ {
			size += calculateSize(v.Field(i))
		}
		return size
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return 0
		}
		return calculateSize(v.Elem())
	default:
		return 0
	}
}
