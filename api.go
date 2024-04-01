package main

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var (
	waitGroup  sync.WaitGroup
	mutex      sync.Mutex
	proxyPools []string = []string{
		"https://openproxylist.xyz/http.txt",
		"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt",
		"https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/http.txt",
		"https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/https.txt",
		"https://raw.githubusercontent.com/MuRongPIG/Proxy-Master/main/http.txt",
		"https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/main/proxy_files/http_proxies.txt",
		"https://raw.githubusercontent.com/andigwandi/free-proxy/main/proxy_list.txt",
		"https://raw.githubusercontent.com/roosterkid/openproxylist/main/HTTPS_RAW.txt",
		"https://raw.githubusercontent.com/mmpx12/proxy-list/master/https.txt",
		"https://raw.githubusercontent.com/mmpx12/proxy-list/master/http.txt",
		"https://raw.githubusercontent.com/zloi-user/hideip.me/main/connect.txt",
		"https://raw.githubusercontent.com/zloi-user/hideip.me/main/http.txt",
		"https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/http/data.txt",
		"https://raw.githubusercontent.com/zloi-user/hideip.me/main/https.txt",
	}
	length    int = len(proxyPools)
	proxyList []string
)

func parseProxies(proxies string) []string {
	proxyPattern := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}`)
	proxyList := proxyPattern.FindAllString(proxies, -1)
	return proxyList
}

func deleteDuplicateProxies(proxies []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range proxies {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func getProxyContent(proxyPoolURL string, waitGroup *sync.WaitGroup, mutex *sync.Mutex) {
	defer waitGroup.Done()
	response, err := http.Get(proxyPoolURL)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return
	}
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	content := string(responseBytes)
	parsedList := parseProxies(content)
	if len(parsedList) == 0 {
		return
	}
	mutex.Lock()
	proxyList = append(proxyList, parsedList...)
	mutex.Unlock()
}

func apiHandler(response http.ResponseWriter, request *http.Request) {
	waitGroup.Add(length)
	for index := range proxyPools {
		proxyPoolURL := proxyPools[index]
		go getProxyContent(proxyPoolURL, &waitGroup, &mutex)
	}
	waitGroup.Wait()
	proxyListString := strings.Join(deleteDuplicateProxies(proxyList), "\n")
	responseContent := []byte(proxyListString)
	response.Write(responseContent)
}

func main() {
	http.HandleFunc("/api", apiHandler)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}
