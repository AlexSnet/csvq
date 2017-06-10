package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type writeTest struct {
	Name     string
	Filename string
	Content  string
	Result   string
	Error    string
}

var createFileTests = []writeTest{
	{
		Name:     "Create",
		Filename: "create.txt",
		Content:  "write",
		Result:   "write",
	},
	{
		Name:     "Output to Stdout",
		Filename: "",
		Content:  "write",
		Result:   "write",
	},
	{
		Name:     "File Exists Error",
		Filename: "create.txt",
		Error:    fmt.Sprintf("file %s already exists", GetTestFilePath("create.txt")),
	},
	{
		Name:     "File Open Error",
		Filename: path.Join("notexistdir", "create.txt"),
		Error:    fmt.Sprintf("open %s: no such file or directory", GetTestFilePath(path.Join("notexistdir", "create.txt"))),
	},
}

func TestCreateFile(t *testing.T) {
	for _, v := range createFileTests {
		if len(v.Filename) < 1 {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			ToStdout(v.Content)

			w.Close()
			os.Stdout = oldStdout

			buf, _ := ioutil.ReadAll(r)
			if string(buf) != v.Result {
				t.Errorf("%s: content = %q, want %q", v.Name, string(buf), v.Result)
			}
		} else {
			filename := GetTestFilePath(v.Filename)
			err := CreateFile(filename, v.Content)
			if err != nil {
				if len(v.Error) < 1 {
					t.Errorf("%s: unexpected error %q", v.Name, err)
				} else if err.Error() != v.Error {
					t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
				}
				continue
			}
			if 0 < len(v.Error) {
				t.Errorf("%s: no error, want error %q", v.Name, v.Error)
				continue
			}

			fp, _ := os.Open(filename)
			buf, _ := ioutil.ReadAll(fp)
			if string(buf) != v.Result {
				t.Errorf("%s: content = %q, want %q", v.Name, string(buf), v.Result)
			}
		}
	}
}

var updateFileTests = []writeTest{
	{
		Name:     "Update",
		Filename: "create.txt",
		Content:  "truncate and write",
		Result:   "truncate and write",
	},
	{
		Name:     "File Not Found Error",
		Filename: "notexist.txt",
		Error:    fmt.Sprintf("open %s: no such file or directory", GetTestFilePath("notexist.txt")),
	},
}

func TestUpdateFile(t *testing.T) {
	for _, v := range updateFileTests {
		filename := GetTestFilePath(v.Filename)
		err := UpdateFile(filename, v.Content)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}

		fp, _ := os.Open(filename)
		buf, _ := ioutil.ReadAll(fp)
		if string(buf) != v.Result {
			t.Errorf("%s: content = %q, want %q", v.Name, string(buf), v.Result)
		}
	}
}
