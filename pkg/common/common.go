package common

import (
	"encoding/json"
	"os"
	"io"
	"bytes"
	"math"
    "errors"
	"strings"
	"time"
	"fmt"
)
var (
	Quit = make(chan []byte)
	// Supported interval units and their time.Duration equivalent
	unitMap = map[string]float64{
		"s":       float64(time.Second),
		"sec":     float64(time.Second),
		"secs":    float64(time.Second),
		"second":  float64(time.Second),
		"seconds": float64(time.Second),
		"m":       float64(time.Minute),
		"min":     float64(time.Minute),
		"mins":     float64(time.Minute),
		"minute":  float64(time.Minute),
		"minutes": float64(time.Minute),
		"h":       float64(time.Hour),
		"hr":      float64(time.Hour),
		"hrs":      float64(time.Hour),
		"hour":    float64(time.Hour),
		"hours":   float64(time.Hour),
		"d":       float64(time.Hour * 24),
		"day":     float64(time.Hour * 24),
		"days":    float64(time.Hour * 24),
	}
)

// Function to convert structs to maps. This requies
// json field tags to be set.  Missing json tags
// will result in field not getting added.
func StructToMap(i interface{}) (smap map[string]interface{}) {
	inrec, _ := json.Marshal(i)
	json.Unmarshal(inrec, &smap)
	return
}

// Function to parse a string into a time.Duration
func ParseDuration(s string) (time.Duration, error) {
    if len(s) == 0 {
            return 0, nil
    }
    var exp, whole, fraction int64
    var number, totalDuration float64

    // consume digits 
    var done bool
    for !done && len(s) > 0 {
        c := s[0]
        switch {
        case c >= '0' && c <= '9':
                d := int64(c - '0')
                if exp > 0 {
                        exp++
                        fraction = 10*fraction + d
                } else {
                        whole = 10*whole + d
                }
                s = s[1:]
        case c == '.':
                if exp > 0 {
                        return 0, fmt.Errorf("invalid floating point number format: two decimal points found")
                }
                exp = 1
                fraction = 0
                s = s[1:]
        default:
                done = true
        }
    }
    // adjust number
    number = float64(whole)
    if exp > 0 {
            number += float64(fraction) * math.Pow(10, float64(1-exp))
    }

    // find end of unit
    var i int
    for ; i < len(s) && s[i] != '+' && s[i] != '-' && (s[i] < '0' || s[i] > '9'); i++ {
            // identifier bytes: no-op
    }
    unit := strings.TrimSpace(s[:i])
    // fmt.Printf("number: %f; unit: %q\n", number, unit)
    if duration, ok := unitMap[unit]; ok {
            totalDuration += number * duration
    } else {
            if unit == "" {
                return 0, errors.New("duration missing units")
            } else {
                return 0, fmt.Errorf("unrecognized unit in duration: %q", unit)
            }
    }
    return time.Duration(totalDuration), nil
}

func CountLines(fileName string) (int) {
	file, err := os.Open(fileName)

	if err != nil {
		return 0
	}

	buf := make([]byte, 1024)
	lines := 0

	for {
		readBytes, err := file.Read(buf)
		if err != nil {
			if readBytes == 0 && err == io.EOF {
				err = nil
			}
			return lines
		}
		lines += bytes.Count(buf[:readBytes], []byte{'\n'})
	}
	return lines
}

// returns a list of files found in the given directory
func GetFileList(dir string) (files []string) {
	f, err := os.Open(dir)
	if err != nil {
		return nil
	}
	dirlist, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil
	}
	for _, e := range dirlist {
		files = append(files, e.Name())
	}
	return
}

