package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	feedUrl      = "https://thorsten-hans.com/index.xml"
	templatePath = "./template.md"
	readmePath   = "../README.md"
	maxPosts     = 10
	today        = "today"
	day          = "day"
	days         = "days"
	month        = "month"
	months       = "months"
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
	tpl, err := ioutil.ReadFile(templatePath)
	if err != nil {
		log.Fatalf("Error while reading template file. %v", err)
		os.Exit(1)
	}
	p, err := getRecentPosts(maxPosts)
	if err != nil {
		log.Fatalf("Error while reading recent posts %v", err)
		os.Exit(1)
	}
	readme := Readme{
		Posts:   p,
		Updated: time.Now().Format("Mon, 02 Jan 2006"),
	}

	file, err := os.Create(readmePath)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	t := template.Must(template.New("readme").Parse(string(tpl)))
	if err = t.Execute(file, readme); err != nil {
		log.Fatalf("error processing template: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func getRecentPosts(max int) ([]Post, error) {
	p := gofeed.NewParser()
	feed, err := p.ParseURL(feedUrl)
	if err != nil {
		return nil, err
	}

	posts := make([]Post, max)
	for i := 0; i < max; i++ {
		item := feed.Items[i]
		posts[i] = Post{
			Title: item.Title,
			Link:  item.Link,
			Date:  toRelativeDate(item.Published),
		}
	}
	return posts, nil
}

func toRelativeDate(dateString string) string {
	dt, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", dateString)
	if err != nil {
		log.Fatalf("error parsing article date: %v", err)
	}
	now := time.Now().Unix()
	d := (now - dt.Unix()) / 86400
	m := (now - dt.Unix()) / 2592000

	if d == 0 {
		return today
	}

	date := ""
	if d < 31 {
		if d == 1 {
			date = fmt.Sprintf("%d %s", int(d), day)
		} else {
			date = fmt.Sprintf("%d %s", int(d), days)
		}
	} else {
		if m == 1 {
			date = fmt.Sprintf("%d %s", int(m), month)
		} else {
			date = fmt.Sprintf("%d %s", int(m), months)
		}
	}
	return fmt.Sprintf("%s ago", date)
}
