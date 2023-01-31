// Package uti provides general, low priority methods and functions for different purposes
package uti

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Nesvilab/philosopher/lib/msg"
)

// Round serves the rol of the missing math.Round function
func Round(val float64, roundOn float64, places int) (newVal float64) {

	var round float64

	pow := math.Pow(10, float64(places))
	digit := pow * val

	_, div := math.Modf(digit)
	_div := math.Copysign(div, val)
	_roundOn := math.Copysign(roundOn, val)

	if _div >= _roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}

	newVal = round / pow

	return
}

// ToFixed truncates float numbers to a specific position
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(toFixedRound(num*output)) / output
}

func toFixedRound(num float64) int {
	return int(num + math.Copysign(0.05, num))
}

// ParseFloat converts scientific notation values from string format to float64
func ParseFloat(str string) (float64, error) {

	val, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return val, nil
	}

	//Some number may be separated by comma, for example, 23,120,123, so remove the comma firstly
	str = strings.Replace(str, ",", "", -1)

	//Some number is specified in scientific notation
	pos := strings.IndexAny(str, "eE")
	if pos < 0 {
		return strconv.ParseFloat(str, 64)
	}

	var baseVal float64
	var expVal int64

	baseStr := str[0:pos]
	baseVal, err = strconv.ParseFloat(baseStr, 64)
	if err != nil {
		return 0, err
	}

	expStr := str[(pos + 1):]
	expVal, err = strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return baseVal * math.Pow10(int(expVal)), nil
}

// GetLabelNames add custom names adds to the label structures user-defined names to be used on the isobaric labels
func GetLabelNames(annot string) map[string]string {

	var labels = make(map[string]string)

	file, e := os.Open(annot)
	if e != nil {
		msg.ReadFile(e, "error")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// does the line has at least an iso tag?
		if len(scanner.Text()) >= 3 {

			names := strings.Fields(scanner.Text())
			labels[names[0]] = names[1]
		}
	}

	if e = scanner.Err(); e != nil {
		msg.ReadFile(e, "error")
	}

	return labels
}

// FindFile locates a file based on a name pattern
func FindFile(targetDir string, pattern string) string {

	match, e := filepath.Glob(targetDir + string(filepath.Separator) + pattern)

	if e != nil {
		fmt.Println(e)
	}

	return match[0]
}

// WalkMatch recursively looks for files with a certain extension in a specific folder
func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// IOReadDir looks for files with a certain extension in a specific folder
func IOReadDir(root, ext string) []string {

	var files []string

	fileInfo, err := os.Open(root)
	if err != nil {
		return files
	}

	fileList, _ := fileInfo.Readdir(0)

	for _, file := range fileList {
		if strings.Contains(file.Name(), ext) {
			files = append(files, filepath.Join(root, file.Name()))
		}
	}

	return files
}

// GetMaxNumber returns the highest (string) number from an array in string format
func GetMaxNumber(list []string) string {

	var max = -1.0

	for _, i := range list {
		f, _ := strconv.ParseFloat(i, 64)
		if f > max {
			max = f
		}
	}

	s := fmt.Sprintf("%f", max)

	if s == "-1.000000" {
		return ""
	}

	return s
}

// RemoveDuplicateStrings removes duplicates from a slice
func RemoveDuplicateStrings(slice []string) []string {
	keys := make(map[string]struct{}, len(slice))
	list := make([]string, 0, len(slice))
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = struct{}{}
			list = append(list, entry)
		}
	}
	list2 := make([]string, len(list))
	copy(list2, list)
	return list2
}
