package result

import (
	"errors"
	"testing"
)

func TestResult_IsSuccess(t *testing.T) {
	// Test success result
	successResult := NewSuccess("test")
	if !successResult.IsSuccess() {
		t.Errorf("Expected success result to be success, got failure")
	}
	
	// Test error result
	errorResult := NewError[string](errors.New("test error"))
	if errorResult.IsSuccess() {
		t.Errorf("Expected error result to be failure, got success")
	}
}

func TestResult_IsError(t *testing.T) {
	// Test success result
	successResult := NewSuccess("test")
	if successResult.IsError() {
		t.Errorf("Expected success result to not be error, got error")
	}
	
	// Test error result
	errorResult := NewError[string](errors.New("test error"))
	if !errorResult.IsError() {
		t.Errorf("Expected error result to be error, got not error")
	}
}

func TestResult_GetError(t *testing.T) {
	// Test success result
	successResult := NewSuccess("test")
	if successResult.GetError() != nil {
		t.Errorf("Expected success result error to be nil, got %v", successResult.GetError())
	}
	
	// Test error result
	testErr := errors.New("test error")
	errorResult := NewError[string](testErr)
	if errorResult.GetError() != testErr {
		t.Errorf("Expected error result error to be %v, got %v", testErr, errorResult.GetError())
	}
}

func TestResult_GetData(t *testing.T) {
	// Test success result
	testData := "test"
	successResult := NewSuccess(testData)
	if successResult.GetData() != testData {
		t.Errorf("Expected success result data to be %v, got %v", testData, successResult.GetData())
	}
	
	// Test error result
	errorResult := NewError[string](errors.New("test error"))
	if errorResult.GetData() != "" {
		t.Errorf("Expected error result data to be empty string, got %v", errorResult.GetData())
	}
}

func TestResult_Unwrap(t *testing.T) {
	// Test success result
	testData := "test"
	successResult := NewSuccess(testData)
	
	// Use a defer/recover to catch the panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected no panic, got %v", r)
		}
	}()
	
	if successResult.Unwrap() != testData {
		t.Errorf("Expected unwrapped data to be %v, got %v", testData, successResult.Unwrap())
	}
}

func TestResult_UnwrapPanic(t *testing.T) {
	// Test error result
	testErr := errors.New("test error")
	errorResult := NewError[string](testErr)
	
	// Use a defer/recover to catch the panic
	defer func() {
		if r := recover(); r != testErr {
			t.Errorf("Expected panic with %v, got %v", testErr, r)
		}
	}()
	
	// This should panic
	_ = errorResult.Unwrap()
	
	// If we get here, the test failed
	t.Errorf("Expected panic, but no panic occurred")
}

func TestResult_UnwrapOr(t *testing.T) {
	// Test success result
	testData := "test"
	fallback := "fallback"
	successResult := NewSuccess(testData)
	if successResult.UnwrapOr(fallback) != testData {
		t.Errorf("Expected unwrapped data to be %v, got %v", testData, successResult.UnwrapOr(fallback))
	}
	
	// Test error result
	errorResult := NewError[string](errors.New("test error"))
	if errorResult.UnwrapOr(fallback) != fallback {
		t.Errorf("Expected unwrapped data to be %v, got %v", fallback, errorResult.UnwrapOr(fallback))
	}
}

func TestResult_UnwrapOrElse(t *testing.T) {
	// Test success result
	testData := "test"
	successResult := NewSuccess(testData)
	if successResult.UnwrapOrElse(func(err error) string {
		return "fallback"
	}) != testData {
		t.Errorf("Expected unwrapped data to be %v, got %v", testData, successResult.UnwrapOrElse(func(err error) string {
			return "fallback"
		}))
	}
	
	// Test error result
	testErr := errors.New("test error")
	errorResult := NewError[string](testErr)
	if errorResult.UnwrapOrElse(func(err error) string {
		if err != testErr {
			t.Errorf("Expected error to be %v, got %v", testErr, err)
		}
		return "fallback"
	}) != "fallback" {
		t.Errorf("Expected unwrapped data to be fallback, got %v", errorResult.UnwrapOrElse(func(err error) string {
			return "fallback"
		}))
	}
}

func TestMap(t *testing.T) {
	// Test success result
	testData := "test"
	successResult := NewSuccess(testData)
	mappedResult := Map(successResult, func(s string) int {
		return len(s)
	})
	if !mappedResult.IsSuccess() {
		t.Errorf("Expected mapped result to be success, got failure")
	}
	if mappedResult.GetData() != 4 {
		t.Errorf("Expected mapped data to be 4, got %v", mappedResult.GetData())
	}
	
	// Test error result
	testErr := errors.New("test error")
	errorResult := NewError[string](testErr)
	mappedErrorResult := Map(errorResult, func(s string) int {
		return len(s)
	})
	if !mappedErrorResult.IsError() {
		t.Errorf("Expected mapped result to be error, got success")
	}
	if mappedErrorResult.GetError() != testErr {
		t.Errorf("Expected mapped error to be %v, got %v", testErr, mappedErrorResult.GetError())
	}
}

func TestFlatMap(t *testing.T) {
	// Test success result
	testData := "test"
	successResult := NewSuccess(testData)
	flatMappedResult := FlatMap(successResult, func(s string) Result[int] {
		return NewSuccess(len(s))
	})
	if !flatMappedResult.IsSuccess() {
		t.Errorf("Expected flat mapped result to be success, got failure")
	}
	if flatMappedResult.GetData() != 4 {
		t.Errorf("Expected flat mapped data to be 4, got %v", flatMappedResult.GetData())
	}
	
	// Test error in original result
	testErr := errors.New("test error")
	errorResult := NewError[string](testErr)
	flatMappedErrorResult := FlatMap(errorResult, func(s string) Result[int] {
		return NewSuccess(len(s))
	})
	if !flatMappedErrorResult.IsError() {
		t.Errorf("Expected flat mapped result to be error, got success")
	}
	if flatMappedErrorResult.GetError() != testErr {
		t.Errorf("Expected flat mapped error to be %v, got %v", testErr, flatMappedErrorResult.GetError())
	}
	
	// Test error in mapping function
	mappingErr := errors.New("mapping error")
	flatMappedMappingErrorResult := FlatMap(successResult, func(s string) Result[int] {
		return NewError[int](mappingErr)
	})
	if !flatMappedMappingErrorResult.IsError() {
		t.Errorf("Expected flat mapped result to be error, got success")
	}
	if flatMappedMappingErrorResult.GetError() != mappingErr {
		t.Errorf("Expected flat mapped error to be %v, got %v", mappingErr, flatMappedMappingErrorResult.GetError())
	}
}

func TestResult_String(t *testing.T) {
	// Test success result
	testData := "test"
	successResult := NewSuccess(testData)
	expectedString := "Success: test"
	if successResult.String() != expectedString {
		t.Errorf("Expected string to be %q, got %q", expectedString, successResult.String())
	}
	
	// Test error result
	testErr := errors.New("test error")
	errorResult := NewError[string](testErr)
	expectedErrorString := "Error: test error"
	if errorResult.String() != expectedErrorString {
		t.Errorf("Expected string to be %q, got %q", expectedErrorString, errorResult.String())
	}
}
