package webScrape

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func scrape(searchURL, serviceName string, webScraper Scraper) {
	fmt.Println("Collecting Initial HTML doc")
	html := getHTML(searchURL)
	defer html.Body.Close()
	doc, error := goquery.NewDocumentFromReader(html.Body)

	if error != nil {
		fmt.Println(error)
	}
	if serviceName == "alla items" {
		itemID, itemName := getItemIDandNameByItemSearchDoc(doc)
		itemAttributes := getItemInfoByItemSearchDoc(itemID)
		webScraper.results = append(webScraper.results, itemName)
		webScraper.results = append(webScraper.results, string(rune(itemID)))
		fmt.Printf("Detected Item Name: %s\n", itemName)
		fmt.Printf("Detected Item ID: %d\n", itemID)

		for i, attribute := range itemAttributes {
			fmt.Printf("Detected Attribute %d: %s\n", i+1, attribute)
			webScraper.results = append(webScraper.results, attribute)
		}
	}
}

func getItemInfoByItemSearchDoc(itemID int) []string {
	attributes := []string{}
	//effects := []string{}
	fmt.Println("Scraping for Item Attributes....")
	itemSearchString := "https://alla.clumsysworld.com/?a=item&id=" + strconv.Itoa(itemID)

	html := getHTML(itemSearchString)
	defer html.Body.Close()
	itemDoc, err := goquery.NewDocumentFromReader(html.Body)

	if err != nil {
		fmt.Println(err)
	}
	var selection *goquery.Selection
	itemDoc.Find("table.container_div").Find("table").Each(func(index int, subSelection *goquery.Selection) {
		//fmt.Printf("Table # %d: %s\n", index+1, subSelection.Text())
		if index == 1 {
			selection = subSelection
			//fmt.Println("Setting Fultered Seletion: ", selection.Text())
		}
	})

	// Retireve item attributes
	selection.Find("tr").Each(func(subi int, subSelection *goquery.Selection) {
		attributeText := strings.TrimSpace(subSelection.Text())
		if attributeText != "" && attributeText != " " && !strings.Contains(attributeText, "Slot") && (strings.Count(attributeText, ":") == 1 || strings.Count(attributeText, "Level for effect:") == 1) {
			if strings.Contains(attributeText, "Level for effect:") {
				levelIndex := strings.Index(attributeText, "Level for effect:")
				effectString := string(attributeText[0:levelIndex])
				levelString := string(attributeText[levelIndex:])
				//fmt.Println("Adding Attribute: ", effectString)
				attributes = append(attributes, effectString)
				//fmt.Println("Adding Attribute: ", levelString)
				attributes = append(attributes, levelString)
			} else {
				//fmt.Println("Adding Attribute: ", attributeText)
				attributes = append(attributes, attributeText)
			}

		}
	})
	if len(attributes) >= 1 {
		fmt.Println("Scrape Succesful")
	} else {
		fmt.Println("Scrape Failed")
	}
	return attributes
}

func getItemIDandNameByItemSearchDoc(doc *goquery.Document) (int, string) {
	returnItemID := ""
	returnItemName := ""
	fmt.Println("Scraping for Item Name/ID....")
	doc.Find("table.display_table").Find("tr").Each(func(i int, selection *goquery.Selection) {
		//itemName := selection.Find("td.sorting_1").Text()
		itemName := selection.Find("a").Text()
		itemID := selection.Find("a").AttrOr("id", "")
		if i == 21 {
			returnItemID = itemID
			returnItemName = itemName
		}
	})

	convertedItemID, _ := strconv.Atoi(returnItemID)
	if convertedItemID > 0 && returnItemName != "" {
		fmt.Println("Scrape Succesful")
	} else {
		fmt.Println("Scrape Failed")
	}
	return convertedItemID, returnItemName
}

func getHTML(url string) *http.Response {
	resp, err := http.Get(url)
	fmt.Println("Response Code: ", resp.StatusCode)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return resp
}

