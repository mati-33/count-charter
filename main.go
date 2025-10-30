package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

type CountTracker struct {
	targetToCount map[string]int
	allCounts     int
}

func (ct *CountTracker) Increment(value string) {
	ct.targetToCount[value]++
	ct.allCounts++
}

func (ct CountTracker) AllCounts() int {
	return ct.allCounts
}

func (ct CountTracker) TargetToCount() map[string]int {
	return ct.targetToCount
}

type Counter interface {
	Count(line string, ct *CountTracker)
}

type Config struct {
	limit          int
	barChar        ColumnCharacter
	chartWidth     int
	showPercentage bool
	timeStep       int
}

func NewDefaultConfig() *Config {
	return &Config{
		limit:          25,
		barChar:        "█",
		chartWidth:     36,
		showPercentage: false,
		timeStep:       30,
	}
}

func main() {
	config := NewDefaultConfig()
	flag.Usage = Usage
	target := ChartTarget("chars")

	flag.Var(&target, "target", "target todo")
	flag.IntVar(&config.limit, "limit", config.limit, "limit todo")
	flag.Var(&config.barChar, "char", "chat todo")
	flag.IntVar(&config.chartWidth, "width", config.chartWidth, "width todo")
	flag.BoolVar(&config.showPercentage, "show-percentage", config.showPercentage, "show-percentage todo")
	flag.IntVar(&config.timeStep, "step", config.timeStep, "step todo")

	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		exitWithErrMessage("required file name")
	}

	filepath := flag.Arg(0)
	f, err := os.Open(filepath)
	if err != nil {
		exitWithErrMessage("could not open the file")
	}
	defer f.Close()

	countTracker := CountTracker{map[string]int{}, 0}
	scanner := bufio.NewScanner(f)
	var counter Counter

	switch target {
	case "chars":
		counter = NewCharsCounter()
	case "words":
		counter = NewWordsCounter()
	case "go":
		counter = NewGoKeywordsCounter()
	}

	var nlines int
	fmt.Printf("\n")
	for scanner.Scan() {
		counter.Count(scanner.Text(), &countTracker)
		chart, chartLen := generateChart(countTracker, *config)

		createTerminalSpace(chartLen + 1)
		saveCursorPosition()
		eraseTerminalBottom()

		fmt.Println(chart)
		time.Sleep(time.Duration(config.timeStep) * time.Millisecond)

		moveCursorToSavedPosition()
		nlines = chartLen
	}
	moveCursorDown(nlines)
}

func generateChart(ct CountTracker, config Config) (string, int) {
	targetToCount := ct.TargetToCount()
	targets := make([]string, 0, len(targetToCount))
	for t := range targetToCount {
		targets = append(targets, t)
	}
	sort.Slice(targets, func(i, j int) bool {
		return targetToCount[targets[i]] > targetToCount[targets[j]]
	})

	limit := min(config.limit, len(targets))
	if limit < 0 {
		limit = len(targets)
	}
	targets = targets[:limit]

	var b strings.Builder
	targetsLen := longestTargetLen(targets)
	var ratio float64
	for i, t := range targets {
		perc := float64(targetToCount[t]) / float64(ct.AllCounts())
		if i == 0 {
			ratio = float64(config.chartWidth) / perc
		}
		spaces := targetsLen - utf8.RuneCountInString(t)
		fmt.Fprintf(&b, "%s %s ", strings.Repeat(" ", spaces), t)
		fmt.Fprintf(&b, strings.Repeat(string(config.barChar), int(ratio*perc)))
		if config.showPercentage {
			fmt.Fprintf(&b, " %.2f%%", 100.00*perc)
		}
		fmt.Fprint(&b, "\n")
	}
	chart := b.String()
	return chart, len(targets)
}

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS] filename\n", os.Args[0])
	fmt.Println("OPTIONS:")
	flag.PrintDefaults()
}

func longestTargetLen(targets []string) int {
	longestLen := 0
	for _, t := range targets {
		if tlen := utf8.RuneCountInString(t); tlen > longestLen {
			longestLen = tlen
		}
	}
	return longestLen
}

func createTerminalSpace(nlines int) {
	fmt.Print(strings.Repeat("\n", nlines))
	fmt.Printf("\033[%dA", nlines)
}

func saveCursorPosition() { fmt.Print("\033[s") }

func eraseTerminalBottom() { fmt.Printf("\033[0J") }

func moveCursorToSavedPosition() { fmt.Print("\033[u") }

func moveCursorDown(nlines int) { fmt.Printf("\033[%dB", nlines) }

func exitWithErrMessage(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
