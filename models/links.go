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
	ID            int `gorm:"primary_key"`
	Title         string
	Desc          string
	AdLinks       []AdLink       // 广告点击跳转链接
	DownloadLinks []DownloadLink // 真正的下载链接
}

type AdLink struct {
	gorm.Model
	ArticleID int
	PkgIndex  int    // 分包序号
	Url       string `gorm:"not null;unique"`
}

type DownloadLink struct {
	gorm.Model
	ArticleID int
	PkgIndex  int    // 分包序号
	Url       string `gorm:"not null;unique"`
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

func AddArticle(articleId int, title, desc string) {
	// 存在则更新
	sqliteDb.Model(Article{ID: articleId}).Updates(Article{Title: title, Desc: desc})
}

func AddArticleAdUrl(articleId int, pkgIndex int, adUrl string) {
	// 是否已经在数据库中
	var count int
	sqliteDb.Model(&AdLink{}).Where("url = ?", adUrl).Count(&count)
	if count > 0 {
		log.Printf("%s exist. Please delete it first.\n", adUrl)
		return
	}
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
	// 是否已经在数据库中
	var count int
	sqliteDb.Model(&DownloadLink{}).Where("url = ?", downloadUrl).Count(&count)
	if count > 0 {
		log.Printf("%s exist. Please delete it first.\n", downloadUrl)
		return
	}

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

func DeleteAdUrl(url string) {
	sqliteDb.Where("url = ?", url).Delete(&AdLink{})
}

func DeleteDownloadUrl(url string) {
	sqliteDb.Where("url = ?", url).Delete(&DownloadLink{})
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
