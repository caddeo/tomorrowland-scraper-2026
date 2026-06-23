package diff_test

import (
	"fmt"
	"testing"

	"tomorrowland-scraper/internal/diff"
	"tomorrowland-scraper/internal/storage"
)

func TestDiffDemo(t *testing.T) {
	old, err := storage.LoadPath("../../archive/tomorrowland-w1-2026-06-22-23010115.json")
	if err != nil {
		t.Fatal(err)
	}
	new, err := storage.LoadPath("../../testdata/tomorrowland-w1-demo-modified.json")
	if err != nil {
		t.Fatal(err)
	}

	d := diff.Diff(old, new)
	fmt.Print(d.Format(
		"tomorrowland-w1-2026-06-22-23010115.json",
		"tomorrowland-w1-demo-modified.json",
	))
}
