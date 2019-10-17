package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/raomuyang/baku/operator"
)

const (
	versionInfo = "1.1"
)

var (
	src        string
	dst        string
	link       bool
	overwrite  bool
	ignore     string
	version    bool
	customCopy string // custom copy command
)

func parse() {
	flag.StringVar(&src, "src", "", "the src path to backup")
	flag.StringVar(&dst, "dst", "", "the dst path")
	flag.BoolVar(&link, "link", false, "create hard link")
	flag.BoolVar(&overwrite, "overwrite", false, "overwrite the existing files")
	flag.StringVar(&ignore, "ignore", "", "ignore file(s) by regex pattern")
	flag.StringVar(&customCopy, "cmd", "", "custom copy command (optional)")
	flag.BoolVar(&version, "v", false, "show version info")
	flag.Parse()
}

func main() {

	parse()

	if version {
		fmt.Printf("baku-%s\n", versionInfo)
		os.Exit(0)
	}

	if len(src) == 0 || len(dst) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	var options []operator.BackupOption
	if overwrite {
		options = append(options, operator.OverwriteOption)
	}
	if len(ignore) > 0 {
		ignoreOpt, err := operator.GetFilterOption(ignore)
		if err != nil {
			fmt.Printf("Failed to compile the regexp (%s): %v\n", ignore, err)
			os.Exit(1)
		}
		options = append(options, ignoreOpt)
	}

	var act operator.CopyAction
	fmt.Println("======================")
	if len(customCopy) != 0 {
		fmt.Printf("Copy command: %s <src path> <dst path>\n", customCopy)
		act = operator.GetCustomCopyAction(customCopy)
	} else if link {
		fmt.Printf("Copy command: builtin hard link\n")
		act = operator.CreateLinkAction
	} else {
		fmt.Printf("Copy command: builtin copy file\n")
		act = operator.CopyFileAction
	}
	fmt.Printf("======================\n\n")

	expand := func(p string) string {
		p, err := operator.ExpandUserHome(p)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		return p
	}

	src = expand(src)
	dst = expand(dst)

	err := operator.BackupDirectory(src, dst, act, options...)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
