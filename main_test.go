package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRSS(t *testing.T) {
	in := `<rss version="2.0">
	<channel>
	<title>Hacker News</title>
	<link>https://news.ycombinator.com/</link>
	<description>Links for the intellectually curious, ranked by readers.</description>
	<item><title>something uninteresting</title>
	</item>
	</channel>
	</rss>`

	expect := `<rss version="2.0">
	<channel>
	<title>Hacker News</title>
	<link>https://news.ycombinator.com/</link>
	<description>Links for the intellectually curious, ranked by readers.</description>
	
	</channel>
	</rss>`

	r := strings.NewReader(in)
	w := bytes.Buffer{}
	keyword := "keyword"
	filterFeed(r, &w, keyword)

	assert(t, w.String(), expect)
}

func TestRSSManyItems(t *testing.T) {
	in := `<rss version="2.0">
	<channel>
	<title>Hacker News</title>
	<link>https://news.ycombinator.com/</link>
	<description>Links for the intellectually curious, ranked by readers.</description>
	<item><title>cat</title></item>
	<item><title>cat1</title></item>
	<item><title>Dog</title></item>
	<item><title>dog1</title></item>
	</channel>
	</rss>`

	expect := `<rss version="2.0">
	<channel>
	<title>Hacker News</title>
	<link>https://news.ycombinator.com/</link>
	<description>Links for the intellectually curious, ranked by readers.</description>
	
	
	<item><title>Dog</title></item>
	<item><title>dog1</title></item>
	</channel>
	</rss>`

	r := strings.NewReader(in)
	w := bytes.Buffer{}
	keyword := "dog"
	filterFeed(r, &w, keyword)

	assert(t, w.String(), expect)
}

func assert(t *testing.T, got, expect string) {
	t.Helper()
	if got != expect {
		t.Fatalf("--- got --\n%v\n-- expected --\n%v", got, expect)
	}
}
