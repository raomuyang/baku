package operator

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestGetFilterOption(t *testing.T) {

	option, err := GetFilterOption("\\.git/config")
	if err != nil {
		t.Errorf("failed to get option: %v", err)
	}

	f := func(old, new string, info os.FileInfo) error {
		return fmt.Errorf("test")
	}

	path := "/path/to/.git/config.yml"
	result := option(f)(path, "", nil)
	if result.Error() != fileIgnored {
		t.Errorf("%s should be ignored: %v", path, result)
	}

	path = "/path/to/.gitt/config.yml"
	result = option(f)(path, "", nil)
	if result.Error() != "test" {
		t.Errorf("%s shouldn't be ignored: %v", path, result)
	}
}

func TestExpandUserHome(t *testing.T) {
	path := "~/~/123"
	newPath, err := ExpandUserHome(path)
	if err != nil {
		t.Errorf("failed to expand the tilde path: %v", err)
	}

	if !strings.HasSuffix(newPath, "~/123") || strings.HasPrefix(newPath, "~/") {
		t.Errorf("expand failed: %s", newPath)
	}
}

func TestCustomCopyFileAction(t *testing.T) {
	testFile := ".test_file"
	err := tmpCreateFile(testFile, "data1")
	if err != nil {
		t.Errorf("Failed to create file: %v", err)
	}
	defer os.Remove(testFile)

	info, _ := os.Stat(testFile)

	target := testFile + ".copy"
	defer os.Remove(target)

	err = GetCustomCopyAction("cp")(testFile, target, info)
	if err != nil {
		t.Errorf("coyp failed: %v", err)
	}

	data, err := ioutil.ReadFile(target)
	if string(data) != "data1" {
		t.Errorf("unexpected copy value: %s (%v)", data, err)
	}

	err = GetCustomCopyAction("cp")(testFile, target, info)
	if err == nil || (err != nil && err.Error() != fileExists) {
		t.Errorf("copy should be skipped: %v", err)
	}

}
