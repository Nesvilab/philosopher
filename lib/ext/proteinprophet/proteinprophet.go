package proteinprophet

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/prvst/philosopher/lib/err"
	unix "github.com/prvst/philosopher/lib/ext/proteinprophet/unix"
	wPoP "github.com/prvst/philosopher/lib/ext/proteinprophet/win"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
)

// ProteinProphet is tool configuration
type ProteinProphet struct {
	DefaultBatchCoverage  string
	DefaultDatabaseParser string
	DefaultProteinProphet string
	WinBatchCoverage      string
	WinDatabaseParser     string
	WinProteinProphet     string
	UnixBatchCoverage     string
	UnixDatabaseParser    string
	UnixProteinProphet    string
	Prot2html             string
	LibgccDLL             string
	Zlib1DLL              string
}

// New constructor
func New() ProteinProphet {

	var self ProteinProphet

	temp, _ := sys.GetTemp()

	self.UnixBatchCoverage = temp + string(filepath.Separator) + "batchcoverage"
	self.UnixDatabaseParser = temp + string(filepath.Separator) + "DatabaseParser"
	self.UnixProteinProphet = temp + string(filepath.Separator) + "ProteinProphet"
	self.WinBatchCoverage = temp + string(filepath.Separator) + "batchcoverage.exe"
	self.WinDatabaseParser = temp + string(filepath.Separator) + "DatabaseParser.exe"
	self.WinProteinProphet = temp + string(filepath.Separator) + "ProteinProphet.exe"
	self.LibgccDLL = temp + string(filepath.Separator) + "libgcc_s_dw2-1.dll"
	self.Zlib1DLL = temp + string(filepath.Separator) + "zlib1.dll"

	return self
}

// Run is the main entry point for ProteinProphet
func Run(m met.Data, args []string) met.Data {

	var pop = New()

	if len(args) < 1 {
		logrus.Fatal("No input file provided")
	}

	// deploy the binaries
	e := pop.Deploy(m.OS, m.Distro)
	if e != nil {
		logrus.Fatal(e.Message)
	}

	// run ProteinProphet
	xml, e := pop.Execute(m.ProteinProphet, m.Home, m.Temp, args)
	if e != nil {
		logrus.Fatal(e.Message)
	}
	_ = xml

	m.ProteinProphet.InputFiles = args

	return m
}

// Deploy generates comet binary on workdir bin directory
func (p *ProteinProphet) Deploy(os, distro string) *err.Error {

	if os == sys.Windows() {
		wPoP.WinBatchCoverage(p.WinBatchCoverage)
		p.DefaultBatchCoverage = p.WinBatchCoverage
		wPoP.WinDatabaseParser(p.WinDatabaseParser)
		p.DefaultDatabaseParser = p.WinDatabaseParser
		wPoP.WinProteinProphet(p.WinProteinProphet)
		p.DefaultProteinProphet = p.WinProteinProphet
		wPoP.LibgccDLL(p.LibgccDLL)
		wPoP.Zlib1DLL(p.Zlib1DLL)
	} else {
		if strings.EqualFold(distro, sys.Debian()) {
			unix.UnixBatchCoverage(p.UnixBatchCoverage)
			p.DefaultBatchCoverage = p.UnixBatchCoverage
			unix.UnixDatabaseParser(p.UnixDatabaseParser)
			p.DefaultDatabaseParser = p.UnixDatabaseParser
			unix.UnixProteinProphet(p.UnixProteinProphet)
			p.DefaultProteinProphet = p.UnixProteinProphet
		} else if strings.EqualFold(distro, sys.Redhat()) {
			unix.UnixBatchCoverage(p.UnixBatchCoverage)
			p.DefaultBatchCoverage = p.UnixBatchCoverage
			unix.UnixDatabaseParser(p.UnixDatabaseParser)
			p.DefaultDatabaseParser = p.UnixDatabaseParser
			unix.UnixProteinProphet(p.UnixProteinProphet)
			p.DefaultProteinProphet = p.UnixProteinProphet
		} else {
			return &err.Error{Type: err.UnsupportedDistribution, Class: err.FATA, Argument: "dont know how to deploy ProteinProphet"}
		}
	}

	return nil
}

