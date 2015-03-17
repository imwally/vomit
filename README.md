# vomit

_A repulsive markdown to static html blog generator._

`vomit` feeds on markdown posts and regurgitates a very simple static HTML blog.
The blog is nothing more then a single index page that lists each blog post.

## How to use

Two things are required -- `posts/` and `templates/`.

A working directory would look something like this:

```
blog
|-- posts
|   `-- 2015-03-05-vomit.md
`-- templates
    |-- index.html
    |-- post.html
    `-- style.css
```

Run `vomit`.

````
~/blog$ vomit
```

This will generate the static HTML blog inside a newly created `site` directory.

```
blog
|-- posts
|   `-- 2015-03-05-vomit.md
|-- site
|   |-- 2015-03-05-vomit.html
|   |-- index.html
|   `-- style.css
`-- templates
    |-- index.html
    |-- post.html
    `-- style.css

```

## Posts

Each post must have the file name format `YYYY-MM-DD-some-title.md`. Two
different extensions are permitted, md and markdown.

## Templates

Templates make use of Go's [text/template](http://golang.org/pkg/text/template)
package. You can find examples inside this repo's own `templates` directory.
Both templates are applied to the `Post` struct.

### post.html variables

```
{{ .Title }}
{{ .FormattedDate }}
{{ .Content }}
```

### index.html variables

The index.html template is applied to a slice of `Post`'s. You can range over
them like such:

```
{{ range . }}
    {{ .Title }}
    {{ .FormattedDate }}
    {{ .Content }}
{{ end }}
```
