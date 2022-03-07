package logger

import "fmt"

// List of pre-defined fields.
const (
	FieldNamePackage = "pkg"
	FieldNameMethod  = "fn"
)

// normalize key and value for logging.
// TODO: possible check for screening characters.
func normalize(k string, v interface{}) (string, string) {
	if k == "" {
		k = "unknown_field"
	}

	return k, fmt.Sprint(v)
}
