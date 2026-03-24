package models

import "time"

type SiteAnalytics struct {
	PostID string `json:"postId"`
	Views int `json:"views"`
	UniqueViews int `json:"uniqueViews"`
	AvgReadingTime int `json:"avgReadingTime"`
	AvgScrollDepth int `json:"avgScrollDepth"`
	Date time.Time `json:"date"`
}

func CalcReadingTime(body string) int {
	words := 0
	stillWord := false

	for _, wr := range body {
		if wr == ' ' || wr == '\t' || wr == '\n' {
			if !stillWord {
				words++
			}
			stillWord = false
		} else {
			stillWord = true
		}
	}

	if mins := words / 200; mins < 1 {
		return 1
	} else {
		return mins
	}

}

