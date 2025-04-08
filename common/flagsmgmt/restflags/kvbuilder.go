package restflags

import (
	"fmt"
	"io"

	kvbuilder "github.com/hashicorp/go-secure-stdlib/kv-builder"
)

func ParseArgsData(stdin io.Reader, args []string) (map[string]interface{}, error) {
	builder := &kvbuilder.Builder{Stdin: stdin}
	if err := builder.Add(args...); err != nil {
		return nil, err
	}

	return builder.Map(), nil
}

func ParseArgsCustomMetadata(stdin io.Reader, args []string) (map[string]string, error) {
	builder := &kvbuilder.Builder{Stdin: stdin}
	if err := builder.Add(args...); err != nil {
		return nil, err
	}

	m := map[string]string{}
	for key, value := range builder.Map() {
		m[key] = fmt.Sprintf("%v", value)
	}
	return m, nil
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
