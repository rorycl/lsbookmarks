package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/frioux/leatherman/pkg/mozlz4"
)

// root is the bookmark json root
type root struct {
	Bookmark
}

// Bookmark is the essential structure of a firefox bookmark
type Bookmark struct {
	GUID         string     `json:"guid"`
	Title        string     `json:"title"`
	Index        int        `json:"index"`
	DateAdded    int64      `json:"dateAdded"`
	LastModified int64      `json:"lastModified"`
	ID           int        `json:"id"`
	Tags         string     `json:"tags"`
	URI          string     `json:"uri"`
	Children     []Bookmark `json:"children,omitempty"`
}

// String returns a string representation of a bookmark
func (b Bookmark) String() string {
	var buf bytes.Buffer
	tpl, err := template.New("tpl").Parse(
		"{{ .Title }}\n{{ if .Tags }}   {{ .Tags }}\n{{ end }}   {{ .URI }}\n",
	)
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(&buf, b)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// getSearchWords gets the words to be used for searching over a
// bookmark
func (b Bookmark) getSearchWords() []string {
	return strings.Split(strings.ToLower(b.Tags+" "+b.Title), " ")
}

// Bookmarks is a slice of Bookmark
type Bookmarks []Bookmark

var bookmarks Bookmarks

// Filter returns a filtered set of Bookmarks based on being provided a
// slice of filter expressions
func (b Bookmarks) Filter(filters []string) (int, Bookmarks) {
	i := 0
	if len(filters) == 0 {
		return i, b
	}
	filters = func() []string {
		var s []string
		for _, f := range filters {
			s = append(s, strings.ToLower(f))
		}
		return s
	}()

	bf := []Bookmark{}
	adder := func(k Bookmark) (ok bool) {
		matches := 0
		for _, filter := range filters {
			for _, word := range k.getSearchWords() {
				if filter == word {
					i++
					matches++
					if matches == len(filters) {
						bf = append(bf, k)
						return true
					}
					break
				}
			}
		}
		return false
	}

	for _, k := range b {
		adder(k)
	}
	return i, bf
}

func listBookmarks(b Bookmark) {
	if b.URI != "" {
		bookmarks = append(bookmarks, b)
	}
	for _, c := range b.Children {
		listBookmarks(c)
	}
	return
}

// get latest bookmark backup file
func getLatestBookmarkBackup(path string) (string, error) {
	var filename string
	var err error
	files, err := os.ReadDir(path)
	if err != nil {
		return filename, err
	}
	// sort.Slice(files, func(i, j int) bool { return files[i].Info.ModTime() > files[j].Info.ModTime() })
	var file os.DirEntry
	file, err = func() (de os.DirEntry, err error) {
		var latestFile os.DirEntry
		var latestFileInfo fs.FileInfo
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			if !strings.Contains(f.Name(), "bookmarks-") {
				continue
			}
			if !strings.HasSuffix(f.Name(), "jsonlz4") {
				continue
			}
			info, err := f.Info()
			if err != nil {
				return latestFile, err
			}
			if latestFile == nil || info.ModTime().After(latestFileInfo.ModTime()) {
				latestFile = f
				latestFileInfo, err = f.Info()
				if err != nil {
					return latestFile, err
				}
			}
		}
		return latestFile, nil
	}()
	return file.Name(), nil
}

func extractBookmarks(path string) (Bookmarks, error) {

	file, err := getLatestBookmarkBackup(path)
	if err != nil {
		return bookmarks, fmt.Errorf("error finding backup file: %s", err)
	}

	f, err := os.Open(filepath.Join(path, file))
	if err != nil {
		return bookmarks, fmt.Errorf("could not open file %s: %s", file, err)
	}
	defer f.Close()

	jsonReader, err := mozlz4.NewReader(f)
	if err != nil {
		return bookmarks, fmt.Errorf("could not read lz4 file %s: %s", file, err)
	}

	j, err := io.ReadAll(jsonReader)
	if err != nil {
		return bookmarks, fmt.Errorf("error reading backup file %s: %s", file, err)
	}

	r := root{}
	err = json.Unmarshal(j, &r)
	if err != nil {
		return bookmarks, fmt.Errorf("json decode error on file %s: %s", file, err)
	}
	listBookmarks(r.Bookmark)

	return bookmarks, nil
}
