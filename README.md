# lsbookmarks

List and interactively search firefox bookmarks from the Linux terminal
using the wonderful
[bubbletea](https://github.com/charmbracelet/bubbletea) TUI library.
Hitting `enter` copies the selected url to the clipboard.

	Usage:
	  lsbookmarks 

	list firefox bookmarks with interactive search.
	provide the backup path or save the path in ~/.lsbookmarksrc

	lsbookmarks [BookmarkBackupPath]

	Help Options:
	  -h, --help                Show this help message

	Arguments:
	  BookmarkBackupPath:       firefox bookmark backup path

## Example

![lsbookmarks](lsbookmarks.gif)

## Requirements

1. Linux
2. Providing the path to the Firefox backup directory either on the
   command line or in a `~/.lsbookmarksrc` file
3. The xlib requirements set out in `golang.design/x/clipboard`,
   currently `libx11-dev`, `xorg-dev` or `libX11-devel`.

## License

This project is licensed under the [MIT Licence](LICENCE).

Rory Campbell-Lange 25 September 2021
