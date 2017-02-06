package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

// Review struct for JSON values
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

const (
	inputFile          = "data/yelp_academic_dataset_review.json"
	outputFile         = "data/yelp_academic_dataset_review.tsv"
	delimiter          = "{\"votes\":"
	errorMsgUnmarshal  = "failed to unmarshal json"
	errorMsgWriteOut   = "failed to write output"
	errorMsgReadFile   = "failed to open input file"
	loadingMsgProcess  = "Processing data "
	loadingMsgWrite    = "Writing output "
	loadingMsgComplete = "Ouput file: yelp_academic_dataset_review.tsv\n"
)

var loading = spinner.New(spinner.CharSets[24], 100*time.Millisecond)

// Run triggers data processing
func Run() {
	loading.Prefix = loadingMsgProcess
	loading.Start()

	var review []*Review

	lineCh := make(chan string, 100)

	go readFile(inputFile, lineCh)

	for line := range lineCh {
		index := getIndexPosition(line, delimiter)
		if index != -1 {
			r, err := unmarshal(line[index:])
			if err != nil {
				log.Fatalf("%s => %s", errorMsgUnmarshal, err)
			}
			review = append(review, r)
		}
	}
	loading.Stop()
	loading.Prefix = loadingMsgWrite
	loading.FinalMSG = loadingMsgComplete
	loading.Start()

	if err := writeOut(review, outputFile); err != nil {
		log.Fatalf("%s => %s", errorMsgWriteOut, err)
	}
	loading.Stop()
}

// Read dataset and send lines through string channel
func readFile(path string, lineCh chan string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("%s => %s", errorMsgWriteOut, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCh <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	close(lineCh)
}

// Parse JSON to Review struct
func unmarshal(jsonStr string) (*Review, error) {
	var r Review
	if err := json.Unmarshal([]byte(jsonStr), &r); err != nil {
		return &Review{}, err
	}
	return &r, nil
}

// Write stars and review text on a CSV file
func writeOut(review []*Review, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	_, err = w.WriteString("star\treview\n")
	if err != nil {
		return err
	}

	w.Flush()

	for _, r := range review {
		text := removeNewlines(r.Text)
		text = removeTabs(text)
		star := normaliseStars(r.Stars)
		row := fmt.Sprintf("%d\t%s\n", star, strings.ToLower(text))
		_, err := w.WriteString(row)
		if err != nil {
			return err
		}
		w.Flush()
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
	if star >= 3 {
		return 1
	}
	return 0
}

// Get index position of given substring
func getIndexPosition(str, substr string) int {
	return strings.Index(str, substr)
}
