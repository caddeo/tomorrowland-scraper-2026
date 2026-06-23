package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"tomorrowland-scraper/internal/diff"
	"tomorrowland-scraper/internal/models"
	"tomorrowland-scraper/internal/scraper"
	"tomorrowland-scraper/internal/storage"
)

func exportToCSV(lineup models.Lineup, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	if err := w.Write([]string{"Day", "Date", "Stage", "Start", "End", "Artist"}); err != nil {
		return err
	}

	for _, p := range lineup.Performances {

		start := p.StartTime
		end := p.EndTime
		if len(start) >= 16 {
			start = start[11:16]
		}
		if len(end) >= 16 {
			end = end[11:16]
		}

		if err := w.Write([]string{p.Day, p.Date, p.Stage.Name, start, end, p.Name}); err != nil {
			return err
		}
	}

	return w.Error()
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
	lineup, err := scraper.Scrape()
	if err != nil {
		return err
	}

	fmt.Printf("%-10s %-30s %-8s %-8s %s\n", "Day", "Stage", "Start", "End", "Artist")
	fmt.Println(strings.Repeat("-", 80))
	for _, p := range lineup.Performances {
		start := p.StartTime[11:16]
		end := p.EndTime[11:16]
		fmt.Printf("%-10s %-30s %-8s %-8s %s\n", p.Day, p.Stage.Name, start, end, p.Name)
	}

	if err := storage.Archive(lineup); err != nil {
		return err
	}

	files, err := storage.List()
	if err != nil || len(files) < 2 {
		return err
	}

	if err := runDiff(files[1], files[0]); err != nil {
		return err
	}

	return nil
}

func runDiff(oldFile, newFile string) error {
	old, err := storage.Load(oldFile)
	if err != nil {
		return fmt.Errorf("load %s: %w", oldFile, err)
	}
	new, err := storage.Load(newFile)
	if err != nil {
		return fmt.Errorf("load %s: %w", newFile, err)
	}

	d := diff.Diff(old, new)
	fmt.Print(d.Format(oldFile, newFile))
	return nil
}
