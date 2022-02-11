package webScrape

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func scrape(searchURL, serviceName string) {
	html := getHTML(searchURL)
	defer html.Body.Close()
	doc, error := goquery.NewDocumentFromReader(html.Body)

	if error != nil {
		fmt.Println(error)
	}

	scrapePageData(doc, serviceName)

}

func scrapePageData(doc *goquery.Document, serviceName string) string {
	if serviceName == "ebay" {
		doc.Find("ul.srp-results>li.s-item").Each(func(i int, s *goquery.Selection) {
			title := s.Find("a.s-item__link").Text()
			price := s.Find("span.s-item__price").Text()
			fmt.Println(title)
			fmt.Println(price)
		})
	} else if serviceName == "alla items" {
		returnItemID := ""
		fmt.Println("HTML Title: ", doc.Find("title").Text())
		fmt.Println("TEST: ", doc.Find("table.display_table.tbody").Text())
		doc.Find("table.display_table").Find("tr").Each(func(i int, selection *goquery.Selection) {
			//itemName := selection.Find("td.sorting_1").Text()
			itemName := selection.Find("a").Text()
			itemID := selection.Find("a").AttrOr("id", "")
			if i == 21 {
				returnItemID = selection.Find("a").AttrOr("id", "")
			}
			if itemName != "" && i >= 21 && i <= 26 {
				fmt.Printf("Item Name: %s\n", itemName)
				fmt.Printf("Item ID: %s\n", itemID)
			}
		})
		return returnItemID
	}
	return ""
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
	if searchService == "ebay" {
		baseURL := "https://www.ebay.com/sch/i.html?_from=R40&_trksid=p2334524.m570.l1313&_nkw="
		searchTerm = strings.Replace(searchTerm, " ", "+", -1)
		fmt.Println("Search URL: " + baseURL + searchTerm)
		return baseURL + searchTerm
	} else if searchService == "alla itemID" { // Search by Item ID
		baseURL := "https://alla.clumsysworld.com/?a=item&id="
		searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
		fmt.Println("Search URL: " + baseURL + searchTerm)
		return baseURL + searchTerm
	} else if searchService == "alla spellID" { // Search by Spell ID
		baseURL := "https://alla.clumsysworld.com/?a=spell&id="
		searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
		fmt.Println("Search URL: " + baseURL + searchTerm)
		return baseURL + searchTerm
	} else if searchService == "alla npcID" { // Search by NPC ID
		baseURL := "https://alla.clumsysworld.com/?a=npc&id="
		searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
		fmt.Println("Search URL: " + baseURL + searchTerm)
		return baseURL + searchTerm
	} else if searchService == "alla items" { // Search items by name
		baseURL := "https://alla.clumsysworld.com/?a=items_search&&a=items&iname="
		searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
		alla_url := baseURL + searchTerm + "&iclass=0&irace=0&islot=0&istat1=&istat1comp=%3E%3D&istat1value=&istat2=&istat2comp=%3E%3D&istat2value=&iresists=&iresistscomp=%3E%3D&iresistsvalue=&iheroics=&iheroicscomp=%3E%3D&iheroicsvalue=&imod=&imodcomp=%3E%3D&imodvalue=&itype=-1&iaugslot=0&ieffect=&iminlevel=0&ireqlevel=0&inodrop=0&iavailability=0&iavaillevel=0&ideity=0&isearch=1"
		fmt.Println("Search URL: " + alla_url)
		return alla_url
	} else if searchService == "alla spells" { // Search spells by name
		baseURL := "https://alla.clumsysworld.com/?a=spells&name="
		searchTerm = strings.Replace(searchTerm, " ", "%20", -1)
		fmt.Println("Search URL: " + baseURL + searchTerm)
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
}

func (webScraper *Scraper) Scrape() {
	if webScraper.timerEnabled {
		fmt.Printf("Timer Enabled: Initiailizing Scrape Timer at %d minute interval for %d minutes\n", webScraper.timerMinutes, webScraper.timerMinutesDuration)
		ticker := time.NewTicker(time.Duration(webScraper.timerMinutes) * time.Minute)
		minutesElapsed := 0
		// for every `tick` that our `ticker`
		// emits, we print `tock`
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
		// for every `tick` that our `ticker`
		// emits, we print `tock`
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
	fmt.Println("Scraping: " + webScraper.url)
	scrape(webScraper.url, webScraper.searchService)
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
