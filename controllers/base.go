package controllers

import (
	"common"
	"strconv"
	"strings"
)

const (
	kArticleKeyPrefix = "article-"
)

func genArticleKey(articleId int32) string {
	return kArticleKeyPrefix + strconv.Itoa(int(articleId))
}

func getIdFromArticleKey(articleKey string) int32 {
	return common.Atoi32(strings.TrimPrefix(articleKey, kArticleKeyPrefix))
}
