/*
	Experiment "main03" [pz] - Finding chars with most in common

	Wish list:
		- start with pinyin group
		- based on several criteria, make "cascading neighborhoods" of each char, then subsets, to connect each char to criterial "neighbors"
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"bufio"
	"os"
	"strings"
)

type Radicals map[string][]string	// "的":["白","勹","丶"], ...
type Sounds map[string]string		// "a":"A阿吖呵啊嗄腌锕", ...

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

	// characters => pinyin
	characterLookup := map[string]string{}
	for sound, chars := range s {
		for _, c := range chars {
			characterLookup[string(c)] = sound
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter character: ")
	obChar, _ := reader.ReadString('\n')
	obChar = strings.Trim(obChar, "\n")
	obChar = strings.Trim(obChar, "\r")					// Windows/DOS CR
	fmt.Printf("%s (%s)\n", obChar, characterLookup[obChar])
	com := map[string]int{}
	
	obRads, ok := r[obChar]
	if ok {
	
		// tally each char with at least one component (radical) in common
		for _, rad := range obRads {
			for _, ch := range radicalLookup[rad] {
				com[string(ch)]++
			}
		}

		// cheapsort chars by frequency of common components
		// needs a smarter sorting technique
		i := 10
		for i > 0 {
			seen := false
			for ch, j := range com {
				if (i == j) {
					if !seen {
						seen = true
						fmt.Printf("%d in common: ", i)
					}
					fmt.Printf("%s%s", string(ch), characterLookup[string(ch)])
				}	
			}
			if seen { fmt.Println() }
			i--
		}
	}
	
}	// main()