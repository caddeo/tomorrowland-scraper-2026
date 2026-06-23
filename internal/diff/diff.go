package diff

import (
	"fmt"
	"strings"

	"tomorrowland-scraper/internal/models"
)

type FieldChange struct {
	Field string
	Old   string
	New   string
}

type PerformanceDiff struct {
	Old            models.Performance
	New            models.Performance
	Changes        []FieldChange
	ArtistsAdded   []models.Artist
	ArtistsRemoved []models.Artist
}

type LineupDiff struct {
	Added   []models.Performance
	Removed []models.Performance
	Changed []PerformanceDiff
}

func Diff(old, new models.Lineup) LineupDiff {
	oldMap := make(map[string]models.Performance, len(old.Performances))
	for _, p := range old.Performances {
		oldMap[p.ID] = p
	}

	newMap := make(map[string]models.Performance, len(new.Performances))
	for _, p := range new.Performances {
		newMap[p.ID] = p
	}

	var result LineupDiff

	for id, oldP := range oldMap {
		newP, exists := newMap[id]
		if !exists {
			result.Removed = append(result.Removed, oldP)
			continue
		}
		if pd, changed := comparePerformances(oldP, newP); changed {
			result.Changed = append(result.Changed, pd)
		}
	}

	for id, newP := range newMap {
		if _, exists := oldMap[id]; !exists {
			result.Added = append(result.Added, newP)
		}
	}

	return result
}

func comparePerformances(old, new models.Performance) (PerformanceDiff, bool) {
	pd := PerformanceDiff{Old: old, New: new}

	check := func(field, oldVal, newVal string) {
		if oldVal != newVal {
			pd.Changes = append(pd.Changes, FieldChange{Field: field, Old: oldVal, New: newVal})
		}
	}

	check("stage", old.Stage.Name, new.Stage.Name)
	check("day", old.Day, new.Day)
	check("date", old.Date, new.Date)
	check("start_time", old.StartTime, new.StartTime)
	check("end_time", old.EndTime, new.EndTime)

	oldArtists := make(map[string]models.Artist, len(old.Artists))
	for _, a := range old.Artists {
		oldArtists[a.ID] = a
	}
	newArtists := make(map[string]models.Artist, len(new.Artists))
	for _, a := range new.Artists {
		newArtists[a.ID] = a
	}

	for id, a := range oldArtists {
		if _, exists := newArtists[id]; !exists {
			pd.ArtistsRemoved = append(pd.ArtistsRemoved, a)
		}
	}
	for id, a := range newArtists {
		if _, exists := oldArtists[id]; !exists {
			pd.ArtistsAdded = append(pd.ArtistsAdded, a)
		}
	}

	changed := len(pd.Changes) > 0 || len(pd.ArtistsAdded) > 0 || len(pd.ArtistsRemoved) > 0
	return pd, changed
}

func fmtTime(t string) string {
	if len(t) >= 16 {
		return t[11:16]
	}
	return t
}

func fmtPerformance(p models.Performance) string {
	return fmt.Sprintf("%s — %s, %s %s–%s", p.Name, p.Stage.Name, p.Day, fmtTime(p.StartTime), fmtTime(p.EndTime))
}

// Format returns a diff-style string. Pass the source filenames for the header.
func (d LineupDiff) Format(oldFile, newFile string) string {
	if len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0 {
		return "No differences."
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "--- %s\n", oldFile)
	fmt.Fprintf(&sb, "+++ %s\n", newFile)

	for _, p := range d.Removed {
		fmt.Fprintf(&sb, "\n- %s\n", fmtPerformance(p))
	}

	for _, p := range d.Added {
		fmt.Fprintf(&sb, "\n+ %s\n", fmtPerformance(p))
	}

	for _, pd := range d.Changed {
		fmt.Fprintf(&sb, "\n@@ %s @@\n", pd.Old.Name)
		fmt.Fprintf(&sb, "- %s\n", fmtPerformance(pd.Old))
		fmt.Fprintf(&sb, "+ %s\n", fmtPerformance(pd.New))
		for _, a := range pd.ArtistsRemoved {
			fmt.Fprintf(&sb, "-   artist: %s\n", a.Name)
		}
		for _, a := range pd.ArtistsAdded {
			fmt.Fprintf(&sb, "+   artist: %s\n", a.Name)
		}
	}

	return sb.String()
}
