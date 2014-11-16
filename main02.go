package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Radicals map[string][]string	// "的":["白","勹","丶"], ...
type Sounds map[string]string		// "a":"A阿吖呵啊嗄腌锕", ...

/*
	(1)	What we are calling "radicals" here might more accurately be called "components" (or, what do linguists call them?)? There is usually only one *radical* to any character, and that is its search key. Hence perhaps the massive number of search results??
	
	(2)	What we probably ultimately need is something more hierarchical, e.g.,:
	
	"好":["女":["?","一","丿"],["子":"了":["㇇","亅"]，"一"]],
		where there are always (always?) only two top-level key-value pairs, the first being by default the *radical* (i.e., 好's dictionary root 女) and the second the *complement* (i.e., the full component that is not the radical), gradually broken down hierarchically (ultimately finishing with each subcomponent's brushstrokes)
	
	Is there anywhere that we can find something like this?? Or should we start building it? Not such a difficult job.
*/

func loadRadicals(file string) (Radicals, error) {
	r := Radicals{}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(b, &r)
	return r, err
}

func loadSounds(file string) (Sounds, error) {
	s := Sounds{}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(b, &s)
	return s, err
}

func main() {
	// characters => radicals
	r, err := loadRadicals("radicals.json")
	if err != nil {
		panic(err)
	}

	// pinyin => characters
	s, err := loadSounds("sounds.json")
	if err != nil {
		panic(err)
	}
	
	// radical => characters
	radicalLookup := map[string]string{}
	for ch, rads := range r {
		for _, rad := range rads {
			if !strings.Contains(radicalLookup[rad], ch) {
				radicalLookup[rad] += ch
			}
		}
	}

	radicals := []string{}
	usedRads := map[string]string{} // character => radicals used in filter
	possibleChars := ""
	i := 0
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("PASS %d Enter pinyin: ",i+1)
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		text = strings.Trim(text, "\r")					// Windows/DOS CR
		fmt.Print(text)
		
		if chars, ok := s[text]; ok {
			radicals = append(radicals, "")
			j := 0
			for _, c := range chars {					// for each of this pinyin's possible chars...
				j++
				fmt.Printf("\nc%d = %d (%s)", j, c, string(c))
				k := 0
				for _, rads := range r[string(c)] {		// get the radical(s) of each of its possible chars...
					k++
					fmt.Printf(" : r%d = %s", k, rads)
					for _, rad := range rads {
						if !strings.Contains(strings.Join(radicals, ""), string(rad)) {
							radicals[i] += string(rad)
						}
						if i == 0 {
							allPossibleChars := radicalLookup[string(rad)]
							possibleChars += allPossibleChars
							for _, apc := range allPossibleChars {
								usedRads[string(apc)] = string(rad)
							}
						}
					}
				}
			}

			// reduce number of possibilities
			if i >= 1 {
				remainingChars := ""
				for _, c := range possibleChars {
					foundAny := false
					// get radicals for possible character
					chRads := r[string(c)]
				filter:
					for _, rad := range radicals[i] {
						// see if character has radical in latest searched list of radicals
						for _, cr := range chRads {
							if cr == string(rad) && !strings.Contains(usedRads[string(c)], string(rad)) {
								foundAny = true
								usedRads[string(c)] = string(rad)
								break filter
							}
						}
					}
					if foundAny {
						remainingChars += string(c)
					}
				}
				possibleChars = remainingChars
			}

			fmt.Printf("\nAbout to print %d characters\n", len(possibleChars)/3)
			fmt.Println(radicals)
			fmt.Println(possibleChars)
			fmt.Printf("\n")

			i++
		}

		if text == "" {
			break
		}
	}
}
