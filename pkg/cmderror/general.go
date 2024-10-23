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