func formatSearchURL(searchTerm, searchService string) string {
	if searchService == "alla itemID" { // Search by Item ID
		baseURL := "https://alla.clumsysworld.com/?a=item&id="
		searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
		//fmt.Println("Search URL: " + baseURL + searchTerm)
		return baseURL + searchTerm
	} else if searchService == "alla spellID" { // Search by Spell ID
		baseURL := "https://alla.clumsysworld.com/?a=spell&id="
		searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
		//fmt.Println("Search URL: " + baseURL + searchTerm)
		return baseURL + searchTerm
	} else if searchService == "alla npcID" { // Search by NPC ID
		baseURL := "https://alla.clumsysworld.com/?a=npc&id="
		searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
		//fmt.Println("Search URL: " + baseURL + searchTerm)
		return baseURL + searchTerm
	} else if searchService == "alla items" { // Search items by name
		baseURL := "https://alla.clumsysworld.com/?a=items_search&&a=items&iname="
		searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
		alla_url := baseURL + searchTerm + "&iclass=0&irace=0&islot=0&istat1=&istat1comp=%3E%3D&istat1value=&istat2=&istat2comp=%3E%3D&istat2value=&iresists=&iresistscomp=%3E%3D&iresistsvalue=&iheroics=&iheroicscomp=%3E%3D&iheroicsvalue=&imod=&imodcomp=%3E%3D&imodvalue=&itype=-1&iaugslot=0&ieffect=&iminlevel=0&ireqlevel=0&inodrop=0&iavailability=0&iavaillevel=0&ideity=0&isearch=1"
		//fmt.Println("Search URL: " + alla_url)
		return alla_url
	} else if searchService == "alla spells" { // Search spells by name
		baseURL := "https://alla.clumsysworld.com/?a=spells&name="
		searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
		//fmt.Println("Search URL: " + baseURL + searchTerm)
		return baseURL + searchTerm
	}

	return searchTerm
}

type Scraper struct {
	url                  string
	searchService        string
	searchTerm           string
	timerMinutes         int
	timerMinutesDuration int
	timerEnabled         bool
	continuous           bool
	results              []string
}

func (webScraper *Scraper) Scrape() {
	if webScraper.timerEnabled {
		fmt.Printf("Timer Enabled: Initiailizing Scrape Timer at %d minute interval for %d minutes\n", webScraper.timerMinutes, webScraper.timerMinutesDuration)
		ticker := time.NewTicker(time.Duration(webScraper.timerMinutes) * time.Minute)
		minutesElapsed := 0
		if minutesElapsed <= webScraper.timerMinutesDuration {
			minutesElapsed++
			for t := range ticker.C {
				if webScraper.timerEnabled {
					fmt.Printf("Scrape: %d\n", t)
					scrapeURL(webScraper)
				} else {
					fmt.Println("Scrape Timer Disabled")
					ticker.Stop()
				}
			}
		} else {
			ticker.Stop()
			fmt.Printf("Scrape timer stopped after %d minutes\n", webScraper.timerMinutesDuration)
		}

	} else if webScraper.continuous {
		fmt.Println("Continuous Scraping Enabled: Initiailizing Scrape Timer...")
		ticker := time.NewTicker(time.Duration(webScraper.timerMinutes) * time.Minute)
		for t := range ticker.C {
			if webScraper.continuous {
				fmt.Printf("Scrape: %d\n", t)
				scrapeURL(webScraper)
			} else {
				fmt.Println("Continuous Scraping Disabled...")
				ticker.Stop()
			}
		}

	} else {
		scrapeURL(webScraper)
	}
}

func scrapeURL(webScraper *Scraper) {
	scrape(webScraper.url, webScraper.searchService, *webScraper)
}

func (webScraper *Scraper) SetUrl(searchTerm, searchService string) {
	fmt.Println("Setting Url for: " + searchService)
	webScraper.searchService = searchService
	webScraper.searchTerm = searchTerm
	webScraper.url = formatSearchURL(searchTerm, searchService)
}

func (webScraper *Scraper) SetTimer(minutes int) {
	fmt.Printf("Setting Scrape Timer to: %d\n", minutes)
	webScraper.timerMinutes = minutes
}

func (webScraper *Scraper) SetTimerDuration(minutes int) {
	fmt.Printf("Setting Scrape Timer Duration to: %d\n", minutes)
	webScraper.timerMinutesDuration = minutes
}

func (webScraper *Scraper) EnableTimer() {
	fmt.Println("Enabling scrape timer")
	webScraper.timerEnabled = true
}

func (webScraper *Scraper) DisableTimer() {
	fmt.Println("Disabling scrape timer")
	webScraper.timerEnabled = false
}

func (webScraper *Scraper) SetContinuous() {
	fmt.Println("Enabling continuous scraping")
	webScraper.timerEnabled = false
	webScraper.continuous = true
}
