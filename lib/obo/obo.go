package obo

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/met"
	unmd "github.com/prvst/philosopher/lib/obo/unimod"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/vmihailenco/msgpack"
)

// DataFormat defines different data type from PSI
type DataFormat interface {
	Deploy() *err.Error
	Serialize() *err.Error
	Restore() *err.Error
}

// Onto contains Ontology terms
type Onto struct {
	met.Data
	OboFile string
	Version string
	Date    string
	Terms   []Term
}

// Term refers to an atomic ontology definition
type Term struct {
	ID               string
	RecordID         int
	Name             string
	Definition       string
	DateTimePosted   string
	DateTimeModified string
	Comments         string
	Synonyms         string
	IsA              string
	MonoIsotopicMass float64
	AverageMass      float64
	Composition      string
}

// NewUniModOntology constructucst a set of UniMod ontologies
func NewUniModOntology() (Onto, *err.Error) {

	var m met.Data
	var o Onto
	var e *err.Error

	m.Restore(sys.Meta())

	o.UUID = m.UUID
	o.Distro = m.Distro
	o.Home = m.Home
	o.MetaFile = m.MetaFile
	o.MetaDir = m.MetaDir
	o.DB = m.DB
	o.Temp = m.Temp
	o.TimeStamp = m.TimeStamp

	// Deploy
	e = o.Deploy()
	if e != nil {
		return o, e
	}

	// Read
	o.Parse()

	// Serielize
	o.Serialize()

	return o, nil
}

// Deploy deploys the OBO file to the temp folder
func (m *Onto) Deploy() *err.Error {

	m.OboFile = fmt.Sprintf("%s%sunimod.obo", m.Temp, string(filepath.Separator))

	unmd.Deploy(m.OboFile)

	return nil
}

// Parse reads the unimod.obo file and creates the data structure
func (m *Onto) Parse() *err.Error {

	file, err := os.Open(m.OboFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var flag = -1
	var term Term

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		if strings.HasPrefix(scanner.Text(), "format-version:") {
			m.Version = splitAndCollect(scanner.Text(), false)
		}

		if strings.HasPrefix(scanner.Text(), "date:") {
			m.Date = splitAndCollect(scanner.Text(), false)
		}

		if strings.Contains(scanner.Text(), "[Term]") {
			flag = 1
		}

		if strings.HasPrefix(scanner.Text(), "id:") && flag == 1 {
			term.ID = splitAndCollect(scanner.Text(), false)
		} else if strings.HasPrefix(scanner.Text(), "name:") && flag == 1 {
			term.Name = splitAndCollect(scanner.Text(), false)
		} else if strings.HasPrefix(scanner.Text(), "def:") && flag == 1 {
			term.Definition = splitAndCollect(scanner.Text(), false)
		} else if strings.HasPrefix(scanner.Text(), "comment:") && flag == 1 {
			term.Comments = splitAndCollect(scanner.Text(), false)
		} else if strings.HasPrefix(scanner.Text(), "synonym:") && flag == 1 {
			term.Synonyms = splitAndCollect(scanner.Text(), false)
		}

		if strings.HasPrefix(scanner.Text(), "comment") && flag == 1 {
			term.Comments = splitAndCollect(scanner.Text(), false)
		}

		if strings.HasPrefix(scanner.Text(), "synonym") && flag == 1 {
			term.Synonyms = splitAndCollect(scanner.Text(), false)
		}

		if strings.HasPrefix(scanner.Text(), "xref: record_id") && flag == 1 {
			i, _ := strconv.Atoi(splitAndCollect(scanner.Text(), true))
			term.RecordID = i
		} else if strings.HasPrefix(scanner.Text(), "xref: delta_mono_mass") && flag == 1 {
			i, _ := strconv.ParseFloat(splitAndCollect(scanner.Text(), true), 64)
			term.MonoIsotopicMass = i
		} else if strings.HasPrefix(scanner.Text(), "xref: delta_avge_mass") && flag == 1 {
			i, _ := strconv.ParseFloat(splitAndCollect(scanner.Text(), true), 64)
			term.AverageMass = i
		} else if strings.HasPrefix(scanner.Text(), "xref: delta_composition") && flag == 1 {
			term.Composition = splitAndCollect(scanner.Text(), true)
		} else if strings.HasPrefix(scanner.Text(), "xref: date_time_posted") && flag == 1 {
			term.DateTimePosted = splitAndCollect(scanner.Text(), true)
		} else if strings.HasPrefix(scanner.Text(), "xref: date_time_modified") && flag == 1 {
			term.DateTimeModified = splitAndCollect(scanner.Text(), true)
		}

		if strings.HasPrefix(scanner.Text(), "is_a") && flag == 1 {
			term.IsA = splitAndCollect(scanner.Text(), false)
		}

		if strings.EqualFold(scanner.Text(), "//") && flag == 1 {
			flag = 0
			m.Terms = append(m.Terms, term)
			term = Term{}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	spew.Dump(m)
	os.Exit(1)

	return nil
}

// Serialize UniMod data structure
func (m Onto) Serialize() *err.Error {

	b, er := msgpack.Marshal(&m)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.MODBin(), b, sys.FilePermission())
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (m *Onto) Restore() *err.Error {

	b, e := ioutil.ReadFile(sys.MODBin())
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: ": database data may be corrupted"}
	}

	e = msgpack.Unmarshal(b, &m)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: ": database data may be corrupted"}
	}

	return nil
}

func splitAndCollect(s string, isXref bool) string {

	if isXref == false {
		l := strings.Split(s, ": ")
		return l[1]
	}

	l := strings.Split(s, "\"")
	return l[1]
}
