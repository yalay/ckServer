package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type MyDb struct {
	*gorm.DB
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

func NewMyDb() *MyDb {
	return &MyDb{}
}

func (m *MyDb) OpenDataBase(dbType, dbFile string) error {
	myDb, err := gorm.Open(dbType, dbFile)
	if err != nil {
		return err
	}

	var article = Article{}
	var adLink = AdLink{}
	var downloadLink = DownloadLink{}
	if !myDb.HasTable(&article) {
		myDb.CreateTable(&article)
	}
	if !myDb.HasTable(&adLink) {
		myDb.CreateTable(&adLink)
		myDb.Model(&article).Related(&adLink)
	}
	if !myDb.HasTable(&downloadLink) {
		myDb.CreateTable(&downloadLink)
		myDb.Model(&article).Related(&downloadLink)
	}

	m.DB = myDb
	return nil
}

func (m *MyDb) AddArticle(articleId int, title, desc string) {
	m.DB.Model(Article{ID: articleId}).Updates(Article{Title: title, Desc: desc})
}

func (m *MyDb) AddArticleAdUrl(articleId int, pkgIndex int, adUrl string) error {
	// 是否已经在数据库中
	var count int
	m.DB.Model(&AdLink{}).Where("url = ?", adUrl).Count(&count)
	if count > 0 {
		return fmt.Errorf("%s exist. Please delete it first.\n", adUrl)
	}

	article := Article{ID: articleId}
	associton := m.DB.Model(&article).Association("AdLinks")
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
	m.DB.Save(&article)
	return nil
}

func (m *MyDb) AddArticleDownloadUrl(articleId int, pkgIndex int, downloadUrl string) error {
	// 是否已经在数据库中
	var count int
	m.DB.Model(&DownloadLink{}).Where("url = ?", downloadUrl).Count(&count)
	if count > 0 {
		return fmt.Errorf("%s exist. Please delete it first.\n", downloadUrl)
	}

	article := Article{ID: articleId}
	associton := m.DB.Model(&article).Association("DownloadLinks")
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
	m.DB.Save(&article)
	return nil
}

func (m *MyDb) GetArticleAdUrls(id int) map[int][]string {
	associton := m.DB.Model(&Article{ID: id}).Association("AdLinks")
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

func (m *MyDb) GetArticleDownloadUrls(id int) map[int][]string {
	associton := m.DB.Model(&Article{ID: id}).Association("DownloadLinks")
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

func (m *MyDb) DeleteAdUrl(url string) {
	m.DB.Unscoped().Where("url = ?", url).Delete(&AdLink{})
}

func (m *MyDb) DeleteDownloadUrl(url string) {
	m.DB.Unscoped().Where("url = ?", url).Delete(&DownloadLink{})
}
