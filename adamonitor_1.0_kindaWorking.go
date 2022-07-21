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
	"strconv"
	"time"

	"strings"
	"syscall"

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

type adafruit struct {
	Context     string   `json:"@context"`
	Type        string   `json:"@type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Sku         int      `json:"sku"`
	Image       []string `json:"image"`
	Offers      struct {
		Type          string `json:"@type"`
		Availability  string `json:"availability"`
		Price         int    `json:"price"`
		PriceCurrency string `json:"priceCurrency"`
		ItemCondition string `json:"itemCondition"`
		URL           string `json:"url"`
		Description   string `json:"description"`
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.

	channelID := "844802341148426260" // Test channel

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "!adafruit_start") {
		for {
			price0, name0, availability0, stockn0 := sendrequest("https://www.adafruit.com/product/4295")
			time.Sleep(1 * time.Second)
			price1, name1, availability1, stockn1 := sendrequest("https://www.adafruit.com/product/4292")
			time.Sleep(1 * time.Second)
			price2, name2, availability2, stockn2 := sendrequest("https://www.adafruit.com/product/4296")
			time.Sleep(1 * time.Second)
			price3, name3, availability3, stockn3 := sendrequest("https://www.adafruit.com/product/4564")
			if availability0 == `http://schema.org/OutOfStock` {
				fmt.Printf("--REPORT--\nItem:\t%v\nStatus:\t%v\n", name0, availability0)
				printTime()
			} else if availability0 == `http://schema.org/InStock` {
				stocks := fmt.Sprintf("%v", stockn0)
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
							Value:  name0,
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   "Link",
							Value:  "https://www.adafruit.com/product/4295",
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   "Stock Left",
							Value:  stocks,
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   "Price",
							Value:  "$" + price0,
							Inline: true,
						},
						&discordgo.MessageEmbedField{
							Name:   "Status",
							Value:  "IN STOCK",
							Inline: true,
						},
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&942283362054320178>")
			} else {
				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "@yamsbot",
					},
					Color: 0x92cded,
					Title: "Error Occured",
					Fields: []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:   "Log",
							Value:  "An error occured and the application was terminated",
							Inline: false,
						},
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.ChannelMessageSendEmbed(channelID, embed)
				os.Exit(1)
			}
			if availability1 == `http://schema.org/OutOfStock` {
				fmt.Printf("--REPORT--\nItem:\t%v\nStatus:\t%v\n", name1, availability1)
				printTime()
			} else if availability1 == `http://schema.org/InStock` {
				stocks := fmt.Sprintf("%v", stockn1)
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
							Value:  name1,
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   "Link",
							Value:  "https://www.adafruit.com/product/4292",
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   "Stock Left",
							Value:  stocks,
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   "Price",
							Value:  "$" + price1,
							Inline: true,
						},
						&discordgo.MessageEmbedField{
							Name:   "Status",
							Value:  "IN STOCK",
							Inline: true,
						},
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&942283362054320178>")
			} else {
				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "@yamsbot",
					},
					Color: 0x92cded,
					Title: "Error Occured",
					Fields: []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:   "Log",
							Value:  "An error occured and the application was terminated",
							Inline: false,
						},
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.ChannelMessageSendEmbed(channelID, embed)
				os.Exit(1)
			}
			if availability2 == `http://schema.org/OutOfStock` {
				fmt.Printf("--REPORT--\nItem:\t%v\nStatus:\t%v\n", name2, availability2)
				printTime()
			} else if availability2 == `http://schema.org/InStock` {
				stocks := fmt.Sprintf("%v", stockn2)
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
							Value:  name2,
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   "Link",
							Value:  "https://www.adafruit.com/product/4296",
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   "Stock Left",
							Value:  stocks,
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   "Price",
							Value:  "$" + price2,
							Inline: true,
						},
						&discordgo.MessageEmbedField{
							Name:   "Status",
							Value:  "IN STOCK",
							Inline: true,
						},
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&942283362054320178>")
			} else {
				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "@yamsbot",
					},
					Color: 0x92cded,
					Title: "Error Occured",
					Fields: []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:   "Log",
							Value:  "An error occured and the application was terminated",
							Inline: false,
						},
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.ChannelMessageSendEmbed(channelID, embed)
				os.Exit(1)
			}
			if availability3 == `http://schema.org/OutOfStock` {
				fmt.Printf("--REPORT--\nItem:\t%v\nStatus:\t%v\n", name3, availability3)
				printTime()
			} else if availability3 == `http://schema.org/InStock` {
				stocks := fmt.Sprintf("%v", stockn3)
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
							Value:  name3,
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   "Link",
							Value:  "https://www.adafruit.com/product/4564",
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   "Stock Left",
							Value:  stocks,
							Inline: false,
						},
						&discordgo.MessageEmbedField{
							Name:   "Price",
							Value:  "$" + price3,
							Inline: true,
						},
						&discordgo.MessageEmbedField{
							Name:   "Status",
							Value:  "IN STOCK",
							Inline: true,
						},
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.ChannelMessageSendEmbed(channelID, embed)
				s.ChannelMessageSend(channelID, "<@&942283362054320178>")
			} else {
				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "@yamsbot",
					},
					Color: 0x92cded,
					Title: "Error Occured",
					Fields: []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:   "Log",
							Value:  "An error occured and the application was terminated",
							Inline: false,
						},
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.ChannelMessageSendEmbed(channelID, embed)
				os.Exit(1)
			}
			time.Sleep(30 * time.Second)
		}
	}
}

func sendrequest(urlA string) (string, string, string, int) {
	proxy := ""
	proxyparsed := strings.Split(proxy, ":")
	ip, port, user, pass := proxyparsed[0], proxyparsed[1], proxyparsed[2], proxyparsed[3]
	url1, _ := url.Parse("http://" + user + ":" + pass + "@" + ip + ":" + port)
	cl := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(url1),
		},
	}
	resp, err := cl.Get(urlA)
	if err != nil {
		log.Print(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	resp.Body.Close()

	jsonExtractreg, _ := regexp.Compile(`\{"@context".*"\}\}`)
	stockFinder, _ := regexp.Compile(`content="\d+ in stock"`)
	stockNumber, _ := regexp.Compile(`\d+`)
	jsonString := fmt.Sprintf(jsonExtractreg.FindString(string(body)))
	stockStr := fmt.Sprintf(stockFinder.FindString(string(body)))
	stockNumb := fmt.Sprintf(stockNumber.FindString(stockStr))
	jsonByte := []byte(jsonString)
	var stock adafruit
	if err := json.Unmarshal(jsonByte, &stock); err != nil {
		log.Print(err)
	}

	price := stock.Offers.Price
	//fprice, _ := strconv.ParseFloat(price, 64)

	sprice := fmt.Sprintf("%v", price)
	name := stock.Name
	availability := stock.Offers.Availability
	stockInt, _ := strconv.Atoi(stockNumb)
	return sprice, name, availability, stockInt
}

func printTime() {
	dt := time.Now()
	fmt.Printf("Time:\t%v\n\n", dt.Format("01-02-2006 15:04:05"))
}
