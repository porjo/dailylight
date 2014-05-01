package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Regex struct {
	Title  *regexp.Regexp
	V1     *regexp.Regexp
	Verses *regexp.Regexp
}

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

var R Regex

func init() {
	R.Title = regexp.MustCompile(`([A-Z]+ [0-9]+) (MORNING|EVENING)`)
	//R.V1 = regexp.MustCompile(`--([A-Z]{3,4})\. ([0-9]{1,3}):([0-9]{1,2})`)
	R.V1 = regexp.MustCompile(`--((?:I{0,2} ?)[A-Z]{3,5})\.? ([0-9]{1,3}):([0-9]{1,2})[-,]?([0-9]{0,3}).`)
	R.Verses = regexp.MustCompile(`\s*-?((?:I{0,2} ?)[a-zA-Z]{3,5})\.? ([0-9]{1,3}):([0-9]{1,2})[-,]?([0-9]{0,3}).`)
}

func main() {

	file, _ := os.Open("daily_light.txt")
	scanner := bufio.NewScanner(file)

	var entries []*Entry
	var entry *Entry

	for scanner.Scan() {

		tmatches := R.Title.FindStringSubmatch(scanner.Text())
		vmatches := R.V1.FindStringSubmatch(scanner.Text())
		vvmatches := R.Verses.FindAllStringSubmatch(scanner.Text(), -1)

		if len(tmatches) == 3 {
			var err error
			if entry != nil {
				entries = append(entries, entry)
				//fmt.Printf("Entry %#v\n", entry)
			}
			entry = &Entry{}
			entry.Date, err = time.Parse("January 2", tmatches[1])
			if err != nil {
				panic(err)
			}
			if tmatches[2] == "MORNING" {
				entry.AMPM = "AM"
			} else {
				entry.AMPM = "PM"
			}
			//fmt.Printf("TMatch %s, %s\n", tmatches[1], tmatches[2])
		}
		if len(vmatches) == 5 {
			var err error
			var chapter, verseStart, verseEnd int
			chapter, err = strconv.Atoi(vmatches[2])
			if err != nil {
				panic(err)
			}
			verseStart, err = strconv.Atoi(vmatches[3])
			if err != nil {
				panic(err)
			}
			if vmatches[4] != "" {
				verseEnd, err = strconv.Atoi(vmatches[4])
				if err != nil {
					panic(err)
				}
			} else {
				verseEnd = verseStart
			}
			entry.Verse1 = Ref{vmatches[1], chapter, verseStart, verseEnd}
			//fmt.Printf("VMatch %s %s:%s\n", vmatches[1], vmatches[2], vmatches[3])
		}

		if len(vvmatches) > 0 {
			//fmt.Printf("VVMatch: ")
			for i := 0; i < len(vvmatches); i++ {
				var err error
				var chapter, verseStart, verseEnd int
				chapter, err = strconv.Atoi(vvmatches[i][2])
				if err != nil {
					panic(err)
				}
				verseStart, err = strconv.Atoi(vvmatches[i][3])
				if err != nil {
					panic(err)
				}

				if vvmatches[i][4] != "" {
					verseEnd, err = strconv.Atoi(vvmatches[i][4])
					if err != nil {
						panic(err)
					}
				} else {
					verseEnd = verseStart
				}

				ref := Ref{vvmatches[i][1], chapter, verseStart, verseEnd}

				if ref != entry.Verse1 {
					entry.Verses = append(entry.Verses, ref)
				}
			}
			//fmt.Println()
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	for i := 0; i < len(entries); i++ {
		fmt.Printf("entry %d: %v\n", i, entries[i])
	}
}
