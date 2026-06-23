package diff_test

import (
	"strings"
	"testing"

	"tomorrowland-scraper/internal/diff"
	"tomorrowland-scraper/internal/storage"
)

func TestDiff(t *testing.T) {
	old, err := storage.LoadPath("../../testdata/old.json")
	if err != nil {
		t.Fatal(err)
	}
	new, err := storage.LoadPath("../../testdata/new.json")
	if err != nil {
		t.Fatal(err)
	}

	d := diff.Diff(old, new)

	// 0003 and 0005 removed, 0006 added, 0001 stage changed, 0002 time changed, 0004 artist added
	if len(d.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(d.Removed))
	}
	if len(d.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(d.Added))
	}
	if len(d.Changed) != 3 {
		t.Errorf("expected 3 changed, got %d", len(d.Changed))
	}

	out := d.Format("old.json", "new.json")

	for _, want := range []string{
		// removed
		"- Artist C — MAINSTAGE, SUNDAY 14:00–15:30",
		"- Artist F — MAINSTAGE, FRIDAY 16:00–17:30",
		// added
		"+ Artist H — MAINSTAGE, SUNDAY 14:00–15:30",
		// stage move
		"@@ Artist A @@",
		"- Artist A — MAINSTAGE, SATURDAY 20:00–21:30",
		"+ Artist A — CRYSTAL GARDEN, SATURDAY 20:00–21:30",
		// time change
		"@@ Artist B @@",
		"- Artist B — CRYSTAL GARDEN, SATURDAY 18:00–19:30",
		"+ Artist B — CRYSTAL GARDEN, SATURDAY 19:00–20:30",
		// b2b artist added
		"@@ Artist D b2b Artist E @@",
		"+   artist: Artist G",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n\nfull output:\n%s", want, out)
		}
	}
}
