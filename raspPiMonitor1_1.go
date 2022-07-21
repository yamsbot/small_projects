package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

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
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

var proxies = []string{
	"",
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	channelID := ""      // channelID
	monitorErrorID := "" // error channelID
	//
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "!adafruitUFIJG") {
		for {
			embed, availability := sendReq("https://www.adafruit.com/product/4564", proxies[0])
			if availability == "InStock" {
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&942283362054320178>")
			} else if availability == "OutOfStock" {
				log.Printf("8GB OOS\n%v\n", time.Now().Format(time.RFC3339))
			} else {
				s.ChannelMessageSend(monitorErrorID, "There may be an issue with adafruit monitor. <@!287785033983459348>")
			}
			embed, availability = sendReq("https://www.adafruit.com/product/4295", proxies[1])
			if availability == "InStock" {
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&942283362054320178>")
			} else if availability == "OutOfStock" {
				log.Printf("1GB OOS\n%v\n", time.Now().Format(time.RFC3339))
			} else {
				s.ChannelMessageSend(monitorErrorID, "There may be an issue with adafruit monitor. <@!287785033983459348>")
			}
			embed, availability = sendReq("https://www.adafruit.com/product/4292", proxies[2])
			if availability == "InStock" {
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&942283362054320178>")
			} else if availability == "OutOfStock" {
				log.Printf("2GB OOS\n%v\n", time.Now().Format(time.RFC3339))
			} else {
				s.ChannelMessageSend(monitorErrorID, "There may be an issue with adafruit monitor. <@!287785033983459348>")
			}
			embed, availability = sendReq("https://www.adafruit.com/product/4296", proxies[3])
			if availability == "InStock" {
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&942283362054320178>")
			} else if availability == "OutOfStock" {
				log.Printf("4GB OOS\n%v\n", time.Now().Format(time.RFC3339))
			} else {
				s.ChannelMessageSend(monitorErrorID, "There may be an issue with adafruit monitor. <@!287785033983459348>")
			}
			time.Sleep(45 * time.Second)
		}
	}
}

func sendReq(prodURL string, proxyRaw string) (*discordgo.MessageEmbed, string) {
	proxySplt := strings.Split(proxyRaw, ":")
	ip, port, user, pass := proxySplt[0], proxySplt[1], proxySplt[2], proxySplt[3]
	proxy, err := url.Parse("http://" + user + ":" + pass + "@" + ip + ":" + port)
	if err != nil {
		log.Println(err)
	}
	cl := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			IdleConnTimeout: (10 * time.Second),
		},
	}
	req, err := http.NewRequest("GET", prodURL, nil)
	if err != nil {
		log.Println(err)
	}
	req.Close = true
	resp, err := cl.Do(req)
	fmt.Println("StatusCode: ", resp.StatusCode) // Prints status code in terminal after request is sent
	if err != nil {
		log.Println(err)
		log.Println(resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	//cl.CloseIdleConnections()
	jsonRegex, _ := regexp.Compile(`{"@context".*"@type":"Product".*}}`)
	jsonString := fmt.Sprintf(jsonRegex.FindString(string(body)))
	jsonByte := []byte(jsonString)
	var results AdafruitJSON
	if err := json.Unmarshal(jsonByte, &results); err != nil {
		log.Println(err)
	}

	parseAval, _ := regexp.Compile(`\w+$`)
	aval := fmt.Sprintf(parseAval.FindString(results.Offers.Availability))
	parseStockReg, _ := regexp.Compile(`"twitter:data2".*>`)
	stockReg := fmt.Sprintf(parseStockReg.FindString(string(body)))
	parseContent, _ := regexp.Compile(`content=".*"`)
	content := fmt.Sprintf(parseContent.FindString(stockReg))
	parseStock, _ := regexp.Compile(`".*stock"`)
	stockStatus := fmt.Sprintf(parseStock.FindString(content))

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "@yamsbot",
		},
		Color: 0x92cded,
		Title: "Stock Alert",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Item",
				Value:  results.Name,
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Link",
				Value:  prodURL,
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Stock",
				Value:  stockStatus,
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Price",
				Value:  "$" + string(results.Offers.Price),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Status",
				Value:  aval,
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	return embed, aval
}
