package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"tomorrowland-scraper/internal/models"
)

func Archive(lineup models.Lineup) (err error) {
	t := time.Now()
	filename := fmt.Sprintf("tomorrowland-w1-%s-%s%02d.json",
		t.Format("2006-01-02"),
		t.Format("150405"),
		t.Nanosecond()/10000000,
	)
	file, err := os.Create(fmt.Sprintf("archive/%s", filename))
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	lineup.Hash = lineup.ComputeHash()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(lineup); err != nil {
		return err
	}

	log.Printf("Saved to %s", filename)
	return nil
}

func List() ([]string, error) {
	entries, err := os.ReadDir("archive")
	if err != nil {
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	return names, nil
}

func Load(filename string) (lineup models.Lineup, err error) {
	return LoadPath(fmt.Sprintf("archive/%s", filename))
}

func LoadPath(path string) (lineup models.Lineup, err error) {
	file, err := os.Open(path)
	if err != nil {
		return models.Lineup{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&lineup); err != nil {
		return models.Lineup{}, err
	}

	return lineup, nil
}
