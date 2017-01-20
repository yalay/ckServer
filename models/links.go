package models

import (
	"github.com/jinzhu/gorm"
	"log"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var sqliteDb *gorm.DB

func init() {
	var err error
	sqliteDb, err = gorm.Open("sqlite3", "ckServer.db")
	if err != nil {
		log.Panicf("open db err:%v\n", err)
	}
}

type Article struct {
	gorm.Model
	Title string
	Desc  string
	Pkgs  []Pkg
}

type Pkg struct {
	gorm.Model
	ArticleID int
	Links     []Link
}

type Link struct {
	gorm.Model
	PkgID       int
	DownloadUrl string
}

func init() {
	var article = Article{}
	var pkg = Pkg{}
	var link = Link{}
	if !sqliteDb.HasTable(&article) {
		sqliteDb.CreateTable(&article)
	}
	if !sqliteDb.HasTable(&pkg) {
		sqliteDb.CreateTable(&pkg)
		sqliteDb.Model(&article).Related(&pkg)
	}
	if !sqliteDb.HasTable(&link) {
		sqliteDb.CreateTable(&link)
		sqliteDb.Model(&pkg).Related(&link)
	}
}

func InsertArticleDownloadUrls(articleId int, downloadUrls [][]string) {
	if len(downloadUrls) == 0 {
		return
	}

	pkgs := make([]Pkg, len(downloadUrls))
	for i, urls := range downloadUrls {
		if len(urls) == 0 {
			continue
		}

		for _, url := range urls {
			if url == "" {
				continue
			}
			var link = Link{
				DownloadUrl: url,
			}
			pkgs[i].Links = append(pkgs[i].Links, link)
		}
	}

	var article = Article{
		ArticleId: articleId,
		Pkgs:      pkgs,
	}
	sqliteDb.Save(&article)
}

func GetArticleDownloadUrls(name string) [][]string {
	var article = Article{}
	sqliteDb.Where("name = ?", name).First(&article)
	log.Printf("article:%+v\n", article)
	if len(article.Pkgs) == 0 {
		return nil
	}

	rspUrls := make([][]string, len(article.Pkgs))
	for i, pkg := range article.Pkgs {
		for _, link := range pkg.Links {
			if link.DownloadUrl == "" {
				continue
			}
			rspUrls[i] = append(rspUrls[i], link.DownloadUrl)
		}
	}
	return rspUrls
}
