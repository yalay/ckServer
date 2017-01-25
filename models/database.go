package models

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
)

type MyDb struct {
	sync.RWMutex
	*gorm.DB
}

type Article struct {
	ID            int32  `gorm:"primary_key"`
	AId           int32  `gorm:"index;unique"` // 对应文章生成页id
	Title         string `gorm:"not null"`
	Desc          string
	Cover         string
	AdLinks       []AdLink       `gorm:"ForeignKey:AId;AssociationForeignKey:AId"` // 广告点击跳转链接
	DownloadLinks []DownloadLink `gorm:"ForeignKey:AId;AssociationForeignKey:AId"` // 真正的下载链接
}

type AdLink struct {
	gorm.Model
	AId      int32
	PkgIndex int32 // 分包序号
	Url      string
}

type DownloadLink struct {
	gorm.Model
	AId      int32
	PkgIndex int32 // 分包序号
	Url      string
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
		myDb.Model(&article).Related(&adLink, "AdLinks")
	}
	if !myDb.HasTable(&downloadLink) {
		myDb.CreateTable(&downloadLink)
		myDb.Model(&article).Related(&downloadLink, "DownloadLinks")
	}

	m.DB = myDb
	return nil
}

func (m *MyDb) AddArticle(articleId int32, title, desc, cover string) {
	m.Lock()
	var article = Article{}
	m.DB.Where("a_id = ?", articleId).Find(&article)
	article.AId = articleId
	article.Title = title
	article.Desc = desc
	article.Cover = cover
	m.DB.Save(&article)
	m.Unlock()
}

func (m *MyDb) AddArticleAdUrl(articleId int32, pkgIndex int32, adUrl string) error {
	m.Lock()
	defer m.Unlock()
	// 是否已经在数据库中
	var count int
	m.DB.Where("url = ?", adUrl).Find(&AdLink{}).Count(&count)
	if count > 0 {
		return fmt.Errorf("%s exist. Please delete it first.", adUrl)
	}

	m.DB.Where("a_id = ?", articleId).Find(&Article{}).Count(&count)
	if count == 0 {
		return fmt.Errorf("%d do not exist.", articleId)
	}

	adLink := AdLink{
		AId:      articleId,
		PkgIndex: pkgIndex,
		Url:      adUrl,
	}

	m.DB.Save(&adLink)
	return nil
}

func (m *MyDb) AddArticleDownloadUrl(articleId int32, pkgIndex int32, downloadUrl string) error {
	m.Lock()
	defer m.Unlock()
	// 是否已经在数据库中
	var count int
	m.DB.Where("url = ?", downloadUrl).Find(&DownloadLink{}).Count(&count)
	if count > 0 {
		return fmt.Errorf("%s exist. Please delete it first.", downloadUrl)
	}

	m.DB.Where("a_id = ?", articleId).Find(&Article{}).Count(&count)
	if count == 0 {
		return fmt.Errorf("%d do not exist.", articleId)
	}

	downloadLink := DownloadLink{
		AId:      articleId,
		PkgIndex: pkgIndex,
		Url:      downloadUrl,
	}

	m.DB.Save(&downloadLink)
	return nil
}

func (m *MyDb) GetArticleAttrs(articleId int32) (string, string, string) {
	m.RLock()
	var article = Article{}
	m.DB.Where("a_id = ?", articleId).First(&article)
	m.RUnlock()
	return article.Title, article.Desc, article.Cover
}

func (m *MyDb) GetArticleAdUrls(articleId int32) map[int32][]string {
	m.RLock()
	defer m.RUnlock()
	associton := m.DB.Where("a_id = ?", articleId).First(&Article{}).Association("AdLinks")
	if associton == nil || associton.Error != nil {
		return nil
	}

	adLinks := make([]AdLink, associton.Count())
	associton.Find(&adLinks)
	if len(adLinks) == 0 {
		return nil
	}

	rspUrls := make(map[int32][]string, 0)
	for _, adLink := range adLinks {
		curIndex := adLink.PkgIndex
		if rspUrls[curIndex] == nil {
			rspUrls[curIndex] = make([]string, 0, len(adLinks))
		}
		rspUrls[curIndex] = append(rspUrls[curIndex], adLink.Url)
	}

	return rspUrls
}

func (m *MyDb) GetArticleDownloadUrls(articleId int32) map[int32][]string {
	m.RLock()
	defer m.RUnlock()
	associton := m.DB.Where("a_id = ?", articleId).First(&Article{}).Association("DownloadLinks")
	if associton == nil || associton.Error != nil {
		return nil
	}

	downloadLinks := make([]DownloadLink, associton.Count())
	associton.Find(&downloadLinks)
	if len(downloadLinks) == 0 {
		return nil
	}
	rspUrls := make(map[int32][]string, 0)
	for _, downloadLink := range downloadLinks {
		curIndex := downloadLink.PkgIndex
		if rspUrls[curIndex] == nil {
			rspUrls[curIndex] = make([]string, 0, len(downloadLinks))
		}
		rspUrls[curIndex] = append(rspUrls[curIndex], downloadLink.Url)
	}

	return rspUrls
}

func (m *MyDb) GetArticlePkgCount(articleId int32) int {
	m.RLock()
	defer m.RUnlock()
	associton := m.DB.Where("a_id = ?", articleId).Find(&Article{}).Association("DownloadLinks")
	if associton == nil || associton.Error != nil {
		return 0
	}
	return associton.Count()
}

func (m *MyDb) DeleteAdUrl(url string) {
	m.Lock()
	m.DB.Where("url = ?", url).Delete(&AdLink{})
	m.Unlock()
}

func (m *MyDb) DeleteDownloadUrl(url string) {
	m.Lock()
	m.DB.Where("url = ?", url).Delete(&DownloadLink{})
	m.Unlock()
}
