package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"unsafe"
)

func IfErrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func ToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func MD5(str string) string {
	h := md5.New()
	h.Write(ToBytes(str))
	return hex.EncodeToString(h.Sum(nil))
}

// NextInt [0,n)
func NextInt(n int64) int64 {
	if n == 0 {
		return 0
	}
	return rand.Int63() % n
}

// Between [min, max]
func Between(min, max int64) int64 {
	if min >= max {
		return min
	}
	return NextInt(max-min) + min + 1
}
