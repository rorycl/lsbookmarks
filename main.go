package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// Options define the flag options
type Options struct {
	Args struct {
		BookmarkBackupPath string `description:"firefox bookmark backup path"`
	} `positional-args:"yes"`
}

const usage = `

list firefox bookmarks with interactive search.
provide the backup path or save the path in ~/.lsbookmarksrc

lsbookmarks`

func main() {
	var options Options
	var parser = flags.NewParser(&options, flags.Default)
	parser.Usage = usage

	if _, err := parser.Parse(); err != nil {
		if !flags.WroteHelp(err) {
			parser.WriteHelp(os.Stdout)
		}
		os.Exit(1)
	}

	if options.Args.BookmarkBackupPath == "" {
		bbp, err := importFromConfigFile()
		if err != nil {
			parser.WriteHelp(os.Stdout)
			fmt.Println("\nfirefox bookmark backup path not provided or config file empty")
			os.Exit(1)
		}
		options.Args.BookmarkBackupPath = bbp
	}

	bookmarks, err := getBmarks(options.Args.BookmarkBackupPath)
	if err != nil {
		fmt.Println("bookmarks could not be extracted: ", err)
		os.Exit(1)
	}
	run(bookmarks)
}
