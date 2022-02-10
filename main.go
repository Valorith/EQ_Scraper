package main

import "github.com/Valorith/EQ_Scraper/webScrape"

func main() {

	webScraper := webScrape.Scraper{}
	webScraper.SetUrl("Orb", "alla items") // ebay, alla itemID, alla items, alla spellID, alla spells, all npcID
	webScraper.SetTimer(1)                 // Minutes
	webScraper.SetTimerDuration(10)
	//webScraper.EnableTimer()
	//webScraper.SetContinuous()
	webScraper.Scrape()

}
