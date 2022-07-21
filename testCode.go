package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type AdafruitJSON struct {
	Context     string   `json:"@context"`
	Type        string   `json:"@type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Sku         int      `json:"sku"`
	Image       []string `json:"image"`
	Offers      struct {
		Type          string          `json:"@type"`
		Availability  string          `json:"availability"`
		Price         json.RawMessage `json:"price"`
		PriceCurrency string          `json:"priceCurrency"`
		ItemCondition string          `json:"itemCondition"`
		URL           string          `json:"url"`
		Description   string          `json:"description"`
	} `json:"offers"`
}

func main() {
	prodURL := "https://www.adafruit.com/product/4564"

	proxyRaw := ""
	proxySplt := strings.Split(proxyRaw, ":")
	ip, port, user, pass := proxySplt[0], proxySplt[1], proxySplt[2], proxySplt[3]
	proxy, err := url.Parse("http://" + user + ":" + pass + "@" + ip + ":" + port)
	if err != nil {
		log.Println(err)
	}
	cl := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	resp, err := cl.Get(prodURL)
	if err != nil {
		log.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	jsonRegex, _ := regexp.Compile(`{"@context".*"@type":"Product".*}}`)
	jsonString := fmt.Sprintf(jsonRegex.FindString(string(body)))
	jsonByte := []byte(jsonString)
	var results AdafruitJSON
	if err := json.Unmarshal(jsonByte, &results); err != nil {
		log.Println(err)
	}
	parseAval, _ := regexp.Compile(`\w+$`)
	aval := fmt.Sprintf(parseAval.FindString(results.Offers.Availability))
	fmt.Println(results.Name)
	fmt.Println(aval)
	fmt.Println(string(results.Offers.Price))
}
