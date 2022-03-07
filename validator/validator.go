package validator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	vd "github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
)

type (
	// Validator struct.
	Validator struct {
		*vd.Validate
		ut.Translator
	}

	// validatorOptions is auxilary constructor struct.
	validatorOptions struct {
		translator ut.Translator
		customFns  map[string]validatorFnSet
		tagNameFn  vd.TagNameFunc
	}

	// validatorFnSet is auxilary constructor struct.
	validatorFnSet struct {
		validateFn     vd.Func
		regTranslateFn vd.RegisterTranslationsFunc
		translateFn    vd.TranslationFunc
	}
)

// New сreate new validator instance.
func New(opts ...validatorOption) (*Validator, error) {

	var vo = &validatorOptions{
		customFns: make(map[string]validatorFnSet),
	}

	for _, opt := range opts {
		opt(vo)
	}

	var v10 = vd.New()
	if vo.tagNameFn != nil {
		v10.RegisterTagNameFunc(vo.tagNameFn)
	}

	if vo.translator == nil {
		var (
			en = en.New()
			un = ut.New(en, en)
		)
		vo.translator, _ = un.GetTranslator("en")
	}

	if err := entrans.RegisterDefaultTranslations(v10, vo.translator); err != nil {
		return nil, err
	}

	for tag, fnSet := range vo.customFns {
		if err := v10.RegisterValidation(tag, fnSet.validateFn); err != nil {
			return nil, err
		}
		if err := v10.RegisterTranslation(tag, vo.translator, fnSet.regTranslateFn, fnSet.translateFn); err != nil {
			return nil, err
		}
	}

	return &Validator{
		Validate:   v10,
		Translator: vo.translator,
	}, nil
}

// Struct method return struct methods. If validation is
// get some errors it's return FieldErrors list, or error
// in case ov other errors.
func (v *Validator) Struct(i interface{}) (FieldErrors, error) {
	switch s := v.Validate.Struct(i).(type) {
	case vd.ValidationErrors:
		return v.mkFieldErrors(s), nil
	default:
		return nil, s
	}
}

// mkFieldErrors create field errors map.
func (v *Validator) mkFieldErrors(validationErrs vd.ValidationErrors) FieldErrors {

	if len(validationErrs) == 0 {
		return nil
	}

	var fieldErrors = make(FieldErrors)
	for _, err := range validationErrs {
		fieldErrors[err.Field()] = err.Translate(v.Translator)
	}

	return fieldErrors
}
