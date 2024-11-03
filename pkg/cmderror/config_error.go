package cmderror

type InvalidPath struct{}

func (err *InvalidPath) Error() string {
	return "Provided path is invalid"
}

type AlreadyInitPath struct{}

func (err *AlreadyInitPath) Error() string {
	return `A .mm directory already exist. Cannot reinitialize the tree path.`
}

type UninitializedRoot struct{}

func (err *UninitializedRoot) Error() string {
	return `Current root is uninitialized`
}

type InvalidOperation struct{}

func (err *InvalidOperation) Error() string {
	return `Invalid operation`
}
