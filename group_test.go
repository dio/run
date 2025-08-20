package run

import (
	"errors"
	"testing"
)

func TestAddReturnsError(t *testing.T) {
	g := New()
	expectedErr := errors.New("test error")

	g.Add(func() error {
		return expectedErr
	}, func(error) {})

	err := g.Run()

	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestAddReturnsPanicRecovery(t *testing.T) {
	g := New()
	expectedPanic := "test panic"

	g.Add(func() error {
		panic(expectedPanic)
	}, func(err error) {})

	err := g.Run()

	if err == nil {
		t.Fatal("Expected panic to be caught, but got nil")
	}

	panicErr, ok := err.(*PanicError)
	if !ok {
		t.Fatalf("Expected PanicError, got %T", err)
	}

	if panicErr.Value != expectedPanic {
		t.Errorf("Expected panic value %v, got %v", expectedPanic, panicErr.Value)
	}

	// Test the Error() method to achieve 100% coverage
	errorMsg := panicErr.Error()
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}
}

func TestAddWithCustomPanicHandler(t *testing.T) {
	customError := errors.New("custom panic error")
	handlerCalled := false
	var capturedPanic any
	var capturedStack string

	g := NewWithHandler(func(panicValue any, stackTrace string) error {
		handlerCalled = true
		capturedPanic = panicValue
		capturedStack = stackTrace
		return customError
	})

	g.Add(func() error {
		panic("test panic")
	}, func(err error) {})

	err := g.Run()

	if err != customError {
		t.Errorf("Expected error %v, got %v", customError, err)
	}

	if !handlerCalled {
		t.Error("Expected panic handler to be called, but it wasn't")
	}

	if capturedPanic != "test panic" {
		t.Errorf("Expected captured panic %v, got %v", "test panic", capturedPanic)
	}

	if capturedStack == "" {
		t.Error("Expected captured stack trace to be non-empty")
	}
}
