package clu

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/ext/cdhit"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/rep"
	"github.com/prvst/philosopher/lib/sys"
)

// Cluster struct
type Cluster struct {
	met.Data
	UID                     string
	Level                   float64
	Number                  int
	Centroid                string
	Description             string
	Coverage                float32
	Members                 map[string]uint8
	GeneNames               string
	Peptides                []string
	PeptideIons             []string
	TotalPeptideNumber      int
	SharedPeptides          int
	UniqueClusterPeptides   []string
	UniqueClusterTopPepProb float64
	TopPepProb              float64
}

// List list
type List []Cluster

// New constructor
func New() Cluster {

	var o Cluster
	var m met.Data
	m.Restore(sys.Meta())

	o.UUID = m.UUID
	o.Distro = m.Distro
	o.Home = m.Home
	o.MetaFile = m.MetaFile
	o.MetaDir = m.MetaDir
	o.DB = m.DB
	o.Temp = m.Temp
	o.TimeStamp = m.TimeStamp
	o.OS = m.OS
	o.Arch = m.Arch

	return o
}

// GenerateReport creates the protein report output
func (c *Cluster) GenerateReport() error {

	var err error

	// create clean reference db for clustering
	clusterFasta, err := c.CreateCleanDataBaseReference()
	if err != nil {
		return err
	}

	// run cdhit, create cluster file
	logrus.Info("Clustering")
	clusterFile, clusterFasta := c.Execute()

	// parse the cluster file
	logrus.Info("Parsing clusters")
	clusters, err := c.ParseClusterFile(clusterFile, clusterFasta)

	// maps all proteins from the db against the clusters
	logrus.Info("Mapping proteins to clusters")
	mappedClust := mapProtXML2Clusters(clusters)

	// mapping to functional annotation and save to disk
	c.SavetoDisk(mappedClust)

	if err != nil {
		return err
	}

	return nil
}

// CreateCleanDataBaseReference removes decoys and contaminants from struct
func (c *Cluster) CreateCleanDataBaseReference() (string, error) {

	var err error

	var u dat.Base
	u.Restore()

	clstrFasta := c.Temp + string(filepath.Separator) + c.UUID + ".fasta"

	file, err := os.Create(clstrFasta)
	if err != nil {
		msg := "Could not create output for binning"
		return "", errors.New(msg)
	}
	defer file.Close()

	for _, i := range u.Records {

		if i.IsDecoy == false && i.IsContaminant == false {

			line := fmt.Sprintf(">%s\n%s\n", i.OriginalHeader, i.Sequence)
			_, err = io.WriteString(file, line)

			if err != nil {
				msg := "Could not create output for binning"
				return "", errors.New(msg)
			}

		}
	}

	return clstrFasta, nil
}

// Execute is top function for Comet
func (c *Cluster) Execute() (string, string) {

	cd := cdhit.New()

	cd.ClusterFasta = cd.FileName + ".fasta"
	cd.ClusterFile = cd.ClusterFasta + ".clstr"

	// deploy binary and paramter to workdir
	cd.Deploy()

	// run cdhit and create the clusters
	cd.Run(c.Level)

	return cd.ClusterFile, cd.ClusterFasta
}