// Execute ProteinProphet executes peptideprophet
func (p ProteinProphet) Execute(params met.ProteinProphet, home, temp string, args []string) ([]string, *err.Error) {

	// run
	bin := p.DefaultProteinProphet
	cmd := exec.Command(bin)

	// append pepxml files
	for i := range args {
		file, _ := filepath.Abs(args[i])
		cmd.Args = append(cmd.Args, file)
	}

	// append output file
	output := fmt.Sprintf("%s%s%s.prot.xml", temp, string(filepath.Separator), params.Output)
	output, _ = filepath.Abs(output)

	cmd.Args = append(cmd.Args, output)
	cmd = p.appendParams(params, cmd)

	cmd.Dir = filepath.Dir(output)

	env := os.Environ()
	env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
	env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", temp))
	for i := range env {
		if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
			env[i] = env[i] + ";" + temp
		}
	}
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Start()
	if e != nil {
		return nil, &err.Error{Type: err.CannotExecuteBinary, Class: err.FATA, Argument: "ProteinProphet"}
	}
	_ = cmd.Wait()

	// copy to work directory
	dest := fmt.Sprintf("%s%s%s", home, string(filepath.Separator), filepath.Base(output))
	e = sys.CopyFile(output, dest)
	if e != nil {
		return nil, &err.Error{Type: err.CannotCopyFile, Class: err.FATA, Argument: "ProtXML"}
	}

	// collect all resulting files
	var processedOutput []string
	for _, i := range cmd.Args {
		if strings.Contains(i, "prot.xml") || i == params.Output {
			processedOutput = append(processedOutput, i)
		}
	}

	return processedOutput, nil
}

func (p ProteinProphet) appendParams(params met.ProteinProphet, cmd *exec.Cmd) *exec.Cmd {

	if params.ExcludeZ == true {
		cmd.Args = append(cmd.Args, "EXCLUDE_ZEROS")
	}

	if params.Noplot == true {
		cmd.Args = append(cmd.Args, "NOPLOT")
	}

	if params.Nooccam == true {
		cmd.Args = append(cmd.Args, "NOOCCAM")
	}

	if params.Softoccam == true {
		cmd.Args = append(cmd.Args, "SOFTOCCAM")
	}

	if params.Icat == true {
		cmd.Args = append(cmd.Args, "ICAT")
	}

	if params.Glyc == true {
		cmd.Args = append(cmd.Args, "GLYC")
	}

	if params.Nogroupwts == true {
		cmd.Args = append(cmd.Args, "NOGROUPWTS")
	}

	if params.NonSP == true {
		cmd.Args = append(cmd.Args, "NONSP")
	}

	if params.Accuracy == true {
		cmd.Args = append(cmd.Args, "ACCURACY")
	}

	if params.Asap == true {
		cmd.Args = append(cmd.Args, "ASAP")
	}

	if params.Refresh == true {
		cmd.Args = append(cmd.Args, "REFRESH")
	}

	if params.Normprotlen == true {
		cmd.Args = append(cmd.Args, "NORMPROTLEN")
	}

	if params.Logprobs == true {
		cmd.Args = append(cmd.Args, "LOGPROBS")
	}

	if params.Confem == true {
		cmd.Args = append(cmd.Args, "CONFEM")
	}

	if params.Allpeps == true {
		cmd.Args = append(cmd.Args, "ALLPEPS")
	}

	if params.Unmapped == true {
		cmd.Args = append(cmd.Args, "UNMAPPED")
	}

	if params.Noprotlen == true {
		cmd.Args = append(cmd.Args, "NOPROTLEN")
	}

	if params.Instances == true {
		cmd.Args = append(cmd.Args, "INSTANCES")
	}

	if params.Fpkm == true {
		cmd.Args = append(cmd.Args, "FPKM")
	}

	if params.Protmw == true {
		cmd.Args = append(cmd.Args, "PROTMW")
	}

	if params.Iprophet == true {
		cmd.Args = append(cmd.Args, "IPROPHET")
	}

	if params.Asapprophet == true {
		cmd.Args = append(cmd.Args, "ASAP_PROPHET")
	}

	if params.Delude == true {
		cmd.Args = append(cmd.Args, "DELUDE")
	}

	// // there is an error in the way how the modified version was implemented.
	// // The mod version is *always* active, and the tag makes it normal again.
	// // it should be the oposite, so thats why this block looks like that.
	// if c.Excludemods == true {
	// 	// the program is always trying to process os'es
	// 	//cmd.Args = append(cmd.Args, "ALLOWDIFFPROBS")
	// } else {
	// 	// the tag makes the program running in "normal" mode
	// 	cmd.Args = append(cmd.Args, "ALLOWDIFFPROBS")
	// }

	if params.Maxppmdiff != 20 {
		v := fmt.Sprintf("MAXPPMDIFF%d", params.Maxppmdiff)
		cmd.Args = append(cmd.Args, v)
	}

	if params.Minprob != 0.05 {
		v := fmt.Sprintf("MINPROB%.4f", params.Minprob)
		cmd.Args = append(cmd.Args, v)
	}

	if params.Minindep != 0 {
		v := fmt.Sprintf("MININDEP%d", params.Minindep)
		cmd.Args = append(cmd.Args, v)
	}

	if params.Mufactor != 1 {
		v := fmt.Sprintf("MUFACTOR%d", params.Mufactor)
		cmd.Args = append(cmd.Args, v)
	}

	return cmd
}
