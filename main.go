package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

type NextData struct {
	Props struct {
		PageProps struct {
			Doc struct {
				Blocks []json.RawMessage `json:"blocks"`
			} `json:"doc"`
		} `json:"pageProps"`
	} `json:"props"`
}

type LineUpBlock struct {
	Type  string `json:"type"`
	Event string `json:"event"`
	UUID  string `json:"uuid"`
}

type Performance struct {
	Name      string `json:"name"`
	Stage     Stage  `json:"stage"`
	Day       string `json:"day"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type Stage struct {
	Name string `json:"name"`
}

type Lineup struct {
	Performances []Performance `json:"performances"`
}

func main() {
	log.Println("Starting program..")
	if err := run(); err != nil {
		log.Fatalln(err)
		os.Exit(0)
	}

	log.Print("Stopping..")
	os.Exit(1)
}

func run() (err error) {
	// Step 1: fetch the lineup page and extract event + uuid from __NEXT_DATA__
	c := colly.NewCollector()

	var cdnURL string

	c.OnHTML("script#__NEXT_DATA__", func(e *colly.HTMLElement) {
		var nd NextData
		if err := json.Unmarshal([]byte(e.Text), &nd); err != nil {
			log.Fatal(err)
			return
		}

		for _, raw := range nd.Props.PageProps.Doc.Blocks {
			var block LineUpBlock
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
		return err
	}

	if cdnURL == "" {
		return err
	}

	// Step 2: fetch the CDN JSON and print performances
	c2 := colly.NewCollector()

	c2.OnResponse(func(r *colly.Response) {
		var lineup Lineup
		if err := json.Unmarshal(r.Body, &lineup); err != nil {
			log.Fatal(err)
			return
		}

		fmt.Printf("%-10s %-30s %-8s %-8s %s\n", "Day", "Stage", "Start", "End", "Artist")
		fmt.Println(strings.Repeat("-", 80))
		for _, p := range lineup.Performances {
			start := p.StartTime[11:16]
			end := p.EndTime[11:16]
			fmt.Printf("%-10s %-30s %-8s %-8s %s\n", p.Day, p.Stage.Name, start, end, p.Name)
		}
	})

	err = c2.Visit(cdnURL)
	c2.Wait()

	if err != nil {
		return err
	}

	return nil
}
