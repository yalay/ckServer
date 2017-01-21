package controllers

import (
	"common"
	"strconv"
	"strings"
)

const (
	kArticleKeyPrefix = "article-"
)

func genArticleKey(articleId int) string {
	return kArticleKeyPrefix + strconv.Itoa(articleId)
}

func getIdFromArticleKey(articleKey string) int {
	return common.Atoi(strings.TrimPrefix(articleKey, kArticleKeyPrefix))
}
