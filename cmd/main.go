package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	path := flag.String("path", ".", "Path to scan")
	dryRun := flag.Bool("dry-run", false, "Dry run")
	flag.Parse()

	fmt.Printf("Path: %s\n", *path)

	scanCount := 0
	processCount := 0

	err := filepath.WalkDir(*path, func(path string, d os.DirEntry, err error) error {
		scanCount++
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			if err != nil {
				panic(err)
			}

			return filepath.SkipDir
		}

		processed := mergeKoreanLetters(d.Name())

		if processed != d.Name() {
			processCount++
			fmt.Printf("Rename: %s -> %s\n", d.Name(), processed)

			if *dryRun {
				return nil
			}

			err := os.Rename(path, filepath.Join(filepath.Dir(path), processed))

			if err != nil {
				_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				if err != nil {
					panic(err)
				}
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Scan count: %d\n", scanCount)
	fmt.Printf("Process count: %d\n", processCount)
}

func getIndex(r rune, start rune) int {
	return int(r - start)
}

func mergeKoreanLetters(input string) string {
	initialStart := rune(0x1100)
	medialStart := rune(0x1161)
	finalStart := rune(0x11A8)

	initialEnd := rune(0x1112)
	medialEnd := rune(0x1175)
	finalEnd := rune(0x11C2)

	runes := []rune(input)
	result := strings.Builder{}

	for i := 0; i < len(runes); {
		if runes[i] >= initialStart && runes[i] <= initialEnd { // Initial consonant
			initialIndex := getIndex(runes[i], initialStart)
			i++
			if i < len(runes) && runes[i] >= medialStart && runes[i] <= medialEnd { // Medial vowel
				medialIndex := getIndex(runes[i], medialStart)
				i++
				finalIndex := 0
				if i < len(runes) && runes[i] >= finalStart && runes[i] <= finalEnd { // Final consonant (optional)
					finalIndex = getIndex(runes[i], finalStart) + 1
					i++
				}
				syllable := 0xAC00 + (initialIndex * 21 * 28) + (medialIndex * 28) + finalIndex
				result.WriteRune(rune(syllable))
			}
		} else {
			result.WriteRune(runes[i])
			i++
		}
	}

	return result.String()
}
