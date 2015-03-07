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
	"time"
)

// Post is a struct that holds information about each blog post.
type Post struct {
	Filename string
	Title    string
	Date     string
	Content  string
}

// Index is a struct that holds all posts.
type Index struct {
	Posts []Post
}

const (
	postDir     = "posts"
	templateDir = "templates"
)

// GeneratePostPage takes a Post and generates an HTML page.
func GeneratePostPage(post Post) {
	f, err := os.Create("site/" + post.Filename)
	CheckErr(err)

	t, _ := template.ParseFiles("templates/post.html")
	t.Execute(f, post)
}

func GenerateIndexPage(index Index) {
	f, err := os.Create("site/index.html")
	CheckErr(err)

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(f, index)
}

// GetPost takes a path to a post and gathers the Filename, Title, Date, and
// content of the post. It returns a Post.
func GetPost(p string) Post {
	var post Post

	basename := filepath.Base(p)
	basename = strings.TrimSuffix(basename, filepath.Ext(basename))

	content, err := ioutil.ReadFile(p)
	CheckErr(err)

	date, err := time.Parse("2006-01-02", basename[:10])
	CheckErr(err)

	post.Date = date.Format("January 2, 2006")
	post.Filename = basename + ".html"
	post.Title = basename[11:]
	post.Content = string(blackfriday.MarkdownCommon(content))

	return post
}

// FindPosts takes a path as an argument that will be traversed and searched for
// markdown files. It returns a slice of Posts.
func FindPosts(p string) []Post {
	var posts []Post

	find := func(p string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if filepath.Ext(p) == ".md" {
				posts = append(posts, GetPost(p))
			} else {
				log.Printf("error: %s is not a markdown file\n", p)
			}
		}
		return nil
	}

	err := filepath.Walk(postDir, find)
	CheckErr(err)

	return posts
}

// CheckErr is a helper function that prints errors.
func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {

	// Check for required directories.
	dirs := []string{postDir, templateDir}

	for _, dir := range dirs {
		_, err := os.Stat(dir)
		if err != nil {
			fmt.Printf("no %s directory found.\n", dir)
			return
		}
	}

	// Gather posts.
	posts := FindPosts(postDir)

	// Generate post pages.
	for _, post := range posts {
		GeneratePostPage(post)
	}

	// Generate index page.
	GenerateIndexPage(Index{Posts: posts})
}
