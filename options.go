package main

import (
	"fmt"
	"unicode/utf8"
)

type ChartTarget string

func (t *ChartTarget) String() string {
	return string(*t)
}
func (t *ChartTarget) Set(v string) error {
	switch v {
	case "ch", "chars":
		*t = "chars"
	case "w", "words":
		*t = "words"
	case "go":
		*t = "go"
	default:
		return fmt.Errorf("valid values are: 'ch', 'w', 'go'")
	}
	return nil
}

type ColumnCharacter string

func (c *ColumnCharacter) String() string {
	return string(*c)
}

func (c *ColumnCharacter) Set(v string) error {
	if utf8.RuneCountInString(v) > 1 {
		return fmt.Errorf("to many characters")
	}
	*c = ColumnCharacter(v)
	return nil
}
