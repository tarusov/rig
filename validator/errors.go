package validator

import "strings"

// FieldErrors is field errors type.
type FieldErrors map[string]string

// Error implement error interface method.
func (fe FieldErrors) Error() string {

	var sb strings.Builder
	for fieldName, errMsg := range fe {
		sb.WriteString(fieldName + ": " + errMsg + "\n")
	}

	return sb.String()
}
