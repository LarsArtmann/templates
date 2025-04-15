package result

import "fmt"

// Result represents the result of an operation with a generic data type
type Result[T any] struct {
	Success bool
	Error   error
	Data    T
}

// NewSuccess creates a new success result
func NewSuccess[T any](data T) Result[T] {
	return Result[T]{
		Success: true,
		Data:    data,
	}
}

// NewError creates a new error result
func NewError[T any](err error) Result[T] {
	return Result[T]{
		Success: false,
		Error:   err,
	}
}

// IsSuccess returns true if the result is a success
func (r Result[T]) IsSuccess() bool {
	return r.Success
}

// IsError returns true if the result is an error
func (r Result[T]) IsError() bool {
	return !r.Success
}

// GetError returns the error if the result is an error, nil otherwise
func (r Result[T]) GetError() error {
	return r.Error
}

// GetData returns the data if the result is a success, nil otherwise
func (r Result[T]) GetData() T {
	return r.Data
}

// Unwrap returns the data if the result is a success, or panics with the error
func (r Result[T]) Unwrap() T {
	if r.IsError() {
		panic(r.Error)
	}
	return r.Data
}

// UnwrapOr returns the data if the result is a success, or the fallback value if it's an error
func (r Result[T]) UnwrapOr(fallback T) T {
	if r.IsError() {
		return fallback
	}
	return r.Data
}

// UnwrapOrElse returns the data if the result is a success, or calls the fallback function if it's an error
func (r Result[T]) UnwrapOrElse(fallback func(error) T) T {
	if r.IsError() {
		return fallback(r.Error)
	}
	return r.Data
}

// Map applies a function to the data if the result is a success, or returns the error unchanged
func Map[T, U any](r Result[T], f func(T) U) Result[U] {
	if r.IsError() {
		return NewError[U](r.Error)
	}
	return NewSuccess(f(r.Data))
}

// FlatMap applies a function that returns a Result to the data if the result is a success
func FlatMap[T, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.IsError() {
		return NewError[U](r.Error)
	}
	return f(r.Data)
}

// String returns a string representation of the result
func (r Result[T]) String() string {
	if r.IsSuccess() {
		return fmt.Sprintf("Success: %v", r.Data)
	}
	return fmt.Sprintf("Error: %v", r.Error)
}
