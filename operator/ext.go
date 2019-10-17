package operator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
)

// GetFilterOption can build a new option to filter the path (ignore the specific files by the giving regexp)
func GetFilterOption(ignoresRegex string) (BackupOption, error) {
	exp, err := regexp.Compile(ignoresRegex)
	if err != nil {
		return nil, err
	}
	opt := func(action CopyAction) CopyAction {
		newAct := func(oldPath, newPath string, info os.FileInfo) error {

			matched := exp.Match([]byte(oldPath))
			if matched {
				return fmt.Errorf(fileIgnored)
			}
			return action(oldPath, newPath, info)
		}
		return newAct
	}
	return opt, nil
}

// GetCustomCopyAction allows users to use a custom program to complete the copy operation, such as "cp"
func GetCustomCopyAction(command string) CopyAction {

	return func(oldPath, newPath string, info os.FileInfo) error {
		isDir, err := checkAndCreateDir(oldPath, newPath, info)
		if err != nil {
			return err
		}
		if isDir && info.IsDir() {
			log.Printf("Operated directory: %s to %s", oldPath, newPath)
			return nil
		}

		_, err = os.Stat(newPath)
		if err == nil {
			return innerError{message: fileExists}
		}

		p := exec.Command(command, oldPath, newPath)
		p.Stdout = os.Stdout
		p.Stderr = os.Stderr
		return p.Run()
	}
}

// ExpandUserHome can replace "~" to user home directory
func ExpandUserHome(path string) (string, error) {
	if !strings.HasPrefix(path, "~/") && path != "~" {
		return path, nil
	}

	user, err := user.Current()
	if err != nil {
		return path, err
	}
	newPath := strings.Replace(path, "~", user.HomeDir, 1)
	return newPath, nil
}
