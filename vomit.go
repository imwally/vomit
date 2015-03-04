package main

import (
	"fmt"
	"log"
	//"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	postdir = "posts"
	sitedir = "site"
)

type Post struct {
	Title   string
	Date    string
	Content []byte
}

func GeneratePostPage(post Post) {
	// Write post html page
}

func GenerateIndexPage() {
	// Write index html page
}

func GetPost(p string) Post {
	var post Post

	basename := strings.TrimSuffix(p, filepath.Ext(p))
	post.Date = basename[:10]
	post.Title = basename[11:]
	content, err := ioutil.ReadFile(p)
	CheckErr(err)

	post.Content = content

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

	err := filepath.Walk(postdir, find)
	CheckErr(err)

	return posts
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {

	// Check for posts directory
	_, err := os.Stat(postdir)
	if err != nil {
		fmt.Println("no posts directory found.")
		return
	}

	// Gather posts
	posts := FindPosts(postdir)

	// Generate post pages
	for _, post := range posts {
		GeneratePostPage(post)
	}

	// Generate index page
}
