package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/extrame/xls"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"strings"
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
	filePath := filepath.FromSlash(fmt.Sprintf("%s/%s/%s", dir, DirTemp, ArZir))
	err := DownloadFile(filePath, url)
	if err != nil {
		Logging("file was not downloaded, exit", err)
		return
	}
	dirZip := filepath.FromSlash(fmt.Sprintf("%s/%s/", dir, DirTemp))
	err = Unzip(filePath, dirZip)
	if err != nil {
		Logging("file was not unzipped, exit", err)
		return
	}
	files, err := FilePathWalkDir(dirZip)
	if err != nil {
		Logging("filelist return error, exit", err)
		return
	}
	for _, f := range files {
		if strings.HasSuffix(f, "xls") {
			t.extractXlsData(f)
		}
	}
}

func (t *GrlsReader) extractXlsData(nameFile string) {
	defer SaveStack()
	xlFile, err := xls.Open(nameFile, "utf-8")
	if err != nil {
		Logging("error open excel file, exit", err)
		return
	}
	sheet := xlFile.GetSheet(0)
	t.insertToBase(sheet)

}

func (t *GrlsReader) insertToBase(sheet *xls.WorkSheet) {
	db, err := sql.Open("sqlite3", "grls.db")
	if err != nil {
		Logging(err)
		return
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM grls")
	if err != nil {
		Logging(err)
		return
	}
	for r := 3; r <= int(sheet.MaxRow); r++ {
		col := sheet.Row(r)
		mnn := col.Col(0)
		name := col.Col(1)
		form := col.Col(2)
		owner := col.Col(3)
		atx := col.Col(4)
		quantity := col.Col(5)
		maxPrice := strings.ReplaceAll(col.Col(6), ",", ".")
		firstPrice := strings.ReplaceAll(col.Col(7), ",", ".")
		ru := col.Col(8)
		dateReg := col.Col(9)
		code := col.Col(10)
		_, err := db.Exec("INSERT INTO grls (id, mnn, name, form, owner, atx, quantity, max_price, first_price, ru, date_reg, code) VALUES (NULL, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", mnn, name, form, owner, atx, quantity, maxPrice, firstPrice, ru, dateReg, code)
		if err != nil {
			Logging(err)
		}
	}
}
