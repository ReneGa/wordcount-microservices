package generic

import (
	"io"
	"os"
	"time"

	"github.com/stretchr/testify/mock"
)

type OS interface {
	MkdirAll(path string, perm os.FileMode) error
	OpenFile(name string, flag int, perm os.FileMode) (io.ReadWriteCloser, error)
}

type RealOS struct{}

func (o *RealOS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (o *RealOS) OpenFile(name string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	return os.OpenFile(name, flag, perm)
}

type MockOS struct{ mock.Mock }

func toReadWriteCloser(v interface{}) io.ReadWriteCloser {
	if v == nil {
		return nil
	}
	return v.(io.ReadWriteCloser)
}

func (o *MockOS) MkdirAll(path string, perm os.FileMode) error {
	returns := o.Called(path, perm)
	return returns.Error(0)
}

func (o *MockOS) OpenFile(name string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	returns := o.Called(name, flag, perm)
	return toReadWriteCloser(returns.Get(0)), returns.Error(1)
}

type MockFileInfo struct {
	NameField    string
	SizeField    int64
	ModeField    os.FileMode
	ModTimeField time.Time
	IsDirField   bool
	SysField     interface{}
}

func (f *MockFileInfo) Name() string {
	return f.NameField
}
func (f *MockFileInfo) Size() int64 {
	return f.SizeField
}
func (f *MockFileInfo) Mode() os.FileMode {
	return f.ModeField
}
func (f *MockFileInfo) ModTime() time.Time {
	return f.ModTimeField
}
func (f *MockFileInfo) IsDir() bool {
	return f.IsDirField
}
func (f *MockFileInfo) Sys() interface{} {
	return f.SysField
}
