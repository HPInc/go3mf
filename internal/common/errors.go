package common

// ErrorCode defines a code associated to an error
type ErrorCode int32

// Error contains an error code that has to match with a valid ErrorCode.
// It is mainly used to translate codes to strings
type Error struct {
	errorCode ErrorCode
}

// NewError returns a new error with the desired code
func NewError(errorCode ErrorCode) Error {
	return Error{
		errorCode: errorCode,
	}
}

// GetErrorCode returns the error code
func (e *Error) GetErrorCode() ErrorCode {
	return e.errorCode
}

func (e *Error) Error() string {
	if val, ok := ErrorMap[e.errorCode]; ok {
		return val
	}
	return "unknown error"
}
