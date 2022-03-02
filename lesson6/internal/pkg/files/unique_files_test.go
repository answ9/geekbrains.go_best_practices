package files_test

import (
	"crypto/sha512"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"lesson6/internal/pkg/files"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

//Tests for UniqueFiles
func TestNewUniqueFilesMap(t *testing.T) {
	log := logrus.New()
	uniqueFiles := files.NewUniqueFilesMap(log)

	filepath.Walk("../config", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Errorf(err.Error())
		}
		if !info.IsDir() && info.Name() != ".DS_Store" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				t.Errorf(err.Error())
			}
			digest := sha512.Sum512(data)
			uniqueFiles.Map[digest] = append(uniqueFiles.Map[digest], files.NewFile(path, info.Name()))
		}
		return nil
	})

	assert.Equal(t, len(uniqueFiles.Map), 2, "they should be equal")
	for k, _ := range uniqueFiles.Map {
		assert.Equal(t, len(uniqueFiles.Map[k]), 1, "they should be equal")
	}
}

func ExampleNewUniqueFilesMap() {
	log := logrus.New()
	uniqueFiles := files.NewUniqueFilesMap(log)
	fmt.Println(uniqueFiles.Map)
	// Output:
	//map[]
}

var tests = []struct {
	input []files.File
	want  []files.File
}{
	{
		input: []files.File{
			{"../folder1/1 copy.txt", "1 copy.txt"}, {"../folder1/1 copy 2.txt", "1 copy 2.txt"}, {"../folder1/1.txt", "1.txt"},
		},
		want: []files.File{
			{"../folder1/1.txt", "1.txt"}, {"../folder1/1 copy.txt", "1 copy.txt"}, {"../folder1/1 copy 2.txt", "1 copy 2.txt"},
		},
	},
	{
		input: []files.File{
			{"", "name a.txt"}, {"", "name abc.txt"}, {"", "name ab.txt"},
		},
		want: []files.File{
			{"", "name a.txt"}, {"", "name ab.txt"}, {"", "name abc.txt"},
		},
	},
	{
		input: []files.File{
			{"", "Screenshot 2021-10-23 at 21.32.20.png"}, {"", "Screenshot 2021-10-23 at 21.32.20 3 2.png"}, {"", "Screenshot 2021-10-23 at 21.32.20 2.png"}, {"", "Screenshot 2021-10-23 at 21.32.20 3 3.png"},
		},
		want: []files.File{
			{"", "Screenshot 2021-10-23 at 21.32.20.png"}, {"", "Screenshot 2021-10-23 at 21.32.20 2.png"}, {"", "Screenshot 2021-10-23 at 21.32.20 3 2.png"}, {"", "Screenshot 2021-10-23 at 21.32.20 3 3.png"},
		},
	},
}

func TestUniqueFiles_Sort(t *testing.T) {
	log := logrus.New()
	uniqueFiles := files.NewUniqueFilesMap(log)

	for i, tt := range tests {
		digest := sha512.Sum512([]byte(strconv.Itoa(i)))
		uniqueFiles.Map[digest] = tt.input
		uniqueFiles.Sort()
		assert.Equal(t, tt.want, uniqueFiles.Map[digest], "they should be equal")
	}
}

func BenchmarkUniqueFiles_Sort(b *testing.B) {
	log := logrus.New()
	uniqueFiles := files.NewUniqueFilesMap(log)
	for i, tt := range tests {
		digest := sha512.Sum512([]byte(strconv.Itoa(i)))
		uniqueFiles.Map[digest] = tt.input
	}
	for j := 0; j < b.N; j++ {
		uniqueFiles.Sort()
	}
}
