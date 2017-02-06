package parser

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

const testFile = "../test/input.json"

func TestUnmarshal(t *testing.T) {
	expectedStars := 4
	expectedText := `Mr Hoagie is an institution. Walking in, it does seem like a throwback to 30 years ago, old fashioned menu board, booths out of the 70s, and a large selection of food. Their speciality is the Italian Hoagie, and it is voted the best in the area year after year. I usually order the burger, while the patties are obviously cooked from frozen, all of the other ingredients are very fresh. Overall, its a good alternative to Subway, which is down the road.`

	buf, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Error(err)
	}

	file := string(buf)

	r, err := unmarshal(file)
	if err != nil {
		t.Error(err)
	}

	if r.Stars != 4 {
		t.Errorf("Expecting %d stars. Got %d", expectedStars, r.Stars)
	}

	if r.Text != expectedText {
		t.Errorf("Expecting text:\n%s\nGot:\n%s\n", expectedText, r.Text)
	}
}

func TestWriteOut(t *testing.T) {
	a := &Review{
		Stars: 1,
		Text:  "foo",
	}
	b := &Review{
		Stars: 2,
		Text:  "bar",
	}
	c := &Review{
		Stars: 3,
		Text:  "foobar",
	}

	review := []*Review{a, b, c}

	testOutput := "../test/writeout.txt"
	if err := writeOut(review, testOutput); err != nil {
		t.Error(err)
	}

	file, err := os.Open(testOutput)
	if err != nil {
		t.Error(err)
	}

	expected := []string{"star	review", "0	foo", "0	bar", "1	foobar"}

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		if scanner.Text() != expected[i] {
			t.Errorf("Expecting: %s. Got: %s", expected[i], scanner.Text())
		}
	}
}

func TestRemoveNewLines(t *testing.T) {
	input := "foo\nbar\n"
	expected := "foo bar "
	if observed := removeNewlines(input); observed != expected {
		t.Errorf("Expected: %s. Got: %s", expected, observed)
	}
}

func TestRemoveTabs(t *testing.T) {
	input := "foo\tbar\t"
	expected := "foo bar "
	if observed := removeTabs(input); observed != expected {
		t.Errorf("Expected: %s. Got: %s", expected, observed)
	}
}

var getIndexPositionCases = []struct {
	inputStr, inputSubstr string
	expected              int
}{
	{"foobar", "bar", 3},
	{"foobar", "fizz", -1},
}

func TestGetIndexPosition(t *testing.T) {
	for _, test := range getIndexPositionCases {
		if observed := getIndexPosition(test.inputStr, test.inputSubstr); observed != test.expected {
			t.Errorf("For string %s and substring %s. Expected: %d. Got: %d",
				test.inputStr, test.inputSubstr, test.expected, observed)
		}
	}
}

var normaliseStarsCases = []struct {
	input, expected int
}{
	{4, 1},
	{2, 0},
}

func TestNormaliseStars(t *testing.T) {
	for _, test := range normaliseStarsCases {
		if observed := normaliseStars(test.input); observed != test.expected {
			t.Errorf("For star: %d. Expected: %d. Got: %d",
				test.input, test.expected, observed)
		}
	}
}

func TestReadFile(t *testing.T) {
	expected := `{"votes": {"funny": 0, "useful": 0, "cool": 0}, "user_id": "PUFPaY9KxDAcGqfsorJp3Q", "review_id": "Ya85v4eqdd6k9Od8HbQjyA", "stars": 4, "date": "2012-08-01", "text": "Mr Hoagie is an institution. Walking in, it does seem like a throwback to 30 years ago, old fashioned menu board, booths out of the 70s, and a large selection of food. Their speciality is the Italian Hoagie, and it is voted the best in the area year after year. I usually order the burger, while the patties are obviously cooked from frozen, all of the other ingredients are very fresh. Overall, its a good alternative to Subway, which is down the road.", "type": "review", "business_id": "5UmKMjUEUNdYWqANhGckJw"}`

	lineCh := make(chan string)
	go readFile(testFile, lineCh)

	for line := range lineCh {
		if line != expected {
			t.Errorf("Expecting line:\n%s\nGot:\n%s\n", expected, line)
		}
	}
}
