package slice

import "strings"

func DistinctLowerStrings(list []string) []string {
	dict := make(map[string]struct{}, len(list))
	for _, str := range list {
		dict[strings.ToLower(str)] = struct{}{}
	}
	list = list[:0]
	for str := range dict {
		list = append(list, str)
	}
	return list
}

func CheckElemInSlice[T any](arr []T, equal func(T) bool) bool {
	for _, elem := range arr {
		if equal(elem) {
			return true
		}
	}
	return false
}

func MapSlice[T any, R any](arr []T, toMap func(T) R) []R {
	res := []R{}
	for _, r := range arr {
		res = append(res, toMap(r))
	}
	return res
}

func ToMap[T any, V any, K comparable](arr []T, key func(T) (K, V)) map[K]V {
	res := map[K]V{}
	for _, r := range arr {
		k, v := key(r)
		res[k] = v
	}
	return res
}

func FilterSlice[T any](arr []T, filter func(T) bool) []T {
	res := []T{}
	for _, el := range arr {
		if filter(el) {
			res = append(res, el)
		}
	}
	return res
}

func ChunkSlice[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for {
		if len(slice) == 0 {
			break
		}

		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
