package main

import (
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
	channelID := "" // channelID

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "!start") {
		counter := 0
		for {
			// First item
			if counter == 25 {
				counter = 0
			}
			embed, availability := sendReq("https://store.ui.com/collections/unifi-protect-cameras/products/unifi-protect-ai-360", proxies[counter])
			counter += 1
			if availability == "false" {
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&995422893158703196>")
			} else if availability == "true" {
				fmt.Println("out of stock")
			} else {
				s.ChannelMessageSend(channelID, "There may be an error with the monitor.")
			}
			// Second item
			if counter == 25 {
				counter = 0
			}
			embed, availability = sendReq("https://store.ui.com/collections/unifi-protect-cameras/products/uvc-g4-dome", proxies[counter])
			counter += 1
			if availability == "false" {
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&995422893158703196>")
			} else if availability == "true" {
				fmt.Println("out of stock")
			} else {
				s.ChannelMessageSend(channelID, "There may be an error with the monitor.")
			}
			// Third item
			if counter == 25 {
				counter = 0
			}
			embed, availability = sendReq("https://store.ui.com/collections/unifi-protect-cameras/products/camera-g4-instant", proxies[counter])
			counter += 1
			if availability == "false" {
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&995422893158703196>")
			} else if availability == "true" {
				fmt.Println("out of stock")
			} else {
				s.ChannelMessageSend(channelID, "There may be an error with the monitor.")
			}
			// Fourth item
			if counter == 25 {
				counter = 0
			}
			embed, availability = sendReq("https://store.ui.com/collections/unifi-protect-nvr/products/unvr-pro", proxies[counter])
			counter += 1
			if availability == "false" {
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&995422893158703196>")
			} else if availability == "true" {
				fmt.Println("out of stock")
			} else {
				s.ChannelMessageSend(channelID, "There may be an error with the monitor.")
			}
			// Fifth item
			if counter == 25 {
				counter = 0
			}
			embed, availability = sendReq("https://store.ui.com/collections/unifi-protect-accessories/products/unifi-protect-viewport", proxies[counter])
			counter += 1
			if availability == "false" {
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&995422893158703196>")
			} else if availability == "true" {
				fmt.Println("out of stock")
			} else {
				s.ChannelMessageSend(channelID, "There may be an error with the monitor.")
			}
			time.Sleep(240 * time.Second)
		}
	}
}

func sendReq(produrl string, proxyRaw string) (*discordgo.MessageEmbed, string) {
	// Proxy magic
	proxySplt := strings.Split(proxyRaw, ":")
	ip, port, user, pass := proxySplt[0], proxySplt[1], proxySplt[2], proxySplt[3]
	proxy, err := url.Parse("http://" + user + ":" + pass + "@" + ip + ":" + port)
	if err != nil {
		log.Println(err)
	}
	// Create our HTTP client
	cl := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			IdleConnTimeout: (10 * time.Second),
		},
	}
	// Create our HTTP Request
	req, err := http.NewRequest("GET", produrl, nil)
	if err != nil {
		log.Println(err)
	}
	req.Close = true
	// Send HTTP request and store response in var resp
	resp, err := cl.Do(req)
	if err != nil {
		log.Println(err)
		log.Println(resp.StatusCode)
	}
	// Defer response body close to prevent memory leakage
	defer resp.Body.Close()
	// Close idle connections to prevent memory leakage
	cl.CloseIdleConnections()
	// Print status code to terminal
	fmt.Println(resp.StatusCode)
	//
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	// Regex to parse availability string
	parseAvailability, _ := regexp.Compile(`APP\_DATA\.product\.variants\[0\]\.inventory\_empty = .*\;`)
	availability := fmt.Sprintf(parseAvailability.FindString(string(body)))
	// Regex to parse TRUE/FALSE from string
	parsetf, _ := regexp.Compile(`\w+;$`)
	tf := fmt.Sprintf(parsetf.FindString(availability))
	// Regex to clean up t/f string
	parsetf2, _ := regexp.Compile(`\w+`)
	resulttf := fmt.Sprintf(parsetf2.FindString(tf))
	// Return result
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
				Value:  "Ubiquity Item",
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Link",
				Value:  produrl,
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	return embed, resulttf
}