// ParseClusterFile ...
func (c *Cluster) ParseClusterFile(cls, database string) (List, error) {

	var list List
	var clustermap = make(map[int][]string)
	var centroidmap = make(map[int]string)
	var clusterNumber int
	var seqsName []string
	var err error

	f, err := os.Open(cls)
	if err != nil {
		msg := "[ERROR] Cannot open cluster file" + cls
		return nil, errors.New(msg)
	}
	defer f.Close()

	reheader, err := regexp.Compile(`^>Cluster\s+(.*)`)
	reseq, err := regexp.Compile(`\|(.*)\|.*`)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {

		if strings.HasPrefix(scanner.Text(), ">") {

			cluster := reheader.FindStringSubmatch(scanner.Text())
			num := cluster[1]
			i, ferr := strconv.Atoi(num)
			if ferr != nil {
				return nil, errors.New("FASTA header not found")
			}
			clusterNumber = i

			clustermap[clusterNumber] = append(clustermap[clusterNumber], "")
			centroidmap[clusterNumber] = ""

		} else {

			if strings.Contains(scanner.Text(), "*") {
				centroid := strings.Split(scanner.Text(), "|")
				//centroid := reseq.FindStringSubmatch(scanner.Text())
				centroidmap[clusterNumber] = centroid[1]
			}

			seq := reseq.FindStringSubmatch(scanner.Text())
			seqsName = clustermap[clusterNumber]
			seqsName = append(seqsName, seq[1])
			clustermap[clusterNumber] = seqsName
		}
	}

	var u dat.Base
	u.Restore()

	var fastaMap = make(map[string]string)
	for _, i := range u.Records {
		fastaMap[i.ID] = i.ProteinName
	}

	for i := 0; i < len(clustermap); i++ {
		var memberMap = make(map[string]uint8)
		arr := clustermap[i][1:]
		for j := range arr {
			memberMap[arr[j]] = 0
		}
		c := Cluster{Number: i, Centroid: centroidmap[i], Description: fastaMap[centroidmap[i]], Members: memberMap}
		list = append(list, c)
	}

	if err != nil {
		return nil, err
	}

	return list, nil
}

// SavetoDisk saves functional inference result to disk
func (c *Cluster) SavetoDisk(list List) {

	output := fmt.Sprintf(".%sclusters.tsv", string(filepath.Separator))

	// create result file
	file, err := os.Create(output)
	if err != nil {
		log.Println("[ERROR]:", err)
		os.Exit(2)
	}
	defer file.Close()

	var line string
	line = fmt.Sprintf("Cluster Number\tRepresentative\tTotal Members\tMembers\tPercentage Coverage\tTotal Peptides\tIntra Cluster Peptides\tInter Cluster Peptides\tDescription\n")

	if len(c.UID) > 0 {
		log.Println("Retrieving annotation from UniProt")
		line = fmt.Sprintf("Cluster Number\tRepresentative\tTotal Members\tMembers\tPercentage Coverage\tTotal Peptides\tIntra Cluster Peptides\tInter Cluster Peptides\tDescription\tStatus\tExistence\tGenes\tProtein Domains\tPathways\tGene Ontology\n")
	}

	n, err := io.WriteString(file, line)
	if err != nil {
		log.Println("[ERROR]", n, err)
		os.Exit(2)
	}

	var faMap = make(map[string][]string)
	if len(c.UID) > 0 {
		faMap = getFile(true, c.Home, c.UID)
	} else {
		//faMap, _ = fasta.ParseFastaDescription(rep.DB)
	}

	for i := range list {

		if list[i].TotalPeptideNumber > 0 {

			var members []string
			for k := range list[i].Members {
				members = append(members, k)
			}
			membersString := strings.Join(members, ", ")

			line := fmt.Sprintf("%d\t%s\t%d\t%s\t%.2f\t%d\t%d\t%d\t%s\t",
				list[i].Number,
				list[i].Centroid,
				len(list[i].Members),
				membersString,
				list[i].Coverage,
				list[i].TotalPeptideNumber,
				len(list[i].UniqueClusterPeptides),
				//list[i].UniqueClusterTopPepProb,
				(list[i].TotalPeptideNumber - len(list[i].UniqueClusterPeptides)),
				list[i].Description)

			v, ok := faMap[list[i].Centroid]
			if ok {
				var index int
				if len(c.UID) > 0 {
					index = 1
				} else {
					index = 0
				}
				for i := index; i < len(v); i++ {
					item := v[i] + "\t"
					line += item
				}
			}

			line += "\n"

			n, err := io.WriteString(file, line)
			if err != nil {
				logrus.Println("[ERROR]", n, err)
				os.Exit(2)
			}
		}

	}

	return
}

