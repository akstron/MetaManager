package cmderror

type InvalidPath struct{}

func (err *InvalidPath) Error() string {
	return "Provided path is invalid"
}

type AlreadyInitPath struct{}

func (err *AlreadyInitPath) Error() string {
	return `A .mm directory already exist. Cannot reinitialize the tree path.`
}
