package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"path/filepath"
)

type GrlsReader struct {
	Url string
}

func (t *GrlsReader) reader() {
	p := t.downloadString()
	if p == "" {
		Logging("get empty string", p)
		return
	}
	url := t.extractUrl(p)
	if url == "" {
		Logging("get empty url", p)
		return
	}
	t.downloadArchive(url)

}

func (t *GrlsReader) downloadString() string {
	pageSource := DownloadPage(t.Url)
	return pageSource

}

func (t *GrlsReader) extractUrl(p string) string {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(p))
	if err != nil {
		Logging(err)
		return ""
	}
	aTag := doc.Find("#ctl00_plate_tdzip > a").First()
	if aTag == nil {
		Logging("a tag not found")
		return ""
	}
	href, ok := aTag.Attr("href")
	if !ok {
		Logging("href attr in a tag not found")
		return ""
	}
	return fmt.Sprintf("https://grls.rosminzdrav.ru/%s", href)
}

func (t *GrlsReader) downloadArchive(url string) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	filePath := fmt.Sprintf("%s/%s/%s", dir, DirTemp, "")
}
