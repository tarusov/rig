package validator_test

import (
	"fmt"
	"testing"

	"github.com/tarusov/rig/validator"
)

type validatorTestStruct struct {
	RequiredString string `customtag:"custom_string" validate:"oneof=a b x"`
}

func TestValidator(t *testing.T) {

	var vd, err = validator.New(
		validator.WithTagName("customtag"),
	)
	if err != nil {
		t.Error(err)
	}

	field, err := vd.Struct(validatorTestStruct{})
	if field != nil {
		fmt.Println(field.Error(), fmt.Sprintf("%T", err))
	}
}
