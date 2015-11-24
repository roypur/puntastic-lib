package puntastic

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"golang.org/x/mobile/asset"
	"math/rand"
	"strconv"
	"strings"
)

type dictionary struct {
	//map-of-words-with-n-syllables, slice-of-words, slice-of-syllables
	Words map[int][][]string
}
type Pun struct {
	Pun       string
	BaseWord string
}

var pos map[int]int
var end map[int]bool

var lastPos int

var lastInput string
var lastSyl int

var isLoaded = false

var dict dictionary

func Load(filename string) {

	file, _ := asset.Open(filename)

	scanner := bufio.NewReader(file)

	var maxSize int = 300000

	var fileData []byte = make([]byte, maxSize)

	size, _ := scanner.Read(fileData)

	fileData = fileData[:size]

	var buf bytes.Buffer = *bytes.NewBuffer(fileData)

	dec := gob.NewDecoder(&buf)

	err := dec.Decode(&dict)
	if err == nil {
		isLoaded = true
	}
	pos = make(map[int]int)
	end = make(map[int]bool)
	rewind()
}

func rewind() {
	for k := 2; k < 9; k++ {
		pos[k] = 0
		end[k] = false
	}
}

func Get(input string) *Pun {
	if !isLoaded {
		return &Pun{"Dictionary not loaded", ""}
	}

	input = strings.ToLower(input)
	if len(input) == 0 {
		return &Pun{"Input field is empty!", ""}
	}
	if len(input) < 2 {
		return &Pun{"Input too short", ""}
	}
	if len(input) > 7 {
		return &Pun{"Input too long", ""}
	}
	//if user enters a new word
	if input != lastInput {
		rewind()
		lastInput = input
	}

	var generatedPun string
	var basisWord string

	for {
		var endOfDict bool = true

		for k, v := range pos {
			if v < len(dict.Words[k]) {
				endOfDict = false
				//break
			}
		}
		if endOfDict {
			rewind()
			return &Pun{"end-of-dict", ""}
		}

		_, err := strconv.Atoi(input)

		if err == nil {
			return &Pun{"1337", input}
		}

		var random int = rand.Intn(82)

		if !end[2] && (((len(input) < 4) && (random < 10)) || ((len(input) >= 4) && (random < 3))) {
			// ( 10% / 3% ) generate from two syllables
			generatedPun, basisWord = searchForPun(input, 2)
			if generatedPun == "EOF" {
				end[2] = true
			}
		} else if !end[3] && (((len(input) < 4) && (random < 30)) || ((len(input) >= 4) && (random < 20))) {
			// ( 20% / 17% ) generate from three syllables
			generatedPun, basisWord = searchForPun(input, 3)
			if generatedPun == "EOF" {
				end[3] = true
			}
		} else if !end[4] && (((len(input) < 4) && (random < 50)) || ((len(input) >= 4) && (random < 45))) {
			// ( 20% / 25% ) generate from four syllables
			generatedPun, basisWord = searchForPun(input, 4)
			if generatedPun == "EOF" {
				end[4] = true
			}
		} else if !end[5] && (((len(input) < 4) && (random < 65)) || ((len(input) >= 4) && (random < 65))) {
			// ( 15% / 20% ) generate from five syllables
			generatedPun, basisWord = searchForPun(input, 5)
			if generatedPun == "EOF" {
				end[5] = true
			}
		} else if !end[6] && (((len(input) < 4) && (random < 75)) || ((len(input) >= 4) && (random < 75))) {
			// ( 10% / 10% ) generate from six syllables
			generatedPun, basisWord = searchForPun(input, 6)
			if generatedPun == "EOF" {
				end[6] = true
			}
		} else if !end[7] && (((len(input) < 4) && (random < 80)) || ((len(input) >= 4) && (random < 80))) {
			// ( 5% / 5% ) generate from seven syllables
			generatedPun, basisWord = searchForPun(input, 7)
			if generatedPun == "EOF" {
				end[7] = true
			}
		} else if !end[8] && (((len(input) < 4) && (random < 82)) || ((len(input) >= 4) && (random < 82))) {
			// ( 2% / 2% ) generate from eight syllables
			generatedPun, basisWord = searchForPun(input, 8)
			if generatedPun == "EOF" {
				end[8] = true
			}
		}
		if (generatedPun != "EOF") && (generatedPun != "") {
			break
		}
	}
	return &Pun{generatedPun, basisWord}
}

