package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"

	"strconv"
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
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.Contains(m.Content, "!creationTime") {
		r, _ := regexp.Compile("\\d*$")
		userID := fmt.Sprintf(r.FindString(m.Content))
		userID1 := strings.TrimSpace(userID)
		timeT, err := creationTime(userID1)
		//timeS := ""
		if err != nil {
			return
		}
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{},
			Color:  0x92cded, // Red
			Title:  "Account Creation Time",
			//Description: "",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Created at:",
					Value: timeT.String(),
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
	if strings.Contains(m.Content, "!help") {
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "@yamsbot",
			},
			Color: 0x92cded, // Red
			Title: "Help Menu",
			//Description: "",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Sizerun",
					Value: "Syntax: !sizes 6 12",
				},
				&discordgo.MessageEmbedField{
					Name:  "Delay Calculator",
					Value: "Syntax: !delay proxies tasks",
				},
				&discordgo.MessageEmbedField{
					Name:  "Variant Scrapper",
					Value: "Syntax: !variants {shopify url}",
				},
				&discordgo.MessageEmbedField{
					Name:  "Shopify Stock Puller",
					Value: "Syntax: !stock {shopify url}",
				},
				&discordgo.MessageEmbedField{
					Name:  "Account Creation Time",
					Value: "Syntax: !creationTime {discord userID}",
				},
				&discordgo.MessageEmbedField{
					Name:  "Fee Calculator",
					Value: "Syntax: !fee saleprice feePercentage",
				},
				&discordgo.MessageEmbedField{
					Name:  "Avatar request",
					Value: "Syntax: !av @yamsbot",
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
	// If the message is "!sizes" reply with sizerun.
	if strings.Contains(m.Content, "!sizes") {
		// Regex for sizes
		r, _ := regexp.Compile("\\d+ \\d+")
		r2, _ := regexp.Compile(`^\d+`)
		r3, _ := regexp.Compile(`\d+$`)
		// Remove leading or trailing whitespace
		string := strings.TrimSpace(m.Content)
		sizes := fmt.Sprintf(r.FindString(string))
		// Extract sizes
		a := fmt.Sprintf(r2.FindString(sizes))
		b := fmt.Sprintf(r3.FindString(sizes))
		sizeA, errA := strconv.Atoi(a)
		sizeB, errB := strconv.Atoi(b)
		for {
			if errA != nil && errB != nil {
				break
			} else {
				sizeOut := sizeRun(sizeA, sizeB)
				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "@yamsbot",
					},
					Color: 0x92cded, // Red
					Title: "Size Generator",
					//Description: "",
					Fields: []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:  "Size run:",
							Value: sizeOut,
						},
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
				break
			}
		}
	}
	// Variants Scrapper
	if strings.Contains(m.Content, "!variants") {
		// Regex for variants (this code is for sizes, change)
		r, _ := regexp.Compile(`https://[\w\d.-]*\w+\.[comnetggshp]*/[\w\d/-]*`)
		// Remove leading or trailing whitespace
		urlString := strings.TrimSpace(m.Content)
		// Extract URL
		url := fmt.Sprintf(r.FindString(urlString))
		// Call variantRequest function
		varT, vars, varsS := variantRequest(url)
		// Convert int64 slice into String
		varsString := []string{}
		for i := range vars {
			ints := vars[i]
			strng := strconv.Itoa(ints)
			varsString = append(varsString, strng)
		}
		// Join string slice
		fnl := strings.Join(varsString, "\n")
		fnlS := strings.Join(varsS, "\n")
		// Creating message embed
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "@yamsbot",
			},
			Color: 0x92cded,
			Title: varT,
			//Description: "Varaints returned",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "URL",
					Value:  url,
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   "Variants",
					Value:  fnl,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Sizes",
					Value:  fnlS,
					Inline: true,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	// ShoePalace Stock
	if strings.Contains(m.Content, "!stock") {
		// Checks for shoepalace link
		//r, _ := regexp.Compile(`https://www.shoepalace.com/[\w\d/-]*`)
		r, _ := regexp.Compile(`https://[\w\d.-]*\w+\.[comnetggshp]*/[\w\d/-]*`)
		// Remove leading or trailing whitespace
		urlString := strings.TrimSpace(m.Content)
		// Extract URL
		url := fmt.Sprintf(r.FindString(urlString))
		// Call variantRequest function
		varT, vars, varsS, stocknm := spStockRequest(url)
		// Convert int64 slices into Strings
		varsString := []string{}
		stockString := []string{}
		for i := range vars {
			ints := vars[i]
			strng := strconv.Itoa(ints)
			varsString = append(varsString, strng)
		}
		for i := range stocknm {
			ints := stocknm[i]
			strng := strconv.Itoa(ints)
			stockString = append(stockString, strng)
		}
		// Join string slices
		fnl := strings.Join(varsString, "\n")
		fnlS := strings.Join(varsS, "\n")
		spSN := strings.Join(stockString, "\n")

		// Creating message embed
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "@yamsbot",
			},
			Color: 0x92cded,
			Title: varT,
			//Description: "Varaints returned",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "URL",
					Value:  url,
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   "Variants",
					Value:  fnl,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Sizes",
					Value:  fnlS,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Stock",
					Value:  spSN,
					Inline: true,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	// Delay Calculator
	if strings.Contains(m.Content, "!delay") {
		// Regex for proxies & tasks
		r, _ := regexp.Compile("\\d+ \\d+")
		r2, _ := regexp.Compile(`^\d+`)
		r3, _ := regexp.Compile(`\d+$`)
		// Remove leading or trailing whitespace
		string := strings.TrimSpace(m.Content)
		sizes := fmt.Sprintf(r.FindString(string))
		// Extract proxies and delays
		a := fmt.Sprintf(r2.FindString(sizes))
		b := fmt.Sprintf(r3.FindString(sizes))
		proxies, errA := strconv.Atoi(a)
		tasks, errB := strconv.Atoi(b)
		for {
			if errA != nil && errB != nil {
				break
			} else {
				finaldelay := delayCalculator(proxies, tasks)
				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "@yamsbot",
					},
					Color:       0x92cded, // Red
					Title:       "Delay Calculator",
					Description: "Formula: 3600/(proxies/tasks)",
					Fields: []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:  "Delay:",
							Value: finaldelay,
						},
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
				break
			}
		}
	}

	if strings.Contains(m.Content, "!fee") {
		fee64 := 0.0                                             // Initialize float for fee
		extract, _ := regexp.Compile(`(\d+ \d+$|\d+ \d+\.\d+$)`) // Regex to extract sale and fee values
		convStr := fmt.Sprintf(extract.FindString(m.Content))    //  Store sale and fee in convStr
		regSale, _ := regexp.Compile(`^\d+`)                     // Regex to extract Sale value
		regFees, _ := regexp.Compile(`(\d+$|\d+\.\d+$)`)         // Regex to extract Fee value
		sale := fmt.Sprintf(regSale.FindString(convStr))         // Store sale value in variable sale as TYPE string
		fees := fmt.Sprintf(regFees.FindString(convStr))         // Store fee value in variable fees as TYPE string

		if strings.Contains(fees, ".") {
			fee64, _ = strconv.ParseFloat(fees, 64) // Store float value of fees in fee64
			if fee64 < 10.0 {
				//fee := fmt.Sprintf("%.1f", fee64)      // Convert float64 back to string value for editing
				fee := strings.Replace(fees, ".", "", 1) // Replace decimal point with nothing
				fee = fmt.Sprintf(".0%v", fee)
				fee64, _ = strconv.ParseFloat(fee, 64)
			} else if fee64 > 10.0 {
				fee := strings.Replace(fees, ".", "", 1)
				fee = ("." + fee)
				fee64, _ = strconv.ParseFloat(fee, 64)
			}
		} else {
			feeInt, _ := strconv.Atoi(fees)
			if feeInt < 10 {
				fee := fmt.Sprintf(".0%v", feeInt)
				fee64, _ = strconv.ParseFloat(fee, 64)
			} else if feeInt >= 10 {
				fee := fmt.Sprintf(".%v", feeInt)
				fee64, _ = strconv.ParseFloat(fee, 64)
			}
		}
		sale64, _ := strconv.ParseFloat(sale, 64)
		feeCalc := fmt.Sprintf("%.2f", sale64*fee64)

		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{},
			Color:  0x92cded, // Baby blue
			Title:  "Fee Calculator",
			//Description: "",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Sale",
					Value:  "$" + sale,
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Fee percentage",
					Value:  fees + "%",
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:  "Fees due",
					Value: "$" + feeCalc,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	if strings.Contains(m.Content, "!av") {
		var url string
		var uid string
		fmt.Println(m.Content)
		// Regex to parse UID from message
		parseUser, _ := regexp.Compile(`<@!\d+>`)
		parseUID, _ := regexp.Compile(`\d+`)
		// Parse userID
		user := fmt.Sprintf("%v", parseUser.FindString(m.Content))
		if user == "" {
			uid = fmt.Sprintf("%v", parseUID.FindString(m.Content))
		} else {
			uid = fmt.Sprintf("%v", parseUID.FindString(user))
		}
		// Retrieve user information based on userID supplied
		userData, err := s.User(uid)
		if err != nil {
			log.Println(err)
		} else {
			imghash := userData.Avatar
			if userData.Avatar == "" {
				url = discordgo.EndpointDefaultUserAvatar(userData.Discriminator)
			} else if strings.HasPrefix(userData.Avatar, "a_") {
				url = "https://cdn.discordapp.com/avatars/" + uid + "/" + imghash + ".gif?size=4096"
			} else {
				url = "https://cdn.discordapp.com/avatars/" + uid + "/" + imghash + ".png?size=2048"
			}
			// Craft and send embed
			name := fmt.Sprintf("%v#%v", userData.Username, userData.Discriminator)
			embed := &discordgo.MessageEmbed{
				Color: 0x92cded,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    name,
					IconURL: url,
				},
				Title: "Avatar",
				Image: &discordgo.MessageEmbedImage{
					URL:    url,
					Width:  1000,
					Height: 1000,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "@yamsbot",
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		}
	}
}

func sizeRun(a int, b int) string {
	sizes := []float64{}
	finaltxt := ""
	switch {
	case a < 1:
		break
	case a > 16:
		break
	case b < 1:
		break
	case b > 16:
		break
	default:
		for i := float64(a); i <= float64(b); i += 0.5 {
			sizes = append(sizes, i)
		}
		sizeString := []string{}
		for i := range sizes {
			number := sizes[i]
			numtxt := fmt.Sprintf("%g", number)
			sizeString = append(sizeString, numtxt)
		}
		finaltxt = strings.Join(sizeString, ",")
		//return finaltxt
	}
	return finaltxt
}

func delayCalculator(proxies int, tasks int) string {
	delay := ""
	switch {
	case proxies > 10000:
		break
	case proxies < 1:
		break
	case tasks > 10000:
		break
	case tasks < 1:
		break
	case proxies < tasks:
		break
	default:
		delay = fmt.Sprintf("%.f", float64(3600)/(float64(proxies)/float64(tasks)))
	}
	return delay
}

type Response struct {
	Product struct {
		ID             int64       `json:"id"`
		Title          string      `json:"title"`
		BodyHTML       string      `json:"body_html"`
		Vendor         string      `json:"vendor"`
		ProductType    string      `json:"product_type"`
		CreatedAt      string      `json:"created_at"`
		Handle         string      `json:"handle"`
		UpdatedAt      string      `json:"updated_at"`
		PublishedAt    string      `json:"published_at"`
		TemplateSuffix interface{} `json:"template_suffix"`
		PublishedScope string      `json:"published_scope"`
		Tags           string      `json:"tags"`
		Variants       []struct {
			ID                   int         `json:"id"`
			ProductID            int64       `json:"product_id"`
			Title                string      `json:"title"`
			Price                string      `json:"price"`
			Sku                  string      `json:"sku"`
			Position             int         `json:"position"`
			InventoryPolicy      string      `json:"inventory_policy"`
			CompareAtPrice       string      `json:"compare_at_price"`
			FulfillmentService   string      `json:"fulfillment_service"`
			InventoryManagement  string      `json:"inventory_management"`
			Option1              string      `json:"option1"`
			Option2              interface{} `json:"option2"`
			Option3              interface{} `json:"option3"`
			CreatedAt            string      `json:"created_at"`
			UpdatedAt            string      `json:"updated_at"`
			Taxable              bool        `json:"taxable"`
			Barcode              string      `json:"barcode"`
			Grams                int         `json:"grams"`
			ImageID              interface{} `json:"image_id"`
			Weight               float64     `json:"weight"`
			WeightUnit           string      `json:"weight_unit"`
			InventoryQuantity    int         `json:"inventory_quantity"`
			OldInventoryQuantity int         `json:"old_inventory_quantity"`
			TaxCode              string      `json:"tax_code"`
			RequiresShipping     bool        `json:"requires_shipping"`
		} `json:"variants"`
		Options []struct {
			ID        int64    `json:"id"`
			ProductID int64    `json:"product_id"`
			Name      string   `json:"name"`
			Position  int      `json:"position"`
			Values    []string `json:"values"`
		} `json:"options"`
		Images []struct {
			ID         int64         `json:"id"`
			ProductID  int64         `json:"product_id"`
			Position   int           `json:"position"`
			CreatedAt  string        `json:"created_at"`
			UpdatedAt  string        `json:"updated_at"`
			Alt        interface{}   `json:"alt"`
			Width      int           `json:"width"`
			Height     int           `json:"height"`
			Src        string        `json:"src"`
			VariantIds []interface{} `json:"variant_ids"`
		} `json:"images"`
		Image struct {
			ID         int64         `json:"id"`
			ProductID  int64         `json:"product_id"`
			Position   int           `json:"position"`
			CreatedAt  string        `json:"created_at"`
			UpdatedAt  string        `json:"updated_at"`
			Alt        interface{}   `json:"alt"`
			Width      int           `json:"width"`
			Height     int           `json:"height"`
			Src        string        `json:"src"`
			VariantIds []interface{} `json:"variant_ids"`
		} `json:"image"`
	} `json:"product"`
}

func variantRequest(url string) (string, []int, []string) {
	title := ""
	vars := []int{}
	varswsz := []string{}
	// Check to see if URL is good
	switch {
	case url == ".json":
		break
	case url == "":
		break
	default:
		// Execute Request
		url = url + `.json`
		fmt.Println(url) // testing purposes
		resp, err := http.Get(url)
		if err != nil {
			log.Print(err)
		}
		defer resp.Body.Close()
		// Read Response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
		}
		var result Response
		if err := json.Unmarshal(body, &result); err != nil {
			log.Print(err)
		}
		title = fmt.Sprintf(result.Product.Title)
		for _, rec := range result.Product.Variants {
			//fmt.Println(rec.ID, rec.Title)
			vars = append(vars, rec.ID)
			varswsz = append(varswsz, rec.Title)
		}
	}
	return title, vars, varswsz
}

func creationTime(userID string) (t time.Time, err error) {
	i, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return
	}
	timestamp := (i >> 22) + 1420070400000
	t = time.Unix(timestamp/1000, 0)
	return
}

func spStockRequest(url string) (string, []int, []string, []int) {
	title := ""
	vars := []int{}
	varswsz := []string{}
	stck := []int{}
	// Check to see if URL is good
	switch {
	case url == ".json":
		break
	case url == "":
		break
	default:
		// Execute Request
		url = url + `.json`
		fmt.Println(url) // testing purposes
		resp, err := http.Get(url)
		if err != nil {
			log.Print(err)
		}
		defer resp.Body.Close()
		// Read Response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
		}
		var result Response
		if err := json.Unmarshal(body, &result); err != nil {
			log.Print(err)
		}
		title = fmt.Sprintf(result.Product.Title)
		for _, rec := range result.Product.Variants {
			//fmt.Println(rec.ID, rec.Title)
			vars = append(vars, rec.ID)
			varswsz = append(varswsz, rec.Title)
			stck = append(stck, rec.InventoryQuantity)
		}
	}
	return title, vars, varswsz, stck
}
