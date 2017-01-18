package controllers

import (
	"strings"
)

const (
	kArticleKeyPrefix = "article-"
)

func genArticleKey(articleId string) string {
	return kArticleKeyPrefix + articleId
}

func getIdFromArticleKey(articleKey string) string {
	return strings.TrimPrefix(articleKey, kArticleKeyPrefix)
}
