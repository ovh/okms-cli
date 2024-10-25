package flagsmgmt

import "errors"

type OutputFormat string

const (
	JSON_OUTPUT_FORMAT OutputFormat = "json"
	TEXT_OUTPUT_FORMAT OutputFormat = "text"
)

func (e *OutputFormat) String() string {
	return string(*e)
}

func (e *OutputFormat) Set(v string) error {
	switch v {
	case "json", "text":
		*e = OutputFormat(v)
		return nil
	default:
		return errors.New(`must be one of "text", "json"`)
	}
}

func (e *OutputFormat) Type() string {
	return "text|json"
}
