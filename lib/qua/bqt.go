package qua

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"philosopher/lib/msg"

	"philosopher/lib/dat"
	"philosopher/lib/ext/cdhit"
	"philosopher/lib/rep"
	"philosopher/lib/sys"

	"github.com/sirupsen/logrus"
)

// Cluster struct
type Cluster struct {
	Centroid                string
	Description             string
	Status                  string
	Existence               string
	GeneNames               string
	Number                  int
	TotalPeptideNumber      int
	SharedPeptides          int
	Coverage                float32
	UniqueClusterTopPepProb float64
	TopPepProb              float64
	Peptides                []string
	PeptideIons             []string
	UniqueClusterPeptides   []string
	Members                 map[string]uint8
}

// List list
type List []Cluster

// CreateCleanDataBaseReference removes decoys and contaminants from struct
func createCleanDataBaseReference(uid, temp string) string {

	var u dat.Base
	u.Restore()

	clstrFasta := temp + string(filepath.Separator) + uid + ".fasta"

	file, e := os.Create(clstrFasta)
	if e != nil {
		msg.Custom(errors.New("Could not create output for binning"), "fatal")
	}
	defer file.Close()

	for _, i := range u.Records {

		if i.IsDecoy == false && i.IsContaminant == false {

			line := fmt.Sprintf(">%s\n%s\n", i.OriginalHeader, i.Sequence)
			_, e = io.WriteString(file, line)

			if e != nil {
				msg.Custom(errors.New("Could not create output for binning"), "fatal")
			}

		}
	}

	return clstrFasta
}

// Execute is top function for Comet
func execute(level float64) (string, string) {

	cd := cdhit.New()

	cd.ClusterFasta = cd.FileName + ".fasta"
	cd.ClusterFile = cd.ClusterFasta + ".clstr"

	// deploy binary and parameter to workdir
	cd.Deploy()

	// run cdhit and create the clusters
	cd.Run(level)

	return cd.ClusterFile, cd.ClusterFasta
}

// ParseClusterFile ...
func parseClusterFile(cls, database string) List {

	var list List
	var clustermap = make(map[int][]string)
	var centroidmap = make(map[int]string)
	var clusterNumber int
	var seqsName []string

	f, e := os.Open(cls)
	if e != nil {
		msg.Custom(errors.New("Cannot open cluster file"), "fatal")
	}
	defer f.Close()

	reheader, e1 := regexp.Compile(`^>Cluster\s+(.*)`)
	if e1 != nil {
		msg.Custom(errors.New("Cannot compile Cluster header regex"), "fatal")
	}

	reseq, e2 := regexp.Compile(`\|(.*)\|.*`)
	if e2 != nil {
		msg.Custom(errors.New("Cannot compile Cluster description regex"), "fatal")
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {

		if strings.HasPrefix(scanner.Text(), ">") {

			cluster := reheader.FindStringSubmatch(scanner.Text())
			num := cluster[1]
			i, e := strconv.Atoi(num)
			if e != nil {
				msg.Custom(errors.New("FAST header not found"), "fatal")
			}
			clusterNumber = i

			clustermap[clusterNumber] = append(clustermap[clusterNumber], "")
			centroidmap[clusterNumber] = ""

		} else {

			if strings.Contains(scanner.Text(), "*") {
				centroid := strings.Split(scanner.Text(), "|")
				//centroid := reseq.FindStringSubmatch(scanner.Text())
				if len(centroid) < 2 {
					msg.Custom(errors.New("FASTA file contains non-formatted sequence headers"), "fatal")
				}
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

	return list
}

// MapProtXML2Clusters ...
func mapProtXML2Clusters(clusters List) List {

	var e rep.Evidence
	e.RestoreGranular()

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

func retrieveInfoFromUniProtDB(clusters List) List {

	// collect database information
	var dtb dat.Base
	dtb.Restore()

	for i := range clusters {
		for _, j := range dtb.Records {
			if strings.EqualFold(clusters[i].Centroid, j.ID) && j.IsDecoy == false && j.IsContaminant == false {
				clusters[i].Description = j.ProteinName
				clusters[i].GeneNames = j.GeneNames
				break
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
		query = fmt.Sprintf("%s%s%s", "http://www.uniprot.org/uniprot/?query=organism:", organism, "&columns=id,reviewed,existence,genes,feature(DOMAIN%20EXTENT),comment(PATHWAY),go-id&format=tab")
	}

	outfile := fmt.Sprintf("%s/%s.tab", resultDir, organism)

	// tries to create an output file
	output, e := os.Create(outfile)
	if e != nil {
		msg.WriteFile(e, "fatal")
	}
	defer output.Close()

	// Tries to query data from Uniprot
	response, e := http.Get(query)
	if e != nil {
		msg.Custom(errors.New("Could not find the annotation file"), "fatal")
	}
	defer response.Body.Close()

	// Tries to download data from Uniprot
	_, e = io.Copy(output, response.Body)
	if e != nil {
		msg.Custom(errors.New("Cannot download the annotation file"), "fatal")
	}

	faMap = make(map[string][]string)

	f, e := os.Open(outfile)
	if outfile == "" || e != nil {
		msg.Custom(errors.New("Emty or inexisting file"), "fatal")
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

// SavetoDisk saves functional inference result to disk
func savetoDisk(list List, temp, uid string) {

	output := fmt.Sprintf("%s%sclusters.tsv", temp, string(filepath.Separator))

	// create result file
	file, e := os.Create(output)
	if e != nil {
		msg.WriteFile(e, "fatal")
	}
	defer file.Close()

	var line string
	line = fmt.Sprintf("Cluster Number\tRepresentative\tTotal Members\tMembers\tPercentage Coverage\tTotal Peptides\tIntra Cluster Peptides\tInter Cluster Peptides\tDescription\n")

	if len(uid) > 0 {
		logrus.Info("Retrieving annotation from UniProt")
		line = fmt.Sprintf("Cluster Number\tRepresentative\tTotal Members\tMembers\tPercentage Coverage\tTotal Peptides\tIntra Cluster Peptides\tInter Cluster Peptides\tDescription\tStatus\tExistence\tGenes\tProtein Domains\tPathways\tGene Ontology\n")
	}

	_, e = io.WriteString(file, line)
	if e != nil {
		msg.WriteToFile(e, "fatal")
	}

	var faMap = make(map[string][]string)
	if len(uid) > 0 {
		faMap = getFile(true, temp, uid)
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
				(list[i].TotalPeptideNumber - len(list[i].UniqueClusterPeptides)),
				list[i].Description)

			v, ok := faMap[list[i].Centroid]
			if ok {
				var index int
				if len(uid) > 0 {
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

			_, e := io.WriteString(file, line)
			if e != nil {
				msg.WriteToFile(e, "fatal")
			}
		}

	}

	sys.CopyFile(output, filepath.Base(output))

	return
}
