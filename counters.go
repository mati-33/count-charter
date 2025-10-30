package main

import (
	"regexp"
	"slices"
)

type CharsCounter struct {
	toSkip []string
}

func NewCharsCounter() *CharsCounter {
	return &CharsCounter{toSkip: []string{"\t", "\n", " "}}
}

func (c CharsCounter) Count(line string, ct *CountTracker) {
	for _, r := range line {
		if slices.Contains(c.toSkip, string(r)) {
			continue
		}
		ct.Increment(string(r))
	}
}

type LettersCounter struct{}

type PunctuationCounter struct{}

type WordsCounter struct {
	wordRegexp *regexp.Regexp
}

func NewWordsCounter() *WordsCounter {
	return &WordsCounter{wordRegexp: regexp.MustCompile("[a-zA-Z]{2,}")}
}

func (c WordsCounter) Count(line string, ct *CountTracker) {
	if len(line) != 0 {
		words := c.wordRegexp.FindAllString(line, -1)
		for _, w := range words {
			ct.Increment(w)
		}
	}
}

type GoKeywordsCounter struct {
	wordRegexp *regexp.Regexp
	keywords   []string
}

func NewGoKeywordsCounter() *GoKeywordsCounter {
	return &GoKeywordsCounter{
		wordRegexp: regexp.MustCompile("[a-zA-Z]{2,}"),
		keywords: []string{
			"break", "default", "func", "interface", "select",
			"case", "defer", "go", "map", "struct",
			"chan", "else", "goto", "package", "switch",
			"const", "fallthrough", "if", "range", "type",
			"continue", "for", "import", "return", "var",
		}}
}

func (c GoKeywordsCounter) Count(line string, ct *CountTracker) {
	if len(line) != 0 {
		words := c.wordRegexp.FindAllString(line, -1)
		for _, w := range words {
			if !slices.Contains(c.keywords, w) {
				continue
			}
			ct.Increment(w)
		}
	}
}
