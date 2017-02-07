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

func AddArticle(articleId int32, title, sTitle, desc, cover string) {
	sqliteDb.AddArticle(articleId, title, sTitle, desc, cover)
}

func AddArticleAdUrl(articleId int32, pkgIndex int32, adUrl string) error {
	return sqliteDb.AddArticleAdUrl(articleId, pkgIndex, adUrl)
}

func AddArticleDownloadUrl(articleId int32, pkgIndex int32, downloadUrl string) error {
	return sqliteDb.AddArticleDownloadUrl(articleId, pkgIndex, downloadUrl)
}
func GetArticleAttrs(id int32) (string, string, string, string) {
	return sqliteDb.GetArticleAttrs(id)
}

func GetArticleAdUrls(id int32) map[int32][]string {
	return sqliteDb.GetArticleAdUrls(id)
}

func GetArticlePkgCount(id int32) int {
	return sqliteDb.GetArticlePkgCount(id)
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
