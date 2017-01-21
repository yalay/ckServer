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

const (
	kLinkTypeUnknown  = 0
	kLinkTypeAd       = 1
	kLinkTypeDownload = 2
)

type Article struct {
	ID           int `gorm:"primary_key"`
	Title        string
	Desc         string
	AdPkgs       []Pkg // 广告点击跳转链接
	DownloadPkgs []Pkg // 真正的下载链接
}

type Pkg struct {
	gorm.Model
	ArticleID int
	Type      int
	PkgIndex  int // 分包序号
	Links     []Link
}

type Link struct {
	gorm.Model
	PkgID int
	Url   string
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

func UpdateArticleAttr(articleId int, title, desc string) {
	var article = Article{}
	sqliteDb.First(&article, articleId)
	article.Title = title
	article.Desc = desc
	sqliteDb.Save(&article)
}

func InsertArticleUrls(articleId int, adUrls [][]string, downloadUrls [][]string) {
	if len(adUrls) == 0 || len(adUrls) != len(downloadUrls) {
		log.Printf("invalid urls. articleId%s\n", articleId)
		return
	}

	adPkgs := make([]Pkg, len(adUrls))
	for i, urls := range adUrls {
		if len(urls) == 0 {
			log.Printf("empty adUrls. articleId%s\n", articleId)
			continue
		}

		for _, url := range urls {
			if url == "" {
				log.Printf("invalid ad url. articleId%s\n", articleId)
				continue
			}
			var link = Link{
				Url: url,
			}
			adPkgs[i].PkgIndex = i
			adPkgs[i].Type = kLinkTypeAd
			adPkgs[i].Links = append(adPkgs[i].Links, link)
		}
	}

	downloadPkgs := make([]Pkg, len(downloadUrls))
	for i, urls := range downloadUrls {
		if len(urls) == 0 {
			log.Printf("empty download urls. articleId%s\n", articleId)
			continue
		}

		for _, url := range urls {
			if url == "" {
				log.Printf("invalid url. articleId%s\n", articleId)
				continue
			}
			var link = Link{
				Url: url,
			}
			downloadPkgs[i].PkgIndex = i
			downloadPkgs[i].Type = kLinkTypeDownload
			downloadPkgs[i].Links = append(downloadPkgs[i].Links, link)
		}
	}

	var article = Article{
		ID:           articleId,
		AdPkgs:       adPkgs,
		DownloadPkgs: downloadPkgs,
	}
	sqliteDb.Save(&article)
}

func GetArticleAdUrls(id int) map[int][]string {
	var article = Article{
		ID: id,
	}

	adPkgs := make([]Pkg, 0)
	sqliteDb.Model(&article).Association("AdPkgs").Find(&adPkgs)
	if len(adPkgs) == 0 {
		return nil
	}

	rspUrls := make(map[int][]string, len(adPkgs))
	for _, adPkg := range adPkgs {
		if adPkg.Type != kLinkTypeAd {
			continue
		}
		adLinks := make([]Link, 0)
		sqliteDb.Model(&adPkg).Association("Links").Find(&adLinks)
		if len(adLinks) == 0 {
			continue
		}

		adUrls := make([]string, 0, len(adLinks))
		for _, adLink := range adLinks {
			if adLink.Url == "" {
				continue
			}
			adUrls = append(adUrls, adLink.Url)
		}
		rspUrls[adPkg.PkgIndex] = adUrls
	}

	log.Printf("rspUrls:%+v", rspUrls)
	return rspUrls
}
