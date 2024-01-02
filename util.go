package quartz

import (
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func sortLevels(levels []QuartzLevel) []QuartzLevel {
	sort.Slice(levels, func(i, j int) bool {
		if levels[i] == QUARTZ_LVL_V {
			return true
		}
		if levels[j] == QUARTZ_LVL_V {
			return false
		}
		return string(levels[i]) < string(levels[j])
	})
	return levels
}

func levelsToString(levels []QuartzLevel) string {
	level := ""
	for _, lvl := range levels {
		level += string(lvl)
	}
	return level
}

func sortLevelsToString(levels []QuartzLevel) string {
	levels = sortLevels(levels)
	return levelsToString(levels)
}

func parseLevels(levels string) []QuartzLevel {
	result := make([]QuartzLevel, 0)
	for i := 0; i < len(levels); i++ {
		result = append(result, QuartzLevel(levels[i]))
	}

	return sortLevels(result)
}

func parseResponse(line string) QuartzResponse {
	switch string(line[1]) {
	case "A":
		if len(line) == 3 {
			return &ResponseAcknowledge{RawData: line}
		}
		// Route list, send as an update message for simplicity
		levelsEndIdx := -1
		for i := 2; i < len(line); i++ {
			if unicode.IsDigit(rune(line[i])) {
				levelsEndIdx = i
				break
			}
		}
		commaIdx := strings.Index(line, ",")
		if commaIdx == -1 || levelsEndIdx == -1 {
			return nil
		}
		levels := parseLevels(line[2:levelsEndIdx])
		dest, err := strconv.Atoi(line[levelsEndIdx:commaIdx])
		if err != nil {
			return nil
		}
		src, err := strconv.Atoi(line[commaIdx+1 : len(line)-1])
		if err != nil {
			return nil
		}
		return &ResponseUpdate{
			RawData:     line,
			Levels:      levels,
			Destination: uint(dest),
			Source:      uint(src),
		}
	case "E":
		return &ResponseError{RawData: line}
	case "P":
		return &ResponsePowerOn{RawData: line}
	case "U":
		levelsEndIdx := -1
		for i := 2; i < len(line); i++ {
			if unicode.IsDigit(rune(line[i])) {
				levelsEndIdx = i
				break
			}
		}
		commaIdx := strings.Index(line, ",")
		if commaIdx == -1 || levelsEndIdx == -1 {
			return nil
		}
		levels := parseLevels(line[2:levelsEndIdx])
		dest, err := strconv.Atoi(line[levelsEndIdx:commaIdx])
		if err != nil {
			return nil
		}
		src, err := strconv.Atoi(line[commaIdx+1 : len(line)-1])
		if err != nil {
			return nil
		}
		return &ResponseUpdate{
			RawData:     line,
			Levels:      levels,
			Destination: uint(dest),
			Source:      uint(src),
		}
	case "R":
		switch string(line[3]) {
		case "D":
			// Destination
			commaIdx := strings.Index(line, ",")
			if commaIdx == -1 {
				return nil
			}
			dest, err := strconv.Atoi(line[4:commaIdx])
			if err != nil {
				return nil
			}
			name := line[commaIdx+1 : len(line)-1]
			return &ResponseReadDestination{
				RawData:     line,
				Destination: uint(dest),
				Name:        name,
			}
		case "S":
			// Source
			commaIdx := strings.Index(line, ",")
			if commaIdx == -1 {
				return nil
			}
			src, err := strconv.Atoi(line[4:commaIdx])
			if err != nil {
				return nil
			}
			name := line[commaIdx+1 : len(line)-1]
			return &ResponseReadSource{
				RawData: line,
				Source:  uint(src),
				Name:    name,
			}
		case "L":
			// Level
			commaIdx := strings.Index(line, ",")
			if commaIdx == -1 {
				return nil
			}
			levels := parseLevels(line[4:commaIdx])
			name := line[commaIdx+1 : len(line)-1]
			return &ResponseReadLevel{
				RawData: line,
				Level:   levels[0],
				Name:    name,
			}
		}
	case "B":
		// Lock status
		commaIdx := strings.Index(line, ",")
		if commaIdx == -1 {
			return nil
		}
		dest, err := strconv.Atoi(line[3:commaIdx])
		if err != nil {
			return nil
		}
		lockVal, err := strconv.Atoi(line[commaIdx+1 : len(line)-1])
		if err != nil {
			return nil
		}
		return &ResponseLockStatus{
			RawData:     line,
			Destination: uint(dest),
			Locked:      lockVal == 0,
		}
	}
	return nil
}
