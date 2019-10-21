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

	"github.com/nesvilab/philosopher/lib/msg"

	"github.com/nesvilab/philosopher/lib/met"
	unmd "github.com/nesvilab/philosopher/lib/obo/unimod"
	"github.com/nesvilab/philosopher/lib/sys"
	"github.com/vmihailenco/msgpack"
)

// DataFormat defines different data type from PSI
type DataFormat interface {
	Deploy()
	Serialize()
	Restore()
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
	Sites            map[string]uint8
}

// NewUniModOntology constructucst a set of UniMod ontologies
func NewUniModOntology() Onto {

	var m met.Data
	var o Onto

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
	o.Deploy()

	// Read
	o.Parse()

	// Serielize
	o.Serialize()

	return o
}

// Deploy deploys the OBO file to the temp folder
func (m *Onto) Deploy() {

	m.OboFile = fmt.Sprintf("%s%sunimod.obo", m.Temp, string(filepath.Separator))

	unmd.Deploy(m.OboFile)

	return
}

// Parse reads the unimod.obo file and creates the data structure
func (m *Onto) Parse() {

	file, e := os.Open(m.OboFile)
	if e != nil {
		msg.ReadFile(e, "fatal")
	}
	defer file.Close()

	var flag = -1
	var term Term
	term.Sites = make(map[string]uint8)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		if strings.HasPrefix(scanner.Text(), "format-version:") {
			m.Version = splitAndCollect(scanner.Text(), "common")
		}

		if strings.HasPrefix(scanner.Text(), "date:") {
			m.Date = splitAndCollect(scanner.Text(), "common")
		}

		if strings.Contains(scanner.Text(), "[Term]") {
			flag = 1
		}

		if strings.HasPrefix(scanner.Text(), "id:") && flag == 1 {
			term.ID = splitAndCollect(scanner.Text(), "common")
		} else if strings.HasPrefix(scanner.Text(), "name:") && flag == 1 {
			term.Name = splitAndCollect(scanner.Text(), "common")
		} else if strings.HasPrefix(scanner.Text(), "def:") && flag == 1 {
			term.Definition = splitAndCollect(scanner.Text(), "common")
		} else if strings.HasPrefix(scanner.Text(), "comment:") && flag == 1 {
			term.Comments = splitAndCollect(scanner.Text(), "common")
		} else if strings.HasPrefix(scanner.Text(), "synonym:") && flag == 1 {
			term.Synonyms = splitAndCollect(scanner.Text(), "common")
		}

		if strings.HasPrefix(scanner.Text(), "comment") && flag == 1 {
			term.Comments = splitAndCollect(scanner.Text(), "common")
		}

		if strings.HasPrefix(scanner.Text(), "synonym") && flag == 1 {
			term.Synonyms = splitAndCollect(scanner.Text(), "common")
		}

		if strings.HasPrefix(scanner.Text(), "xref: record_id") && flag == 1 {
			i, _ := strconv.Atoi(splitAndCollect(scanner.Text(), "xref"))
			term.RecordID = i
		} else if strings.HasPrefix(scanner.Text(), "xref: delta_mono_mass") && flag == 1 {
			i, _ := strconv.ParseFloat(splitAndCollect(scanner.Text(), "xref"), 64)
			term.MonoIsotopicMass = i
		} else if strings.HasPrefix(scanner.Text(), "xref: delta_avge_mass") && flag == 1 {
			i, _ := strconv.ParseFloat(splitAndCollect(scanner.Text(), "xref"), 64)
			term.AverageMass = i
		} else if strings.HasPrefix(scanner.Text(), "xref: delta_composition") && flag == 1 {
			term.Composition = splitAndCollect(scanner.Text(), "xref")
		} else if strings.HasPrefix(scanner.Text(), "xref: date_time_posted") && flag == 1 {
			term.DateTimePosted = splitAndCollect(scanner.Text(), "xref")
		} else if strings.HasPrefix(scanner.Text(), "xref: date_time_modified") && flag == 1 {
			term.DateTimeModified = splitAndCollect(scanner.Text(), "xref")
		}

		if strings.HasPrefix(scanner.Text(), "is_a") && flag == 1 {
			term.IsA = splitAndCollect(scanner.Text(), "common")
		}

		if strings.EqualFold(scanner.Text(), "//") && flag == 1 {
			flag = 0
			m.Terms = append(m.Terms, term)
			term = Term{}
			term.Sites = make(map[string]uint8)
		}

		if strings.Contains(scanner.Text(), "xref:") && strings.Contains(scanner.Text(), "_site") && flag == 1 {
			term.Sites[splitAndCollect(scanner.Text(), "site")]++
		}

	}

	if e := scanner.Err(); e != nil {
		log.Fatal(e)
	}

	return
}

// Serialize UniMod data structure
func (m Onto) Serialize() {

	b, e := msgpack.Marshal(&m)
	if e != nil {
		msg.MarshalFile(e, "fatal")
	}

	e = ioutil.WriteFile(sys.MODBin(), b, sys.FilePermission())
	if e != nil {
		msg.WriteFile(e, "fatal")
	}

	return
}

// Restore reads philosopher results files and restore the data sctructure
func (m *Onto) Restore() {

	b, e := ioutil.ReadFile(sys.MODBin())
	if e != nil {
		msg.ReadFile(e, "fatal")
	}

	e = msgpack.Unmarshal(b, &m)
	if e != nil {
		msg.DecodeMsgPck(e, "fatal")
	}

	return
}

func splitAndCollect(s string, target string) string {

	if target == "common" {
		l := strings.Split(s, ": ")
		return strings.Replace(l[1], "\"", "", -1)
	} else if target == "xref" {
		l := strings.Split(s, "\"")
		return strings.Replace(l[1], "\"", "", -1)
	} else if target == "site" {
		l := strings.Split(s, "_site ")
		return strings.Replace(l[1], "\"", "", -1)
	}

	return ""
}
