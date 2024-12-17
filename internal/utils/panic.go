package utils

// for handling unexpected errors and better test coverage
func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