// MapProtXML2Clusters ...
func mapProtXML2Clusters(clusters List) List {

	var e rep.Evidence
	e.Restore()

	for _, i := range e.Proteins {
		if i.IsDecoy == false && i.IsContaminant == false {
			for j := range clusters {

				_, ok := clusters[j].Members[i.ProteinID]
				if ok {

					clusters[j].Members[i.ProteinID]++
					clusters[j].TotalPeptideNumber += len(i.TotalPeptideIons)

					if i.Coverage > clusters[j].Coverage {
						clusters[j].Coverage = i.Coverage
					}

					for _, k := range i.TotalPeptideIons {
						clusters[j].Peptides = append(clusters[j].Peptides, k.Sequence)
					}

					for _, k := range i.TotalPeptideIons {
						if clusters[j].TopPepProb < k.Probability {
							clusters[j].TopPepProb = k.Probability
						}
					}

				}

			}
		}

	}

	// creates a global peptide map
	pepMap := make(map[string]uint8)
	for _, i := range e.Proteins {
		for _, j := range i.TotalPeptideIons {
			_, ok := pepMap[j.Sequence]
			if ok {
				pepMap[j.Sequence]++
			} else {
				pepMap[j.Sequence] = 1
			}
		}
	}

	// for i := range clusters {
	//   for j := range clusters[i].Peptides {
	//
	//     _, ok := pepMap[clusters[i].Peptides[j]]
	//     if ok {
	//       pepMap[clusters[i].Peptides[j]]++
	//     } else {
	//       pepMap[clusters[i].Peptides[j]] = 1
	//     }
	//   }
	// }

	// now runs for each cluster and checks if the peptides appear in other clusters
	for i := range clusters {
		for j := range clusters[i].Peptides {
			v, ok := pepMap[clusters[i].Peptides[j]]
			if ok {
				if v > 1 {
					clusters[i].SharedPeptides++
					clusters[i].UniqueClusterTopPepProb = clusters[i].TopPepProb
				} else {
					clusters[i].UniqueClusterPeptides = append(clusters[i].UniqueClusterPeptides, clusters[i].Peptides[j])

					if clusters[i].UniqueClusterTopPepProb < clusters[i].TopPepProb {
						clusters[i].UniqueClusterTopPepProb = clusters[i].TopPepProb
					}

				}
			}
		}
	}

	return clusters
}

func retrieveInfoFromUniProtDB(clusters List, db dat.Base) List {

	for i := range clusters {
		for _, j := range db.Records {

			if strings.EqualFold(clusters[i].Centroid, j.ID) && j.IsDecoy == false && j.IsContaminant == false {
				clusters[i].Description = j.ProteinName
				clusters[i].GeneNames = j.GeneNames
			}

		}
	}

	return clusters
}

// GetFile is the miun function from annot package. It's responsible for connecting Uniprot
// using ans Organism ID and retrieving functional information.
func getFile(getAll bool, resultDir string, organism string) (faMap map[string][]string) {

	var query string
	query = fmt.Sprintf("%s%s%s", "http://www.uniprot.org/uniprot/?query=organism:", organism, "&columns=id,protein%20names&format=tab")

	if getAll == true {
		query = fmt.Sprintf("%s%s%s", "http://www.uniprot.org/uniprot/?query=organism:", organism, "&columns=id,reviewed,existence,genes,domains,pathway,go-id&format=tab")
	}

	outfile := fmt.Sprintf("%s/%s.tab", resultDir, organism)

	// tries to create an output file
	output, err := os.Create(outfile)
	if err != nil {
		log.Println("[ERROR] Could not create output file", query, "-", err)
		os.Exit(2)
	}
	defer output.Close()

	// Tries to query data from Uniprot
	response, err := http.Get(query)
	if err != nil {
		log.Println("[ERROR] Could not find annotation file", err)
		os.Exit(2)
	}
	defer response.Body.Close()

	// Tries to download data from Uniprot
	n, err := io.Copy(output, response.Body)
	if err != nil {
		log.Println("[ERROR] Could not download annotation file", n, "-", err)
		return
	}

	faMap = make(map[string][]string)

	f, err := os.Open(outfile)
	if outfile == "" || err != nil {
		log.Println("[ERROR] Empty or inexistent file")
		os.Exit(2)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), "\t")
		faMap[arr[0]] = arr
	}

	return
}

func parseFastaFile(db dat.Base) map[string]string {

	var fastaMap = make(map[string]string)

	// get protein ID and description and add them to fastaMap
	for _, k := range db.Records {
		fastaMap[k.ID] = k.ProteinName
	}

	return fastaMap
}
