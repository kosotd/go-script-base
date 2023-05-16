package async

import (
	"fmt"
)

type Triple[P any, V any] struct {
	Param P
	Value V
	Error error
}

func MapValuesNoRecover[P any, V any](params []P, call func(P) (V, error)) []Triple[P, V] {
	return asyncMapValues(params, false, call)
}

func MapValues[P any, V any](params []P, call func(P) (V, error)) []Triple[P, V] {
	return asyncMapValues(params, true, call)
}

func asyncMapValues[P any, V any](params []P, needRecover bool, call func(P) (V, error)) []Triple[P, V] {
	chans := make(chan Triple[P, V], len(params))
	defer close(chans)

	var def V

	for _, param := range params {
		go func(param P, chans chan Triple[P, V]) {
			if needRecover {
				defer func() {
					if r := recover(); r != nil {
						chans <- Triple[P, V]{param, def, fmt.Errorf("%v", r)}
					}
				}()
			}

			value, err := call(param)
			chans <- Triple[P, V]{param, value, err}
		}(param, chans)
	}

	results := make([]Triple[P, V], 0, len(params))
	if len(params) < 1 {
		return results
	}

	cnt := 0
	for val := range chans {
		cnt++

		results = append(results, val)

		if cnt >= len(params) {
			break
		}
	}

	return results
}

func TriplesToMap[P comparable, V any](triples []Triple[P, V]) map[P]V {
	result := make(map[P]V, len(triples))
	for _, triple := range triples {
		if triple.Error != nil {
			continue
		}
		result[triple.Param] = triple.Value
	}
	return result
}
