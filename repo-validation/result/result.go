package result

// Result represents the result of an operation
type Result struct {
	Success bool
	Error   error
	Data    interface{}
}

// NewSuccess creates a new success result
func NewSuccess(data interface{}) Result {
	return Result{
		Success: true,
		Data:    data,
	}
}

// NewError creates a new error result
func NewError(err error) Result {
	return Result{
		Success: false,
		Error:   err,
	}
}

// IsSuccess returns true if the result is a success
func (r Result) IsSuccess() bool {
	return r.Success
}

// IsError returns true if the result is an error
func (r Result) IsError() bool {
	return !r.Success
}

// GetError returns the error if the result is an error, nil otherwise
func (r Result) GetError() error {
	return r.Error
}

// GetData returns the data if the result is a success, nil otherwise
func (r Result) GetData() interface{} {
	return r.Data
}
