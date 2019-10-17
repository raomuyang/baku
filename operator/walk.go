package operator

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	fileExists  = "file already exists"
	fileIgnored = "file ignored"
)

type innerError struct {
	message string
}

func (e innerError) Error() string {
	return e.message
}

// BackupOption is the normal mode of the series of custom options during backup files
type BackupOption func(action CopyAction) CopyAction

// CopyAction is define to express the real action what you wanna do
type CopyAction func(oldPath, newPath string, info os.FileInfo) error

// CreateLinkAction will create an new link to the target path
func CreateLinkAction(oldPath, newPath string, info os.FileInfo) error {
	isDir, err := checkAndCreateDir(oldPath, newPath, info)
	if err != nil {
		return err
	}
	if isDir {
		log.Printf("Operated directory: %s to %s", oldPath, newPath)
		return nil
	}

	log.Printf("Link %s to %s\n", oldPath, newPath)
	_, err = os.Stat(newPath)
	if err == nil {
		return innerError{message: fileExists}
	}
	return os.Link(oldPath, newPath)
}

// CopyFileAction will be copied the file from old path to new path
func CopyFileAction(oldPath, newPath string, info os.FileInfo) error {
	isDir, err := checkAndCreateDir(oldPath, newPath, info)
	if err != nil {
		return err
	}
	if isDir {
		log.Printf("Operated directory: %s to %s", oldPath, newPath)
		return nil
	}
	log.Printf("Copy %s to %s\n", oldPath, newPath)

	_, err = os.Stat(newPath)
	if err == nil {
		return innerError{message: fileExists}
	}
	written, err := copyFile(oldPath, newPath)
	log.Printf("---> written %d\n", written)
	return err
}

// OverwriteOption is an option for backup, it will be delete the target file in the new
func OverwriteOption(action CopyAction) CopyAction {

	newAction := func(oldPath, newPath string, info os.FileInfo) error {
		newFInfo, err := os.Stat(newPath)

		rename := false
		bak := newPath + ".bak"
		if err == nil && !newFInfo.IsDir() {
			rename = true
			err := os.Rename(newPath, bak)
			if err != nil {
				return err
			}
		}
		err = action(oldPath, newPath, info)
		if err != nil && rename {
			return os.Rename(bak, newPath)
		}

		if rename {
			return os.Remove(bak)
		}
		return nil
	}
	return newAction
}

// BackupDirectory backup
func BackupDirectory(root, target string, action CopyAction, options ...BackupOption) error {

	walkFunc := backupFileWalkFunc(root, target, action, options...)
	return filepath.Walk(root, walkFunc)
}

func backupFileWalkFunc(oldPrefix, newPrefix string, action CopyAction, options ...BackupOption) filepath.WalkFunc {

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("[ERROR] walk err: %v\n", err)
		}
		relPath, err := filepath.Rel(oldPrefix, path)
		if err != nil {
			return err
		}
		newPath := filepath.Join(newPrefix, relPath)

		for i := range options {
			action = options[i](action)
		}

		err = action(path, newPath, info)
		if err != nil {
			if err.Error() == fileExists || err.Error() == fileIgnored {
				log.Printf("%s: %s\n", err.Error(), path)
				return nil
			}
		}
		return err

	}
	return walkFunc
}

func copyFile(src, dst string) (written int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	written, err = io.Copy(dstFile, srcFile)
	return
}

func checkAndCreateDir(oldPath, newPath string, info os.FileInfo) (isDir bool, err error) {
	isDir = info.IsDir()
	if isDir {
		f, err := os.Stat(newPath)
		if err == nil {
			if !f.IsDir() {
				return false, fmt.Errorf("%s not a directory", newPath)
			}
			return true, nil
		}
		return true, os.MkdirAll(newPath, 0755)
	}

	realPath, err := filepath.EvalSymlinks(oldPath)
	if err == nil {
		realStat, err := os.Stat(realPath)
		if err != nil {
			return false, err
		}
		if realStat.IsDir() {
			log.Printf("[WARN] unsupported symbol link of directory: %s\n", oldPath)
			return true, nil
		}
	}

	return false, nil
}
