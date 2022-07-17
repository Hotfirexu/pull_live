package bytes

import "unsafe"

type sliceT struct {
	array unsafe.Pointer
	len   int
	cap   int
}

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func Bytes2StringRef(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func String2BytesRef(s string) []byte {
	str := (*stringStruct)(unsafe.Pointer(&s))
	ret := sliceT{array: unsafe.Pointer(str.str), len: str.len, cap: str.len}
	return *(*[]byte)(unsafe.Pointer(&ret))
}

// Sub
//
// 注意，内部会处理`b`大小不够，越界访问等情况
//
func Sub(b []byte, index int, length int) []byte {
	if index >= len(b) {
		return nil
	}

	if index+length > len(b) {
		return b[index:]
	}

	return b[index : index+length]
}

func Prefix(b []byte, length int) []byte {
	return Sub(b, 0, length)
}
