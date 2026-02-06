package cmderror

type SomethingWentWrong struct{}

func (err *SomethingWentWrong) Error() string {
	return "Something went wrong. Please try again."
}

type ActionForbidden struct{}

func (err *ActionForbidden) Error() string {
	return "This action is forbidden"
}

type InvalidNumberOfArguments struct{}

func (err *InvalidNumberOfArguments) Error() string {
	return "Invalid number of arguments"
}

type InvalidOperation struct{}

func (err *InvalidOperation) Error() string {
	return `Invalid operation`
}

type Unexpected struct{}

func (err *Unexpected) Error() string {
	return `something unexpected happened. try again. or report`
}
