package generic

import (
	"io/ioutil"
	"os"

	"github.com/stretchr/testify/mock"
)

type IOUtil interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
}

type RealIOUtil struct{}

func (i *RealIOUtil) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

type MockIOUtil struct{ mock.Mock }

func toSliceOfFileInfo(v interface{}) []os.FileInfo {
	if v == nil {
		return nil
	}
	return v.([]os.FileInfo)
}

func (i *MockIOUtil) ReadDir(dirname string) ([]os.FileInfo, error) {
	returns := i.Called(dirname)
	return toSliceOfFileInfo(returns.Get(0)), returns.Error(1)
}
