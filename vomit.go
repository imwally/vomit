package main

import (
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Add(filename string) {

	// Check if filename exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println(err)
		return
	}

	// Set title
	title := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Set Path
	path := "posts/" + title

	// Create path
	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println(err)
	}

	// Convert
	markup, err := DownToUp(filename)
	if err != nil {
		fmt.Println(err)
	}

	// Write markup
	newpost := path + "/index.html"
	err = ioutil.WriteFile(newpost, markup, 0755)
	if err != nil {
		fmt.Println(err)
	}

}

func DownToUp(fn string) ([]byte, error) {
	file, err := ioutil.ReadFile(fn)
	if err != nil {
		return file, err
	}

	markup := blackfriday.MarkdownBasic(file)

	return markup, err
}

func main() {

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "No command given.\n")
	}

	if flag.Arg(0) == "add" {
		Add(flag.Arg(1))
	}
}
