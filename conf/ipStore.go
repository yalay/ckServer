package conf

import (
	"net/url"
)

var gIpStore *ipStore

type ipStore struct {
	// ip->urlType->count
	ipCounts map[string]map[string]int
}

func init() {
	gIpStore = &ipStore{
		ipCounts: make(map[string]map[string]int, 0),
	}
}

// 获取点击量最少的链接
func GetCkLeastUrl(ip string, ckUrls []string) string {
	return gIpStore.getCkLeastUrl(ip, ckUrls)
}

func (s *ipStore) getCount(ip, host string) int {
	if counts, ok := s.ipCounts[ip]; !ok {
		return 0
	} else {
		if num, ok := counts[host]; !ok {
			return 0
		} else {
			return num
		}
	}
}

func (s *ipStore) counts(ip, ckUrl string) {
	ckParsedUrl, err := url.Parse(ckUrl)
	if err != nil {
		return
	}

	host := ckParsedUrl.Host
	if counts, ok := s.ipCounts[ip]; !ok {
		s.ipCounts[ip] = make(map[string]int, 3)
		s.ipCounts[ip][host] = 1
	} else {
		if curNum, ok := counts[host]; !ok {
			counts[host] = 1
		} else {
			counts[host] = curNum + 1
		}
	}
}

func (s *ipStore) getCkLeastUrl(ip string, ckUrls []string) string {
	if len(ckUrls) == 0 {
		return ""
	}

	if _, ok := s.ipCounts[ip]; !ok {
		s.counts(ip, ckUrls[0])
		return ckUrls[0]
	}

	var minCountsUrl string
	var minCounts int
	for _, ckUrl := range ckUrls {
		ckParsedUrl, err := url.Parse(ckUrl)
		if err != nil {
			continue
		}

		host := ckParsedUrl.Host
		num := s.getCount(ip, host)
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

	s.counts(ip, minCountsUrl)
	return minCountsUrl
}
