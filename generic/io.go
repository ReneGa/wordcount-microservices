package generic

import (
	"github.com/stretchr/testify/mock"
)

type MockReader struct{ mock.Mock }

func (w *MockWriter) Read(p []byte) (n int, err error) {
	returns := w.Called(p)
	return returns.Int(0), returns.Error(1)
}

type MockWriter struct{ mock.Mock }

func (w *MockWriter) Write(p []byte) (n int, err error) {
	returns := w.Called(p)
	return returns.Int(0), returns.Error(1)
}

type MockCloser struct{ mock.Mock }

func (w *MockWriter) Close() error {
	returns := w.Called()
	return returns.Error(0)
}

type MockReadWriteCloser struct {
	MockReader
	MockWriter
	MockCloser
}
