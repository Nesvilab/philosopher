package msfragger

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/met"
	"github.com/sirupsen/logrus"
)

// MSFragger represents the tool configuration
type MSFragger struct {
	DefaultBin   string
	DefaultParam string
}

// New constructor
func New(temp string) MSFragger {

	var self MSFragger

	self.DefaultBin = ""
	self.DefaultParam = ""

	return self
}

// Run is the Fragger main entry point
func Run(m met.Data, args []string) (met.Data, *err.Error) {

	var frg = New(m.Temp)

	// if len(m.MSFragger.Param) < 1 {
	// 	return m, &err.Error{Type: err.CannotRunMSFragger, Class: err.WARN, Argument: "No parameter file found, using values defined via command line"}
	// }

	// collect and store the mz files
	m.MSFragger.RawFiles = args

	if len(m.MSFragger.Param) > 1 {
		// convert the param file to binary and store it in meta
		var binFile []byte
		paramAbs, _ := filepath.Abs(m.MSFragger.Param)
		binFile, e := ioutil.ReadFile(paramAbs)
		if e != nil {
			logrus.Fatal(e)
		}
		m.MSFragger.ParamFile = binFile
	}

	// run comet
	e := frg.Execute(args, m.MSFragger)
	if e != nil {
		//logrus.Fatal(e)
	}

	return m, nil
}

// Execute is the main fucntion to execute MSFragger
func (c *MSFragger) Execute(cmdArgs []string, m met.MSFragger) *err.Error {

	mem := fmt.Sprintf("-Xmx%sG", m.Memmory)
	cmd := exec.Command("java", "-jar", mem, m.JarPath, m.Param)

	for _, i := range cmdArgs {
		file, _ := filepath.Abs(i)
		cmd.Args = append(cmd.Args, file)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Start()
	if e != nil {
		return nil
	}

	_ = cmd.Wait()

	return nil
}

func (c *MSFragger) appendParams(params met.MSFragger, cmd *exec.Cmd) *exec.Cmd {

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
