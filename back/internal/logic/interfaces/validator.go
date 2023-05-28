package logicinterfaces

type Validator interface {
	CommonValidation(i interface{}) error
}
