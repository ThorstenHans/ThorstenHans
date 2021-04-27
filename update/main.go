package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	feedUrl  = "https://thorsten-hans.com/index.xml"
	filename = "../README.md"
)

type Readme struct {
	Posts   []Post
	Updated string
}

type Post struct {
	Title string
	Link  string
	Date  string
}

func main() {
	tpl := `## Hi there, I am Thorsten 👋🏼

- 🇩🇪 I am a cloud consultant from Germany 
- 🔷 I am a Microsoft MVP since 2011
- 🐳 I do quite a bunch of Docker
- ☸️ Kubernetes is my passion
- 🌤 Azure is my datacenter

## Recent posts from [my blog](https://thorsten-hans.com) 

{{range .Posts}}- **[{{.Title}}]({{.Link}})** ({{.Date}})
{{end}}
## Get in touch

Reach out via [🐦 Twitter at @ThorstenHans](https://twitter.com/ThorstenHans) or find me on [LinkedIn](https://linkedin.com/in/ThorstenHans).

_last update_: {{ .Updated }}
`

	p := gofeed.NewParser()
	feed, err := p.ParseURL(feedUrl)
	if err != nil {
		log.Fatalf("error getting feed: %v", err)
	}

	var posts []Post
	for i := 0; i < 5; i++ {
		p := feed.Items[i]
		post := Post{
			Title: p.Title,
			Link:  p.Link,
			Date:  relativeDate(p.Published),
		}
		posts = append(posts, post)
	}

	readme := Readme{
		Posts:   posts,
		Updated: time.Now().Format("Mon, 02 Jan 2006"),
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer file.Close()

	t := template.Must(template.New("readme").Parse(tpl))
	if err = t.Execute(file, readme); err != nil {
		log.Fatalf("error processing template: %v", err)
	}
}

func relativeDate(d string) string {
	dt, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", d)
	if err != nil {
		log.Fatalf("error parsing article date: %v", err)
	}
	now := time.Now().Unix()
	days := (now - dt.Unix()) / 86400
	months := (now - dt.Unix()) / 2592000

	if days == 0 {
		return d
	}

	date := ""
	if days < 31 {
		date = strconv.Itoa(int(days))
		if days == 1 {
			date += " day"
		} else {
			date += " days"
		}
	} else {
		date = strconv.Itoa(int(months))
		if months == 1 {
			date += " month"
		} else {
			date += " months"
		}
	}
	return fmt.Sprintf("%s ago", date)
}
