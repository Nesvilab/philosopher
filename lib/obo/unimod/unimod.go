package obo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/prvst/philosopher/lib/err"
)

// UniMod is the top level struct
type UniMod []Mod

// Mod is the top level struct
type Mod struct {
	ID               string  `toml:"ID"`
	Name             string  `toml:"Name"`
	Definition       string  `toml:"Definition"`
	TermID           int     `toml:"TermID"`
	MonoIsotopicMass float64 `toml:"MonoIsotopicMass"`
	AverageMass      float64 `toml:"AverageMass"`
	Composition      string  `toml:"Composition"`
	DateTimePosted   string  `toml:"DateTimePosted"`
	DateTimeModified string  `toml:"DateTimeModified"`
	IsA              string  `toml:"IsA"`
}

// Parse is the function that will read an OBO file and return the filled structs
func (o *UniMod) Parse(f string) *err.Error {

	oboFile, e := os.Open(f)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: filepath.Base(f)}
	}
	defer oboFile.Close()

	scanner := bufio.NewScanner(oboFile)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}
