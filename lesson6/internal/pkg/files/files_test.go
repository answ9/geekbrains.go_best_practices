package files_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"lesson6/internal/pkg/files"
	"os"
	"path/filepath"
	"testing"
)

//Tests for File
func TestNewFile(t *testing.T) {
	want := []files.File{
		{"../config/config.go", "config.go"},
		{"../config/config_test.go", "config_test.go"}}

	got := []files.File{}

	filepath.Walk("../config", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Errorf(err.Error())
		}
		if !info.IsDir() && info.Name() != ".DS_Store" {
			got = append(got, files.NewFile(path, info.Name()))
		}
		return nil
	})

	assert.Equal(t, want, got, "they should be equal")
}

func ExampleNewFile() {
	file := files.NewFile("../config/config.go", "config.go")
	fmt.Println(file)
	// Output:
	//{../config/config.go config.go}
}
