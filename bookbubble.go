// bmark converts firefox bookmarks (from bookmark.go) to a bmark
// struct. The bookmark.Bookmark structure clashes with the bubbletea
// list interface. The tea list.Item interface contracts are also
// fulfilled here.

package main

import (
	"github.com/rorycl/lsbookmarks/bookmark"
)

type bmark struct {
	title string
	uri   string
	tags  string
}

// FilterValue allows a Bookmark to meet the bubbletea Item interface
func (b bmark) FilterValue() string {
	return b.tags + " " + b.title
}

// Title returns the title for bubbletea
func (b bmark) Title() string {
	return b.title
}

// URI returns the title for bubbletea
func (b bmark) URI() string {
	return b.uri
}

// Description returns the description for bubbletea
func (b bmark) Description() string {
	if b.tags == "" {
		return b.uri
	}
	return b.tags + "\n" + b.uri
}

// Bmarks is a slice of Bookmark
type Bmarks []bmark

func getBmarks(path string) (Bmarks, error) {
	var bmarks Bmarks
	bookmarks, err := bookmark.ExtractBookmarks(path)
	if err != nil {
		return bmarks, err
	}
	for _, b := range bookmarks {
		bmarks = append(bmarks, bmark{b.Title, b.URI, b.Tags})
	}
	return bmarks, nil
}

/*
func main() {

	path := "/home/rory/.mozilla/firefox/bpaqics3.default-esr/bookmarkbackups"
	bmarks, err := getBmarks(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, b := range bmarks {
		fmt.Println(b)
	}

}
*/
