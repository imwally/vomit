package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Post struct {
	Filename string
	Title    string
	Date     string
	Content  string
}

const (
	postDir     = "posts"
	templateDir = "templates"
)

func GeneratePostPage(post Post) {
	f, err := os.Create("site/" + post.Filename)
	CheckErr(err)

	t, _ := template.ParseFiles("templates/post.html")
	t.Execute(f, post)
}

func GenerateIndexPage() {
	// Write index html page
}

func GetPost(p string) Post {
	var post Post

	basename := filepath.Base(p)
	basename = strings.TrimSuffix(basename, filepath.Ext(basename))

	content, err := ioutil.ReadFile(p)
	CheckErr(err)

	post.Filename = basename + ".html"
	post.Title = basename[11:]
	post.Date = basename[:10]
	post.Content = string(blackfriday.MarkdownCommon(content))

	return post
}

func FindPosts(p string) []Post {
	var posts []Post

	find := func(p string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			posts = append(posts, GetPost(p))
		}
		return nil
	}

	err := filepath.Walk(postDir, find)
	CheckErr(err)

	return posts
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {

	// Check for needed directories
	dirs := []string{postDir, templateDir}

	for _, dir := range dirs {
		_, err := os.Stat(dir)
		if err != nil {
			fmt.Printf("no %s directory found.\n", dir)
			return
		}
	}

	// Gather posts
	posts := FindPosts(postDir)

	// Generate post pages
	for _, post := range posts {
		GeneratePostPage(post)
	}

}
