package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	feedUrl      = "https://thorsten-hans.com/index.xml"
	templatePath = "./template.md"
	readmePath   = "../README.md"
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
	p := gofeed.NewParser()
	feed, err := p.ParseURL(feedUrl)
	if err != nil {
		log.Fatalf("error getting feed: %v", err)
		os.Exit(1)
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
