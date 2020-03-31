package cromp

import (
	"bufio"
	"os"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// AuthorRE is a regex to grab our Authors
var AuthorRE = regexp.MustCompile(`^author:\s(.*)$`)

// TitleRE matches our article title
var TitleRE = regexp.MustCompile(`^title:\s(.*)$`)

// DateRE matches our article date
var DateRE = regexp.MustCompile(`^date:\s(.*)$`)

// UUIDRE matches a uuid in our file
var UUIDRE = regexp.MustCompile("^id: ([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})$")

// Header is the header info loaded from a file
type Header struct {
	Author string
	Title  string
	Date   time.Time
	UUID   uuid.UUID
}

//ParseFileHeader grabs the header info out of an existing file
func ParseFileHeader(f string) (*Header, error) {
	h := &Header{}
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		var line = scanner.Bytes()
		if AuthorRE.Match(line) {
			aline := AuthorRE.ReplaceAllString(string(line), "$1")
			h.Author = aline
		}
		if TitleRE.Match(line) {
			h.Title = TitleRE.ReplaceAllString(string(line), "$1")
		}
		if DateRE.Match(line) {
			d := DateRE.ReplaceAllString(string(line), "$1")
			h.Date, err = time.Parse(time.RFC1123, d)
			if err != nil {
				return nil, err
			}
		}

		if UUIDRE.Match(line) {
			u := UUIDRE.ReplaceAllString(string(line), "$1")
			h.UUID, err = uuid.Parse(u)
			if err != nil {
				return nil, err
			}
		}
	}

	if err != nil {
		return nil, err
	}

	return h, nil
}
