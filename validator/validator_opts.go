package validator

import (
	"reflect"
	"strings"

	ut "github.com/go-playground/universal-translator"
	vd "github.com/go-playground/validator/v10"
)

// validatorOption is validator constructor optional modificator.
type validatorOption func(*validatorOptions)

// WithTranslator set custom translator for field messages.
func WithTranslator(translator ut.Translator) validatorOption {
	return func(vo *validatorOptions) {
		vo.translator = translator
	}
}

// WithValidationFunc setup custom validation/translation func for specified tag.
func WithValidationFunc(tag string, validateFn vd.Func, regTranslateFn vd.RegisterTranslationsFunc, translateFn vd.TranslationFunc) validatorOption {
	return func(vo *validatorOptions) {
		vo.customFns[tag] = validatorFnSet{
			validateFn:     validateFn,
			regTranslateFn: regTranslateFn,
			translateFn:    translateFn,
		}
	}
}

// WithTagName set custom tag (from docs, "json", "mapstructure", etc)
func WithTagName(tag string) validatorOption {
	return func(vo *validatorOptions) {
		vo.tagNameFn = func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get(tag), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		}
	}
}

// WithTagNameFunc set custom tag name extract func.
func WithTagNameFunc(tagNameFn vd.TagNameFunc) validatorOption {
	return func(vo *validatorOptions) {
		vo.tagNameFn = tagNameFn
	}
}
