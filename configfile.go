// read the configuration from file if it exists

package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const configFile = "~/.lsbookmarksrc"

// try and read from ~/.lsbookmarks.rc
func importFromConfigFile() (string, error) {

	cfile := configFile
	// https://stackoverflow.com/a/43578461
	if cfile[0] == '~' {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		cfile = filepath.Join(usr.HomeDir, cfile[1:])
	}

	f, err := os.ReadFile(cfile)
	if err != nil {
		return "", fmt.Errorf("could not read file %s", err)
	}

	fTrimmed := strings.TrimSpace(string(f))
	if fTrimmed == "" {
		return "", errors.New("the config file is empgy")
	}
	return fTrimmed, nil
}
