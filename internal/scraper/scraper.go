package scraper

import (
	"encoding/json"
	"fmt"
	"log"

	"tomorrowland-scraper/internal/models"

	"github.com/gocolly/colly"
)

type nextData struct {
	Props struct {
		PageProps struct {
			Doc struct {
				Blocks []json.RawMessage `json:"blocks"`
			} `json:"doc"`
		} `json:"pageProps"`
	} `json:"props"`
}

func Scrape() (lineup models.Lineup, err error) {
	// Step 1: fetch the lineup page and extract event + uuid from __NEXT_DATA__
	c := colly.NewCollector()

	var cdnURL string

	c.OnHTML("script#__NEXT_DATA__", func(e *colly.HTMLElement) {
		var nd nextData
		if err := json.Unmarshal([]byte(e.Text), &nd); err != nil {
			log.Fatal(err)
			return
		}

		for _, raw := range nd.Props.PageProps.Doc.Blocks {
			var block models.Festival
			if err := json.Unmarshal(raw, &block); err != nil {
				continue
			}
			if block.Type == "line-up" {
				cdnURL = fmt.Sprintf("https://artist-lineup-cdn.tomorrowland.com/%s-W1-%s.json", block.Event, block.UUID)
				fmt.Println("CDN URL:", cdnURL)
			}
		}
	})

	err = c.Visit("https://belgium.tomorrowland.com/en/line-up/")
	c.Wait()

	if err != nil {
		return models.Lineup{}, err
	}

	if cdnURL == "" {
		return models.Lineup{}, err
	}

	c2 := colly.NewCollector()
	c2.OnResponse(func(r *colly.Response) {
		if err := json.Unmarshal(r.Body, &lineup); err != nil {
			log.Fatal(err)
			return
		}
	})

	err = c2.Visit(cdnURL)
	c2.Wait()

	if err != nil {
		return models.Lineup{}, err
	}

	return lineup, nil
}
