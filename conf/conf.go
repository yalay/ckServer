package conf

import (
	"flag"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"
)

var gConfig = &Config{}

type Config struct {
	WhiteRefers []string
	LinksWeight []string // 按索引排序，越往前权重越大
	Articles    map[int32]*Article
}

type Article struct {
	Title   string
	Cover   string
	Stitle  string
	Desc    string
	AdLinks []string
	DlLinks []string
}

func init() {
	var configFile string
	flag.StringVar(&configFile, "c", "conf/config.yaml", "conf file path")

	reloadConf(configFile)
	rand.Seed(time.Now().Unix())
}

func GetArticleAdLinks(articleId int32) []string {
	if len(gConfig.Articles) == 0 {
		return nil
	}

	article, ok := gConfig.Articles[articleId]
	if !ok || article == nil {
		return nil
	}
	return article.AdLinks
}

func GetArticleAttrs(articleId int32) (string, string, string, string) {
	if len(gConfig.Articles) == 0 {
		return "", "", "", ""
	}

	article, ok := gConfig.Articles[articleId]
	if !ok || article == nil {
		return "", "", "", ""
	}
	return article.Title, article.Stitle, article.Cover, article.Desc
}

func GetArticleStitle(articleId int32) string {
	if len(gConfig.Articles) == 0 {
		return ""
	}

	article, ok := gConfig.Articles[articleId]
	if !ok || article == nil {
		return ""
	}
	return article.Stitle
}

func GetArticleDlLinks(articleId int32) []string {
	if len(gConfig.Articles) == 0 {
		return nil
	}

	article, ok := gConfig.Articles[articleId]
	if !ok || article == nil {
		return nil
	}
	return article.DlLinks
}

func IsInWhiteList(url string) bool {
	if len(gConfig.WhiteRefers) == 0 {
		return true
	}

	for _, domain := range gConfig.WhiteRefers {
		if strings.Contains(url, domain) {
			return true
		}
	}
	return false
}

func GetHighestWeightLink(adLinks []string) string {
	linksLen := len(adLinks)
	if linksLen == 0 {
		return ""
	}

	if linksLen == 1 {
		return adLinks[0]
	}

	if len(gConfig.LinksWeight) == 0 {
		return adLinks[rand.Intn(linksLen-1)]
	}

	hostsLink := make(map[string]string, linksLen)
	for _, adLink := range adLinks {
		linkUrl, err := url.Parse(adLink)
		if err != nil {
			continue
		}
		hostsLink[linkUrl.Host] = adLink
	}

	for _, link := range gConfig.LinksWeight {
		if selectedLink, ok := hostsLink[link]; ok {
			return selectedLink
		}
	}

	return adLinks[rand.Intn(linksLen-1)]
}

func reloadConf(configFile string) {
	err := ParseYaml(configFile, gConfig)
	if err != nil {
		log.Panicf("ParseYaml err:%v\n", err)
	}

	log.Printf("config:%+v", gConfig)
	go reloadYamlFile(configFile, time.Minute, gConfig)
}

func reloadYamlFile(configFile string, duration time.Duration, serverConf *Config) {
	var lastMtime = getFileMtime(configFile)
	for {
		time.Sleep(duration)
		if curMtime := getFileMtime(configFile); curMtime > lastMtime {
			lastMtime = curMtime
			err := ParseYaml(configFile, serverConf)
			if err != nil {
				log.Panicf("ParseYaml err:%v\n", err)
			}
			log.Printf("config:%+v", serverConf)
		}
	}
}

func getFileMtime(file string) int64 {
	fileInfo, err := os.Stat(file)
	if err != nil {
		log.Fatalf("file stat err:%v\n", err)
		return 0
	}
	return fileInfo.ModTime().Unix()
}
