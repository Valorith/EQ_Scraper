package main

import "github.com/Valorith/EQ_Scraper/webScrape"

func main() {

	// Run a new web scraper
	//-----------------------------------------
	webScraper := webScrape.Scraper{}
	webScraper.SetUrl("Kromzek Kings", "alla items") // ebay, alla itemID, alla items, alla spellID, alla spells, all npcID
	webScraper.SetTimer(1)                           // Minutes
	webScraper.SetTimerDuration(10)
	//webScraper.EnableTimer()
	//webScraper.SetContinuous()
	webScraper.Scrape()
	//-----------------------------------------

}
