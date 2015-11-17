package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aryann/difflib"
)

func TestMain(t *testing.T) {
	build, _ := filepath.Abs("build")
	filepath.Walk("test", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".sh") {
			cmd := exec.Command("bash", filepath.Base(path))
			cmd.Dir = filepath.Dir(path)
			cmd.Env = []string{"PATH=" + build}
			stderr := new(bytes.Buffer)
			cmd.Stderr = stderr
			output, err := cmd.Output()
			if err != nil {
				t.Errorf("FAIL: execution failed: " + path + ": " + err.Error())
			} else {
				outfile := strings.TrimSuffix(path, filepath.Ext(path)) + ".txt"
				expected, err := ioutil.ReadFile(outfile)
				if err != nil {
					t.Errorf("FAIL: error on reading output file: " + outfile)
				} else {
					diffs := difflib.Diff(strings.Split(stderr.String()+string(output), "\n"), strings.Split(string(expected), "\n"))
					help := strings.Contains(string(output), "NAME:")
					differs := false
					for _, diff := range diffs {
						differs = differs || (help && diff.Delta == difflib.RightOnly || !help && diff.Delta != difflib.Common)
					}
					if differs {
						buf := bytes.NewBufferString("")
						for _, diff := range diffs {
							if diff.Delta != difflib.Common {
								buf.WriteString(diff.String() + "\n")
							}
						}
						t.Errorf("FAIL: output differs: " + path + "\n" + buf.String())
					} else {
						t.Logf("PASS: " + path + "\n")
					}
				}
			}
		}
		return nil
	})
}
