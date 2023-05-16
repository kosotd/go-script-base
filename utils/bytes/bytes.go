package bytes

import "encoding/hex"

func HexToBytes(hexStr string) [4]byte {
	res, _ := hex.DecodeString(hexStr)
	var dst [4]byte
	copy(dst[:], res)
	return dst
}

func ParseValue[T any](value interface{}) (t T) {
	if value == nil {
		return t
	}

	if result, ok := value.(T); ok {
		return result
	}

	return t
}
