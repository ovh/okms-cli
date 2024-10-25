package restflags

import (
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
