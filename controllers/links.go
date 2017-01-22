package controllers

import (
	"log"
	"models"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var sqliteDb *models.MyDb

func init() {
	sqliteDb = models.NewMyDb()
	err := sqliteDb.OpenDataBase("sqlite3", "ckServer.db")
	if err != nil {
		log.Panicf("open db err:%v\n", err)
	}
}

func AddArticle(articleId int32, title, desc string) {
	sqliteDb.AddArticle(articleId, title, desc)
}

func AddArticleAdUrl(articleId int32, pkgIndex int32, adUrl string) {
	err := sqliteDb.AddArticleAdUrl(articleId, pkgIndex, adUrl)
	if err != nil {
		log.Println(err.Error())
	}
}

func AddArticleDownloadUrl(articleId int32, pkgIndex int32, downloadUrl string) {
	err := sqliteDb.AddArticleDownloadUrl(articleId, pkgIndex, downloadUrl)
	if err != nil {
		log.Println(err.Error())
	}
}

func GetArticleAdUrls(id int32) map[int32][]string {
	return sqliteDb.GetArticleAdUrls(id)
}

func GetArticleDownloadUrls(id int32) map[int32][]string {
	return sqliteDb.GetArticleDownloadUrls(id)
}

func DeleteAdUrl(url string) {
	sqliteDb.DeleteAdUrl(url)
}

func DeleteDownloadUrl(url string) {
	sqliteDb.DeleteDownloadUrl(url)
}
