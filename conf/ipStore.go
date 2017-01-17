package conf

import (
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/go-yaml/yaml"
)

var gConfig = &Config{}

const (
	kShortenAdf = "adf.ly"
	kShortenBoo = "boo.tw"
	kShortenRef = "ref.so"
)

type ipStore struct {
	// ip->urlType->count
	ipCounts map[string]map[string]int
}

func (i *ipStore) getCount(ip, host string) int {
	if counts, ok := i.ipCounts[ip]; !ok {
		return 0
	} else {
		if num, ok := counts[host]; !ok {
			return 0
		} else {
			return num
		}
	}
}

func (i *ipStore) counts(ip, ckUrl string) {
	ckParsedUrl, err := url.Parse(ckUrl)
	if err != nil {
		return
	}

	host := ckParsedUrl.Host
	if counts, ok := i.ipCounts[ip]; !ok {
		i.ipCounts[ip] = make(map[string]int, 3)
		i.ipCounts[ip][host] = 1
	} else {
		if curNum, ok := counts[host]; !ok {
			counts[host] = 1
		} else {
			counts[host] = curNum + 1
		}
	}
}

func (i *ipStore) getNextUrl(ip string, ckUrls []string) string {
	if len(ckUrls) == 0 {
		return ""
	}

	counts, ok := i.ipCounts[ip]
	if !ok {
		i.counts(ip, ckUrls[0])
		return ckUrls[0]
	}

	var minCountsUrl string
	var minCounts int
	for i, ckUrl := range ckUrls {
		ckParsedUrl, err := url.Parse(ckUrl)
		if err != nil {
			continue
		}

		host := ckParsedUrl.Host
		num := i.getCount(ip, host)
		if minCountsUrl == "" {
			minCountsUrl = ckUrl
			minCounts = num
		} else {
			if num < minCounts {
				minCountsUrl = ckUrl
				minCounts = num
			}
		}
	}

	i.counts(ip, minCountsUrl)
	return minCountsUrl
}
