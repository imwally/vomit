package main

import (
	"bufio"
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

// Common directories.
const (
	postDir     = "posts/"
	templateDir = "templates/"
	siteDir     = "site/"
)

// GeneratePostPage takes a Post and generates a single HTML blog post page.
func GeneratePostPage(post Post) {
	f, err := os.Create(siteDir + post.Filename)
	CheckErr(err)

	t, _ := template.ParseFiles(templateDir + "post.html")
	t.Execute(f, post)
}

// GenerateIndexPage takes a slice of Posts and generates an index page that
// links to all blog posts.
func GenerateIndexPage(index Index) {
	f, err := os.Create(siteDir + "index.html")
	CheckErr(err)

	t, _ := template.ParseFiles(templateDir + "index.html")
	t.Execute(f, index)
}

// CopyStyleSheet will copy the style.css file from the template directory to
// the site directory.
func CopyStyleSheet() {
	f, err := ioutil.ReadFile(templateDir + "style.css")
	CheckErr(err)

	err = ioutil.WriteFile(siteDir+"style.css", f, 0644)
	CheckErr(err)
}

// GetTitleContent takes a os.File and parses the content for the title and body
// of the post.
func GetTitleContent(f *os.File) (string, []byte) {
	var title string
	var fm, ylen int

	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() == "---" {
			fm++
			ylen += len(s.Text())
		} else {
			if fm < 2 {
				if s.Text()[:6] == "title:" {
					title = s.Text()[6:]
				}
				ylen += len(s.Text())
			}
		}
	}

	content, err := ioutil.ReadFile(f.Name())
	CheckErr(err)

	return title, content[ylen+4:]
}

// GetPost takes a path to a post and gathers the Filename, Title, Date, and
// content of the post. It returns a Post.
func GetPost(p string) Post {
	var post Post

	f, err := os.Open(p)
	CheckErr(err)
	defer f.Close()

	// Get title and content
	title, content := GetTitleContent(f)
	post.Title = title
	post.Content = string(blackfriday.MarkdownCommon(content))

	// Get filename
	basename := filepath.Base(f.Name())
	basename = strings.TrimSuffix(basename, filepath.Ext(basename))
	post.Filename = basename + ".html"

	// Get Date
	date, err := time.Parse("2006-01-02", basename[:10])
	CheckErr(err)
	post.Date = date.Format("January 2, 2006")

	return post
}

// FindMarkDown takes a path as an argument that will be traversed and searched for
// markdown files. It returns a slice of Posts.
func FindMarkDown(p string) []Post {
	var posts []Post

	find := func(p string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if filepath.Ext(p) == ".md" || filepath.Ext(p) == ".markdown" {
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

	// Check for required directories and templates.
	dirs := []string{
		postDir,
		templateDir,
		templateDir + "index.html",
		templateDir + "post.html"}

	for _, dir := range dirs {
		_, err := os.Stat(dir)
		if err != nil {
			fmt.Printf("error: %s not found.\n", dir)
			return
		}
	}

	// Create site directory if it doesn't exist.
	if _, err := os.Stat("site"); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("site", 0775)
		}
	}

	// Gather posts.
	posts := FindMarkDown(postDir)

	// Generate post pages.
	for _, post := range posts {
		GeneratePostPage(post)
	}

	// Generate index page.
	GenerateIndexPage(Index{Posts: posts})

	// Copy over style sheet.
	CopyStyleSheet()
}
