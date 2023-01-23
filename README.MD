# Font Subset Creator
Creates optimized web fonts (SVG, WOFF and WOFF2) out of a subset of unicode characters from a OpenType or TTF font file.

## Pre-requisites
* [Go](https://go.dev/) programming language
* Python 3.7 or greater (It is recommended to manage your Python versions with [pyenv](https://github.com/pyenv/pyenv))
* [fonttools](https://pypi.org/project/fonttools/)
* [brotli](https://github.com/google/brotli) to compress WOFF2 files

## Usage
In the terminal of your choice, run the script `create-subset-font.go`. It accepts three arguments:
+ A file with a list of glyph names to extract from the font file. The glyph names must correspond to the glyph names in the SVG font file, in order to extract the glyphs from the font file.  
  The items in the list can be separated by a comma (,) or by a return carriage.
+ The source SVG font file.
+ The source OpenType or TTF font file from where to extract the subset of glyphs.

After running the command, it will create the web font files (SVG, WOFF and WOFF2) in the subset/ folder.

Example:
```
go run create-subset-font.go icons.txt fa-brands.svg fa-brands.ttf
```

## Testing
In the test folder you will find example TXT files and font files to generate font subsets from [Font Awesome](https://fontawesome.com/).

```
go run create-subset-font.go test/icons.txt test/fonts/fa-light-300.svg test/fonts/fa-light-300.ttf
go run create-subset-font.go test/icons.txt test/fonts/fa-regular-400.svg test/fonts/fa-regular-400.ttf
go run create-subset-font.go test/icons.txt test/fonts/fa-solid-900.svg test/fonts/fa-solid-900.ttf
go run create-subset-font.go test/icon-brands.txt test/fonts/fa-brands-400.svg test/fonts/fa-brands-400.ttf
```