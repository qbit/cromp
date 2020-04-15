package cromp

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AuthorRE is a regex to grab our Authors
var AuthorRE = regexp.MustCompile(`^author:\s(.*)$`)

// TitleRE matches our article title for either plain text or org-mode
var TitleRE = regexp.MustCompile(`^(?:\*|title:)\s(.*)$`)

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

// ReadFileBody grabs the entire file
func ReadFileBody(f string) ([]byte, error) {
	return ioutil.ReadFile(f)
}

//ParseHeader grabs the header info out of a string
func ParseHeader(f string) (*Header, error) {
	var err error
	h := &Header{}
	for _, line := range strings.Split(f, "\n") {
		if AuthorRE.MatchString(line) {
			aline := AuthorRE.ReplaceAllString(line, "$1")
			if h.Author == "" {
				h.Author = aline
			}
		}

		if TitleRE.MatchString(line) {
			if h.Title == "" {
				h.Title = TitleRE.ReplaceAllString(line, "$1")
			}
		}

		if DateRE.MatchString(line) {
			if h.Date.String() == "" {
				d := DateRE.ReplaceAllString(line, "$1")
				h.Date, err = time.Parse(time.RFC1123, d)
				if err != nil {
					return nil, err
				}
			}
		}

		if UUIDRE.MatchString(line) {
			u := UUIDRE.ReplaceAllString(line, "$1")
			h.UUID, err = uuid.Parse(u)
			if err != nil {
				return nil, err
			}
		}
	}

	return h, nil
}

//ParseFileHeader grabs the header info out of an existing file
func ParseFileHeader(f string) (*Header, error) {
	h := &Header{}
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()
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
