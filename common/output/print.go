package output

import (
	"encoding/json"
	"fmt"
	"os"
)

func JsonPrint(resp any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(resp); err != nil {
		fmt.Fprintln(os.Stderr, "Fail to marshall response")
		os.Exit(1)
	}
}
