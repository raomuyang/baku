package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/raomuyang/baku/operator"
)

var (
	src       string
	dst       string
	link      bool
	overwrite bool
	ignore    string
)

func parse() {
	flag.StringVar(&src, "src", "", "the src path to backup")
	flag.StringVar(&dst, "dst", "", "the dst path")
	flag.BoolVar(&link, "link", false, "create hard link")
	flag.BoolVar(&overwrite, "overwrite", false, "overwrite the existing files")
	flag.StringVar(&ignore, "ignore", "", "ignore file(s) by regex pattern")
	flag.Parse()
}

func main() {

	parse()

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

	act := operator.CopyFileAction
	if link {
		act = operator.CreateLinkAction
	}

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
