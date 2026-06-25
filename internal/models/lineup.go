package models

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
)

type Lineup struct {
	Hash         string        `json:"hash,omitempty"`
	Performances []Performance `json:"performances"`
}

func (l Lineup) ComputeHash() string {
	sorted := make([]Performance, len(l.Performances))
	copy(sorted, l.Performances)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ID < sorted[j].ID
	})

	data, _ := json.Marshal(sorted)
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum[:8])
}
