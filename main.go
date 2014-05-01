package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	Date   time.Time
	AMPM   string
	Verse1 Ref
	Verses []Ref
}

type Ref struct {
	Book       string
	Chapter    int
	VerseStart int
	VerseEnd   int
}

var (
	titleReg  *regexp.Regexp
	v1Reg     *regexp.Regexp
	versesReg *regexp.Regexp
)

func init() {
	titleReg = regexp.MustCompile(`([A-Z]+ [0-9]+) (MORNING|EVENING)`)
	v1Reg = regexp.MustCompile(`--((?:I{0,2} ?)[A-Z]{3,5})\.? ([0-9]{1,3}):([0-9]{1,2})[-,]?([0-9]{0,3}).`)
	versesReg = regexp.MustCompile(`\s*-?((?:I{0,2} ?)[a-zA-Z]{3,5})\.? ([0-9]{1,3}):([0-9]{1,2})[-,]?([0-9]{0,3}).`)
}

func main() {

	file, err := os.Open("daily_light.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)

	var entries []*Entry
	var entry *Entry

	for scanner.Scan() {

		txt := scanner.Text()
		// The following never matches...why?
		if txt != "" && strings.Contains("Isa. 33:17", txt) {
			fmt.Println("Found text")
		}

		titler := titleReg.FindStringSubmatch(scanner.Text())
		v1r := v1Reg.FindStringSubmatch(scanner.Text())
		versesr := versesReg.FindAllStringSubmatch(scanner.Text(), -1)

		if len(titler) == 3 {
			var err error
			if entry != nil {
				entries = append(entries, entry)
			}
			entry = &Entry{}
			entry.Date, err = time.Parse("January 2", titler[1])
			if err != nil {
				panic(err)
			}
			if titler[2] == "MORNING" {
				entry.AMPM = "AM"
			} else {
				entry.AMPM = "PM"
			}
		}
		if len(v1r) == 5 {
			var err error
			var chapter, verseStart, verseEnd int
			chapter, err = strconv.Atoi(v1r[2])
			if err != nil {
				panic(err)
			}
			verseStart, err = strconv.Atoi(v1r[3])
			if err != nil {
				panic(err)
			}
			if v1r[4] != "" {
				verseEnd, err = strconv.Atoi(v1r[4])
				if err != nil {
					panic(err)
				}
			} else {
				verseEnd = verseStart
			}
			entry.Verse1 = Ref{v1r[1], chapter, verseStart, verseEnd}
		}

		if len(versesr) > 0 {
			for i := 0; i < len(versesr); i++ {
				var err error
				var chapter, verseStart, verseEnd int
				chapter, err = strconv.Atoi(versesr[i][2])
				if err != nil {
					panic(err)
				}
				verseStart, err = strconv.Atoi(versesr[i][3])
				if err != nil {
					panic(err)
				}

				if versesr[i][4] != "" {
					verseEnd, err = strconv.Atoi(versesr[i][4])
					if err != nil {
						panic(err)
					}
				} else {
					verseEnd = verseStart
				}

				ref := Ref{versesr[i][1], chapter, verseStart, verseEnd}

				if ref != entry.Verse1 {
					entry.Verses = append(entry.Verses, ref)
				}
			}
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	for i := 0; i < len(entries); i++ {
		fmt.Printf("entry %d: %v\n", i, entries[i])
	}
}
