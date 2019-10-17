package operator

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func tmpCreateFile(name, content string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(content))
	if err != nil {
		return err
	}
	return nil
}

func TestOverwriteSuccess(t *testing.T) {
	testFile := ".test_overwrite"
	if err := tmpCreateFile(testFile, "test\n"); err != nil {
		t.Errorf("create file failed: %v", err)
	}

	succ := func(old, new string, info os.FileInfo) error {
		return nil
	}

	action := OverwriteOption(succ)
	stat, _ := os.Stat(testFile)

	err := action(testFile, testFile, stat)
	if err != nil {
		t.Errorf("run overwrite-succ failed: %v", err)
	}

	_, err = os.Stat(testFile)
	if err == nil {
		t.Errorf("%s should be deleted", testFile)
	}

}

func TestOverwriteFailed(t *testing.T) {
	testFile := ".test_overwrite"
	if err := tmpCreateFile(testFile, "test\n"); err != nil {
		t.Errorf("create file failed: %v", err)
	}

	failed := func(old, new string, info os.FileInfo) error {
		return fmt.Errorf("test error: %s", testFile)
	}

	action := OverwriteOption(failed)
	stat, _ := os.Stat(testFile)

	err := action(testFile, testFile, stat)
	if err != nil {
		t.Errorf("run overwrite-succ failed: %v", err)
	}

	_, err = os.Stat(testFile)
	if err != nil {
		t.Errorf("%s should existing", testFile)
	}

	os.Remove(testFile)

}

func TestCopyFileAction(t *testing.T) {
	testFile := ".test_file"
	err := tmpCreateFile(testFile, "data1")
	if err != nil {
		t.Errorf("Failed to create file: %v", err)
	}
	defer os.Remove(testFile)

	info, _ := os.Stat(testFile)

	target := testFile + ".copy"
	defer os.Remove(target)

	err = CopyFileAction(testFile, target, info)
	if err != nil {
		t.Errorf("coyp failed: %v", err)
	}

	data, err := ioutil.ReadFile(target)
	if string(data) != "data1" {
		t.Errorf("unexpected copy value: %s (%v)", data, err)
	}

	err = CopyFileAction(testFile, target, info)
	if err == nil || (err != nil && err.Error() != fileExists) {
		t.Errorf("copy should be skipped: %v", err)
	}

}

func TestCreateLinkAction(t *testing.T) {
	testFile := ".test_file"
	err := tmpCreateFile(testFile, "data2")
	if err != nil {
		t.Errorf("Failed to create file: %v", err)
	}
	defer os.Remove(testFile)

	info, _ := os.Stat(testFile)

	target := testFile + ".link"
	defer os.Remove(target)

	err = CreateLinkAction(testFile, target, info)
	if err != nil {
		t.Errorf("link failed: %v", err)
	}

	data, err := ioutil.ReadFile(target)
	if string(data) != "data2" {
		t.Errorf("unexpected link value: %s (%v)", data, err)
	}

	err = CreateLinkAction(testFile, target, info)
	if err == nil || (err != nil && err.Error() != fileExists) {
		t.Errorf("link should be skipped: %v", err)
	}
}
