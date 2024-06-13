package types

type CustomError struct {
	Message    string
	StatusCode int
}

func (err *CustomError) Error() string {
	return err.Message
}

func (err *CustomError) GetStatusCode() int {
	return err.StatusCode
}