// generatedPun, dictWord
func searchForPun(input string, syl int) (string, string) {

	var generatedPun string
	var dictWord string

	for pos[syl] < len(dict.Words[syl]){

		if (pos[syl] == lastPos) && (pos[syl] < len(dict.Words[syl])) && (lastSyl == syl) {
		    return "EOF", "EOF"
		}

		lastPos = pos[syl]
		lastSyl = syl

		var syllableSlices []string = dict.Words[syl][pos[syl]]

		generatedPun = compareInputAndLine(input, syllableSlices)

        pos[syl]++

		dictWord = ""

		for i := 0; i < len(syllableSlices); i++ {
			dictWord += syllableSlices[i]
		}

		if generatedPun != "" {

			//remove occurances of three characters in a row
			for i := 0; i < len(generatedPun)-2; i++ {
				if (generatedPun[i] == generatedPun[i+1]) && (generatedPun[i] == generatedPun[i+2]) {
					generatedPun = generatedPun[:i] + generatedPun[i+1:]
				}
			}

			//removes duplicate letters in line/syllable joints
			var punSeparatedByInput []string = strings.Split(generatedPun, input)
			//if last of input equals first of dictionary
			if syllableSlices[0][0] == input[len(input)-1] {
				if len(punSeparatedByInput) == 3 {
					generatedPun = punSeparatedByInput[0] + input[:len(input)-1] + punSeparatedByInput[1]
					//removes last character
				}
			}
			return generatedPun, dictWord
		}
	}
	return "EOF", "EOF"
}
func compareInputAndLine(input string, line []string) string {
	var compareScore int = 0

	//len(line)-1 to ignore last syllable
	for s := 0; s < len(line)-1; s++ {
		var currentSyllable string = line[s]
		var nextSyllable string = line[s+1]
		var currentAndNextSyllable string = line[s] + line[s+1]

		//case 1: input is three characters shorter than currentSyllable
		if len(input) == len(currentSyllable)-3 {
			//case 1: input is three characters shorter than currentSyllable
			for i := 0; i < len(input); i++ {
				if input[i] == currentSyllable[i] {
					compareScore += 6
				}
				if input[i] == currentSyllable[i+1] {
					compareScore += 3
				}
			}

			if (compareScore >= 12) && (input != currentAndNextSyllable) {
				var approvedPun string
				for i := 0; i < len(line); i++ {
					if i == s {
						approvedPun += input + currentSyllable[len(currentSyllable)-3:len(currentSyllable)-1]
					} else {
						approvedPun += line[i]
					}
				}
				return approvedPun
			}
		} else if len(input) == len(currentSyllable)-2 {
			//case 2: input is two shorter than currentSyllable
			for i := 0; i < len(input); i++ {
				if input[i] == currentAndNextSyllable[i] {
					compareScore += 6
				}
				if input[i] == currentAndNextSyllable[i+1] {
					compareScore += 3
				}
			}
			if (compareScore >= 12) && (input != currentAndNextSyllable) {
				var approvedPun string
				for i := 0; i < len(line); i++ {
					if i == s {
						approvedPun += input + currentSyllable[len(currentSyllable)-2:len(currentSyllable)-1]
					} else {
						approvedPun += line[i]
					}
				}
				return approvedPun
			}

		} else if len(input) == len(currentSyllable)-1 {
			//case 3: input is one shorter than currentSyllable
			for i := 0; i < len(input); i++ {
				if input[i] == currentAndNextSyllable[i] {
					compareScore += 6
				}
				if input[i] == currentAndNextSyllable[i+1] {
					compareScore += 3
				}
			}
			if (input != currentAndNextSyllable) && ((compareScore >= 12) || ((len(input) == 2) && (compareScore >= 6))) {
				var approvedPun string

				for i := 0; i < len(line); i++ {
					if i == s {
						approvedPun += input + currentSyllable[:len(currentSyllable)-1]
					} else {
						approvedPun += line[i]
					}
				}
				return approvedPun
			}
		} else if len(input) == len(currentSyllable) {
			//case 4: input and currentSyllable have equal length
			for i := 0; i < len(input); i++ {
				if (input[i] == currentAndNextSyllable[i]) && (i < len(currentSyllable)+1) {
					compareScore += 6
				} else if (input[i] == currentAndNextSyllable[i]) && (i >= len(currentSyllable)+1) {
					compareScore += 4
				}
			}
			if (input != currentAndNextSyllable) && ((compareScore >= 12) || ((len(input) == 2) && (compareScore >= 6))) {
				var approvedPun string
				for i := 0; i < len(line); i++ {
					if i == s {
						approvedPun += input
					} else {
						approvedPun += line[i]
					}
				}
				return approvedPun
			}
		} else if len(input) == len(currentSyllable)+1 {
			//case 5: input is 1 longer than currentSyllable
			for i := 0; i < len(input); i++ {
				if (input[i] == currentAndNextSyllable[i]) && (i < len(currentSyllable)+1) {
					compareScore += 6
				} else if (input[i] == currentAndNextSyllable[i]) && (i >= len(currentSyllable)+1) {
					compareScore += 4
				}
			}
			if (compareScore >= 12) && (input != currentAndNextSyllable) {
				var approvedPun string
				for i := 0; i < len(line); i++ {
					if i == s {
						approvedPun += input
					} else {
						approvedPun += line[i]
					}
				}
				return approvedPun
			}
		} else if (len(input) == len(currentSyllable)+2) && (len(line) >= 4) {
			//case 6: input is 2 longer than currentSyllable and line has 4 or more syllables
			if len(nextSyllable) == 3 {
				//case 6.1: input is 1 shorter than currentAndNextSyllable
				for i := 0; i < len(currentAndNextSyllable)-1; i++ {
					if (input[i] == currentAndNextSyllable[i]) && (i < len(currentSyllable)+1) {
						compareScore += 6
					} else if (input[i] == currentAndNextSyllable[i]) && (i >= len(currentSyllable)+1) {
						compareScore += 4
					}
				}
			} else if len(nextSyllable) == 2 {
				//case 6.2: input and currentAndNextSyllable has equal length
				for i := 0; i < len(currentAndNextSyllable); i++ {
					if (input[i] == currentAndNextSyllable[i]) && (i < len(currentSyllable)+1) {
						compareScore += 6
					} else if (input[i] == currentAndNextSyllable[i]) && (i >= len(currentSyllable)+1) {
						compareScore += 4
					}
				}
			}
			if (compareScore >= 12) && (input != currentAndNextSyllable) {
				var approvedPun string
				for i := 0; i < len(line); i++ {
					if i == s {
						approvedPun += input
					} else {
						approvedPun += line[i]
					}
				}
				return approvedPun
			}
		} else if (len(input) == len(currentSyllable)+3) && (len(line) >= 4) {
			//case 7: input is 3 longer than currentSyllable and line has 4 or more syllables

			if len(nextSyllable) == 4 {
				// case 7.1: Input is one shorter than currentAndNextSyllable
				for i := 0; i < len(currentAndNextSyllable)-1; i++ {

					if (input[i] == currentAndNextSyllable[i]) && (i < len(currentSyllable)+1) {
						compareScore += 6
					} else if (input[i] == currentAndNextSyllable[i]) && (i >= len(currentSyllable)+1) {
						compareScore += 4
					}
				}
			} else if (len(nextSyllable) == 3) || (len(nextSyllable) == 2) {
				//  case 7.2: input is as long as currentAndNextSyllable or is one longer
				for i := 0; i < len(currentAndNextSyllable); i++ {

					if (input[i] == currentAndNextSyllable[i]) && (i < len(currentSyllable)+1) {
						compareScore += 6
					} else if input[i] == currentAndNextSyllable[i] && (i >= len(currentSyllable)+1) {
						compareScore += 4
					}
				}
			}

			if (compareScore >= 12) && (input != currentAndNextSyllable) {
				var dictword string

				for i := 0; i < len(line); i++ {
					dictword += line[i]
				}

				var approvedPun string

				// Replaces one syllable for shorter words, two syllables for longer words
				for i := 0; i < len(line); i++ {
					if (i == s) && (((len(dictword) - len(currentAndNextSyllable)) < 6) || (!((len(dictword) - len(currentAndNextSyllable)) < 6)) && (s >= 3)) {
						approvedPun += input
					} else if (i == s) && !((len(dictword) - len(currentAndNextSyllable)) < 6) {
						approvedPun += input
						i++
					} else {
						approvedPun += line[i]
					}
				}

				return approvedPun
			}
		} else if (len(input) == len(currentSyllable)+4) && (len(line) >= 5) {
			//case 8: input is 4 longer than currentSyllable and line has 5 or more syllables
			if len(nextSyllable) == 5 {
				// case 8.1: Input is one shorter than currentAndNextSyllable
				for i := 0; i < len(currentAndNextSyllable)-1; i++ {

					if (input[i] == currentAndNextSyllable[i]) && (i < len(currentSyllable)+1) {
						compareScore += 6
					} else if (input[i] == currentAndNextSyllable[i]) && (i >= len(currentSyllable)+1) {
						compareScore += 4
					}
				}
			} else if (len(nextSyllable) == 3) || (len(nextSyllable) == 4) {
				// case 8.2: input is as long as currentAndNextSyllable or is one longer
				for i := 0; i < len(currentAndNextSyllable); i++ {

					if (input[i] == currentAndNextSyllable[i]) && (i < len(currentSyllable)+1) {
						compareScore += 6
					} else if input[i] == currentAndNextSyllable[i] && (i >= len(currentSyllable)+1) {
						compareScore += 4
					}
				}
			}
			if (compareScore >= 12) && (input != currentAndNextSyllable) {
				var dictword string

				for i := 0; i < len(line); i++ {
					dictword += line[i]
				}

				var approvedPun string

				// Replaces one syllable for shorter words, two syllables for longer words if s < 3
				for i := 0; i < len(line); i++ {
					if (i == s) && (((len(dictword) - len(currentAndNextSyllable)) < 6) || (!((len(dictword) - len(currentAndNextSyllable)) < 6)) && (s >= 3)) {
						approvedPun += input
					} else if (i == s) && !((len(dictword) - len(currentAndNextSyllable)) < 6) && (s < 3) {
						approvedPun += input
						i++
					} else {
						approvedPun += line[i]
					}
				}

				return approvedPun
			}
		}
		compareScore = 0
	}

	return ""

}
