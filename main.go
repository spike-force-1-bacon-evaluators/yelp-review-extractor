package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// Review ...
type Review struct {
	Votes struct {
		Funny  int `json:"funny"`
		Useful int `json:"useful"`
		Cool   int `json:"cool"`
	} `json:"votes"`
	UserID     string `json:"user_id"`
	ReviewID   string `json:"review_id"`
	Stars      int    `json:"stars"`
	Date       string `json:"date"`
	Text       string `json:"text"`
	Type       string `json:"type"`
	BusinessID string `json:"business_id"`
}

func main() {

	var review []*Review

	// Read dataset
	file, err := os.Open("data/yelp_academic_dataset_review.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		firstIndex := strings.Index(text, "{\"votes\":")
		if firstIndex != -1 {
			r, err := unmarshal(text[firstIndex:])
			if err != nil {
				log.Fatal(err)
			}
			review = append(review, r)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if err := writeOut(review); err != nil {
		log.Fatal(err)
	}
}

// Parse JSON to Review struct
func unmarshal(item string) (*Review, error) {
	var r Review
	if err := json.Unmarshal([]byte(item), &r); err != nil {
		return &Review{}, err
	}
	return &r, nil
}

// Write stars and review text on a CSV file
func writeOut(review []*Review) error {
	file, err := os.Create("data/yelp_academic_dataset_review.tsv")
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprint(w, "star\treview\n")

	for _, r := range review {
		text := removeNewlines(r.Text)
		text = removeTabs(text)
		star := normaliseStars(r.Stars)
		fmt.Fprintf(w, "%d\t%s\n", star, strings.ToLower(text))
	}
	return nil
}

// Remove newlines from review text
func removeNewlines(text string) string {
	return strings.Replace(text, "\n", " ", -1)
}

// Remove tabs from review text
func removeTabs(text string) string {
	return strings.Replace(text, "\t", " ", -1)
}

// Normalise stars for binary values
func normaliseStars(star int) int {
	if star > 3 {
		return 1
	}
	return 0
}
