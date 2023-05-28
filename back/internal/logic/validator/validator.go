package validator

import (
	"fmt"
	"reflect"
	logicerrors "svalka-service/internal/logic/errors"
	logicinterfaces "svalka-service/internal/logic/interfaces"

	v10 "github.com/go-playground/validator/v10"
)

type commonValidator struct {
	*v10.Validate
}

var v commonValidator

func init() {
	vr := v10.New()
	vr.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tag := fld.Tag.Get("json")
		if tag != "" {
			return tag
		}
		return fld.Name
	})
	v = commonValidator{
		Validate: vr,
	}
}

func NewValidator() logicinterfaces.Validator {
	return v
}

// CommonValidation ...
func (v commonValidator) CommonValidation(i interface{}) error {
	if i == nil {
		return logicerrors.NilInterfaceError
	}
	// by tag
	err := v.Struct(i)
	if err != nil {
		// return the first error
		validationErrors := err.(v10.ValidationErrors)
		for _, nextErr := range validationErrors {
			return fmt.Errorf("\"%s\" must be %s\n",
				nextErr.Field(),
				nextErr.Tag(),
			)
		}
	}
	return nil
}
