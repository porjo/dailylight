package main

import (
	"bufio"
	"encoding/json"
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
	Books     map[string]string
)

func init() {
	titleReg = regexp.MustCompile(`([A-Z]+ [0-9]+) (MORNING|EVENING)`)
	v1Reg = regexp.MustCompile(`--((?:I{0,4} ?)[A-Z]{3,5})\.? ([0-9]{1,3})\:?([0-9]{0,2})[-,]?([0-9]{0,3}).`)
	versesReg = regexp.MustCompile(`\s*-?((?:I{0,3} ?)[A-Z][a-z]{2,5})\.? ([0-9]{1,3}):([0-9]{1,2})[-,]?([0-9]{0,3}).`)

	Books = make(map[string]string)

	Books["prov"] = "Proverbs"
	Books["i chr"] = "1 Chronicles"
	Books["phi"] = "Philemon"
	Books["ii thes"] = "2 Thessalonians"
	Books["i kgs"] = "1 Kings"
	Books["zech"] = "Zechariah"
	Books["tim"] = "Timothy"
	Books["ezra"] = "Ezra"
	Books["rev"] = "Revelations"
	Books["i pet"] = "1 Peter"
	Books["ii cor"] = "2 Corinthians"
	Books["ruth"] = "Ruth"
	Books["john"] = "John"
	Books["neh"] = "Nehemiah"
	Books["col"] = "Colossians"
	Books["eccl"] = "Ecclesiastes"
	Books["amos"] = "Amos"
	Books["job"] = "Job"
	Books["mark"] = "Mark"
	Books["nah"] = "Nahum"
	Books["kgs"] = "Kings"
	Books["hab"] = "Habbakkuk"
	Books["mal"] = "Malachi"
	Books["acts"] = "Acts"
	Books["luke"] = "Luke"
	Books["jude"] = "Jude"
	Books["tit"] = "Titus"
	Books["i cor"] = "1 Corinthians"
	Books["rom"] = "Romans"
	Books["isa"] = "Isaiah"
	Books["judg"] = "Judges"
	Books["ii kgs"] = "2 Kings"
	Books["james"] = "James"
	Books["mic"] = "Micah"
	Books["exo"] = "Exodus"
	Books["jer"] = "Jeremiah"
	Books["matt"] = "Matthew"
	Books["thes"] = "Thessalonians"
	Books["i john"] = "1 John"
	Books["jon"] = "Jonah"
	Books["jas"] = "James"
	Books["hag"] = "Haggai"
	Books["hos"] = "Hoseah"
	Books["ii sam"] = "2 Samuel"
	Books["pet"] = "Peter"
	Books["esth"] = "Esther"
	Books["deut"] = "Deuteronomy"
	Books["i thes"] = "1 Thessalonians"
	Books["gal"] = "Galatians"
	Books["phil"] = "Philippians"
	Books["joel"] = "Joel"
	Books["psa"] = "Psalms"
	Books["josh"] = "Joshua"
	Books["eph"] = "Ephesians"
	Books["song"] = "Song of Songs"
	Books["lam"] = "Lamentations"
	Books["dan"] = "Daniel"
	Books["zeph"] = "Zephania"
	Books["ii tim"] = "2 Timothy"
	Books["ezek"] = "Ezekiel"
	Books["num"] = "Numbers"
	Books["i sam"] = "1 Samuel"
	Books["ii chr"] = "2 Chronicles"
	Books["ii pet"] = "2 Peter"
	Books["i tim"] = "1 Timothy"
	Books["cor"] = "Corinthians"
	Books["eze"] = "Ezekiel"
	Books["lev"] = "Leviticus"
	Books["heb"] = "Hebrews"
	Books["gen"] = "Genesis"
	Books["iii john"] = "3 John"
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
			if v1r[3] != "" {
				verseStart, err = strconv.Atoi(v1r[3])
				if err != nil {
					panic(err)
				}
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
	entries = append(entries, entry)

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	// Expand abbreviated book names
	for i := 0; i < len(entries); i++ {
		entries[i].Verse1.Book = Books[strings.ToLower(entries[i].Verse1.Book)]
		for j := 0; j < len(entries[i].Verses); j++ {
			entries[i].Verses[j].Book = Books[strings.ToLower(entries[i].Verses[j].Book)]
		}
	}

	json, err := json.Marshal(entries)
	if err != nil {

		panic(err)
	}
	fmt.Printf("%s\n", json)
}
