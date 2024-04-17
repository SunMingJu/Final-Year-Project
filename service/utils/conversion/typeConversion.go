package conversion

import (
	"reflect"
	"strings"
	"unsafe"
)

// StringConversionMap 
func StringConversionMap(s string) []string {
	s = strings.TrimSpace(s)
	list := strings.Split(s, ",")
	//Returns an empty array if there is only one item and it is empty
	if len(list) == 1 && list[0] == "" {
		return make([]string, 0)
	}
	return list
}

// MapConversionString 
func MapConversionString(m []string) string {
	var srt string
	if len(m) != 0 {
		for _, v := range m {
			srt = srt + v + ","
		}
		srt = srt[:len(srt)-1]
		return srt
	}
	return ""
}

// StringImgConversionMap 
func StringImgConversionMap(s string) []string {
	list := strings.Split(s, ",")
	for k, v := range list {
		list[k] = FormattingSrc(v)
	}
	return list
}

//BoolTurnInt8 
func BoolTurnInt8(is bool) int8 {
	if is {
		return 1
	} else {
		return 0
	}
}

//Int8TurnBool 
func Int8TurnBool(i int8) bool {
	if i > 0 {
		return true
	} else {
		return false
	}
}

func IntTurnBool(i int) bool {
	if i > 0 {
		return true
	} else {
		return false
	}
}

// String2Bytes 
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// Bytes2String 
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
