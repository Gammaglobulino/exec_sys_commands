package using_socks5_proxy

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestConnectingUsingSocks5(t *testing.T) {
	targetUrl := "https://check.torproject.org"
	torProxy := "socks5://localhost:9051"

	torProxyUrl, err := url.Parse(torProxy)
	if err != nil {
		log.Fatal(err)
	}
	torTransport := &http.Transport{
		Proxy: http.ProxyURL(torProxyUrl),
	}
	client := &http.Client{
		Transport:     torTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 5,
	}
	response, err := client.Get(targetUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}
