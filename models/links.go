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
	ID            int `gorm:"primary_key"`
	Title         string
	Desc          string
	AdLinks       []AdLink       // 广告点击跳转链接
	DownloadLinks []DownloadLink // 真正的下载链接
}

type AdLink struct {
	gorm.Model
	ArticleID int
	PkgIndex  int // 分包序号
	Url       string
}

type DownloadLink struct {
	gorm.Model
	ArticleID int
	PkgIndex  int // 分包序号
	Url       string
}

func init() {
	var article = Article{}
	var adLink = AdLink{}
	var downloadLink = DownloadLink{}
	if !sqliteDb.HasTable(&article) {
		sqliteDb.CreateTable(&article)
	}
	if !sqliteDb.HasTable(&adLink) {
		sqliteDb.CreateTable(&adLink)
		sqliteDb.Model(&article).Related(&adLink)
	}
	if !sqliteDb.HasTable(&downloadLink) {
		sqliteDb.CreateTable(&downloadLink)
		sqliteDb.Model(&article).Related(&downloadLink)
	}
}

func UpdateArticle(articleId int, title, desc string) {
	sqliteDb.Model(Article{ID: articleId}).Update(Article{Title: title, Desc: desc})
}

func AddArticleAdUrl(articleId int, pkgIndex int, adUrl string) {
	article := Article{ID: articleId}
	associton := sqliteDb.Model(&article).Association("AdLinks")
	if associton == nil || associton.Count() == 0 {
		adLinks := make([]AdLink, 1)
		adLinks[0].PkgIndex = pkgIndex
		adLinks[0].Url = adUrl
		article.AdLinks = adLinks
	} else {
		adLinks := make([]AdLink, 0, associton.Count()+1)
		associton.Find(&adLinks)
		adLinks = append(adLinks, AdLink{
			PkgIndex: pkgIndex,
			Url:      adUrl,
		})
		article.AdLinks = adLinks
	}
	sqliteDb.Save(&article)
}

func AddArticleDownloadUrl(articleId int, pkgIndex int, downloadUrl string) {
	article := Article{ID: articleId}
	associton := sqliteDb.Model(&article).Association("DownloadLinks")
	if associton == nil || associton.Count() == 0 {
		downloadLinks := make([]DownloadLink, 1)
		downloadLinks[0].PkgIndex = pkgIndex
		downloadLinks[0].Url = downloadUrl
		article.DownloadLinks = downloadLinks
	} else {
		downloadLinks := make([]DownloadLink, 0, associton.Count()+1)
		associton.Find(&downloadLinks)
		downloadLinks = append(downloadLinks, DownloadLink{
			PkgIndex: pkgIndex,
			Url:      downloadUrl,
		})
		article.DownloadLinks = downloadLinks
	}
	sqliteDb.Save(&article)
}

func GetArticleAdUrls(id int) map[int][]string {
	associton := sqliteDb.Model(&Article{ID: id}).Association("AdLinks")
	if associton == nil || associton.Count() == 0 {
		return nil
	}

	adLinks := make([]AdLink, associton.Count())
	associton.Find(&adLinks)
	if len(adLinks) == 0 {
		return nil
	}

	rspUrls := make(map[int][]string, 0)
	for _, adLink := range adLinks {
		curIndex := adLink.PkgIndex
		if rspUrls[curIndex] == nil {
			rspUrls[curIndex] = make([]string, 0, len(adLinks))
		}
		rspUrls[curIndex] = append(rspUrls[curIndex], adLink.Url)
	}

	return rspUrls
}

func GetArticleDownloadUrls(id int) map[int][]string {
	associton := sqliteDb.Model(&Article{ID: id}).Association("DownloadLinks")
	if associton == nil || associton.Count() == 0 {
		return nil
	}

	downloadLinks := make([]DownloadLink, associton.Count())
	associton.Find(&downloadLinks)
	if len(downloadLinks) == 0 {
		return nil
	}

	rspUrls := make(map[int][]string, 0)
	for _, downloadLink := range downloadLinks {
		curIndex := downloadLink.PkgIndex
		if rspUrls[curIndex] == nil {
			rspUrls[curIndex] = make([]string, 0, len(downloadLinks))
		}
		rspUrls[curIndex] = append(rspUrls[curIndex], downloadLink.Url)
	}

	return rspUrls
}
