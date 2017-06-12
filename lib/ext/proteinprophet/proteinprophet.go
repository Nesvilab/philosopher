package proteinprophet

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	unix "github.com/prvst/philosopher-source/lib/ext/proteinprophet/unix"
	wPoP "github.com/prvst/philosopher-source/lib/ext/proteinprophet/win"
	"github.com/prvst/philosopher-source/lib/meta"
	"github.com/prvst/philosopher-source/lib/sys"
)

// ProteinProphet is tool configuration
type ProteinProphet struct {
	meta.Data
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
	Minprob               string
	Minindep              string
	Mufactor              string
	Output                string
	Maxppmdiff            string
	//Combine               bool
	Noplot      bool
	Nooccam     bool
	Softoccam   bool
	Icat        bool
	Glyc        bool
	Nogroupwts  bool
	NonSP       bool
	Accuracy    bool
	Asap        bool
	Refresh     bool
	Normprotlen bool
	Logprobs    bool
	Confem      bool
	Allpeps     bool
	Unmapped    bool
	Noprotlen   bool
	Instances   bool
	Fpkm        bool
	Protmw      bool
	Iprophet    bool
	Asapprophet bool
	Delude      bool
	Excludemods bool
}

// New constructor
func New() ProteinProphet {

	var o ProteinProphet
	var m meta.Data
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

	o.UnixBatchCoverage = o.Temp + string(filepath.Separator) + "batchcoverage"
	o.UnixDatabaseParser = o.Temp + string(filepath.Separator) + "DatabaseParser"
	o.UnixProteinProphet = o.Temp + string(filepath.Separator) + "ProteinProphet"
	o.WinBatchCoverage = o.Temp + string(filepath.Separator) + "batchcoverage.exe"
	o.WinDatabaseParser = o.Temp + string(filepath.Separator) + "DatabaseParser.exe"
	o.WinProteinProphet = o.Temp + string(filepath.Separator) + "ProteinProphet.exe"
	o.LibgccDLL = o.Temp + string(filepath.Separator) + "libgcc_s_dw2-1.dll"
	o.Zlib1DLL = o.Temp + string(filepath.Separator) + "zlib1.dll"

	return o
}

// Deploy generates comet binary on workdir bin directory
func (c *ProteinProphet) Deploy() error {

	if c.OS == sys.Windows() {
		wPoP.WinBatchCoverage(c.WinBatchCoverage)
		c.DefaultBatchCoverage = c.WinBatchCoverage
		wPoP.WinDatabaseParser(c.WinDatabaseParser)
		c.DefaultDatabaseParser = c.WinDatabaseParser
		wPoP.WinProteinProphet(c.WinProteinProphet)
		c.DefaultProteinProphet = c.WinProteinProphet
		wPoP.LibgccDLL(c.LibgccDLL)
		wPoP.Zlib1DLL(c.Zlib1DLL)
	} else {
		if strings.EqualFold(c.Distro, sys.Debian()) {
			unix.UnixBatchCoverage(c.UnixBatchCoverage)
			c.DefaultBatchCoverage = c.UnixBatchCoverage
			unix.UnixDatabaseParser(c.UnixDatabaseParser)
			c.DefaultDatabaseParser = c.UnixDatabaseParser
			unix.UnixProteinProphet(c.UnixProteinProphet)
			c.DefaultProteinProphet = c.UnixProteinProphet
		} else if strings.EqualFold(c.Distro, sys.Redhat()) {
			unix.UnixBatchCoverage(c.UnixBatchCoverage)
			c.DefaultBatchCoverage = c.UnixBatchCoverage
			unix.UnixDatabaseParser(c.UnixDatabaseParser)
			c.DefaultDatabaseParser = c.UnixDatabaseParser
			unix.UnixProteinProphet(c.UnixProteinProphet)
			c.DefaultProteinProphet = c.UnixProteinProphet
		} else {
			return errors.New("Unsupported distribution for ProteinProphet")
		}
	}

	return nil
}

// Run ProteinProphet executes peptideprophet
func (c *ProteinProphet) Run(args []string) error {

	//if c.Combine == true {

	// run
	bin := c.DefaultProteinProphet
	cmd := exec.Command(bin)

	// append pepxml files
	for i := range args {
		file, _ := filepath.Abs(args[i])
		cmd.Args = append(cmd.Args, file)
	}

	// append output file
	output := fmt.Sprintf("%s%s%s.prot.xml", c.Temp, string(filepath.Separator), c.Output)
	output, _ = filepath.Abs(output)

	cmd.Args = append(cmd.Args, output)
	cmd = c.appendParams(cmd)

	cmd.Dir = filepath.Dir(output)

	env := os.Environ()
	env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
	env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", c.Temp))
	for i := range env {
		if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
			env[i] = env[i] + ";" + c.Temp
		}
	}
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		msg := fmt.Sprintf("Could not run ProteinProphet: %s", err)
		return errors.New(msg)
	}
	_ = cmd.Wait()

	var baseDir string
	baseDir = filepath.Dir(args[0])

	// copy to work directory
	dest := fmt.Sprintf("%s%s%s", baseDir, string(filepath.Separator), filepath.Base(output))
	err = sys.CopyFile(output, dest)
	if err != nil {
		return err
	}

	// } else {
	//
	// 	var files []string
	// 	for _, i := range args {
	// 		file, _ := filepath.Abs(i)
	// 		files = append(files, file)
	// 	}
	//
	// 	// append pepxml files
	// 	for _, i := range files {
	//
	// 		// run
	// 		bin := c.DefaultProteinProphet
	// 		cmd := exec.Command(bin)
	//
	// 		cmd.Args = append(cmd.Args, i)
	//
	// 		var name string
	// 		var base string
	// 		var baseDir string
	//
	// 		if strings.Contains(strings.ToLower(i), ".pep.xml") {
	// 			base = filepath.Base(i)
	// 			baseDir = filepath.Dir(i)
	// 			name = strings.TrimSuffix(base, ".pep.xml")
	// 		} else if strings.Contains(strings.ToLower(i), ".pepxml") {
	// 			base = filepath.Base(i)
	// 			baseDir = filepath.Dir(i)
	// 			name = strings.TrimSuffix(base, ".pepxml")
	// 		}
	//
	// 		fmt.Println(name)
	//
	// 		// append output file
	// 		output := fmt.Sprintf("%s%s%s.prot.xml", c.Temp, string(filepath.Separator), name)
	// 		output, _ = filepath.Abs(output)
	//
	// 		cmd.Args = append(cmd.Args, output)
	// 		cmd = c.appendParams(cmd)
	//
	// 		cmd.Dir = filepath.Dir(output)
	//
	// 		env := os.Environ()
	// 		env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
	// 		env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", c.Temp))
	// 		for i := range env {
	// 			if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
	// 				env[i] = env[i] + ";" + c.Temp
	// 			}
	// 		}
	// 		cmd.Env = env
	//
	// 		cmd.Stdout = os.Stdout
	// 		cmd.Stderr = os.Stderr
	// 		err := cmd.Start()
	// 		if err != nil {
	// 			msg := fmt.Sprintf("Could not run ProteinProphet: %s", err)
	// 			return errors.New(msg)
	// 		}
	// 		_ = cmd.Wait()
	//
	// 		// copy to work directory
	// 		dest := fmt.Sprintf("%s%s%s", baseDir, string(filepath.Separator), filepath.Base(output))
	// 		err = sys.CopyFile(output, dest)
	// 		if err != nil {
	// 			return err
	// 		}
	//
	// 	}
	// }
	//}

	return nil
}

// // Run ProteinProphet executes peptideprophet
// func (c *ProteinProphet) Run(args []string) error {
//
// 	if c.Combine == true {
//
// 		// run
// 		bin := c.DefaultProteinProphet
// 		cmd := exec.Command(bin)
//
// 		// append pepxml files
// 		for i := range args {
// 			file, _ := filepath.Abs(args[i])
// 			cmd.Args = append(cmd.Args, file)
// 		}
//
// 		// append output file
// 		output := fmt.Sprintf("%s%s%s", c.Temp, string(filepath.Separator), "interact.prot.xml")
// 		output, _ = filepath.Abs(output)
//
// 		cmd.Args = append(cmd.Args, output)
// 		cmd = c.appendParams(cmd)
//
// 		cmd.Dir = filepath.Dir(output)
//
// 		env := os.Environ()
// 		env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
// 		env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", c.Temp))
// 		for i := range env {
// 			if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
// 				env[i] = env[i] + ";" + c.Temp
// 			}
// 		}
// 		cmd.Env = env
//
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		err := cmd.Start()
// 		if err != nil {
// 			msg := fmt.Sprintf("Could not run ProteinProphet: %s", err)
// 			return errors.New(msg)
// 		}
// 		_ = cmd.Wait()
//
// 		// copy to work directory
// 		err = sys.CopyFile(output, filepath.Base(output))
// 		if err != nil {
// 			return err
// 		}
//
// 	} else {
//
// 		// append pepxml files
// 		for _, i := range args {
//
// 			bin := c.DefaultProteinProphet
// 			cmd := exec.Command(bin)
//
// 			file, _ := filepath.Abs(i)
// 			cmd.Args = append(cmd.Args, file)
//
// 			var name string
// 			var base string
// 			var baseDir string
// 			if strings.Contains(strings.ToLower(file), ".pep.xml") {
// 				base = filepath.Base(file)
// 				baseDir = filepath.Dir(file)
// 				name = strings.TrimSuffix(base, ".pep.xml")
// 			} else if strings.Contains(strings.ToLower(file), ".pepxml") {
// 				base = filepath.Base(file)
// 				baseDir = filepath.Dir(file)
// 				name = strings.TrimSuffix(base, ".pepxml")
// 			}
//
// 			// append output file
// 			output := fmt.Sprintf("%s%s%s.prot.xml", c.Temp, string(filepath.Separator), name)
// 			output, _ = filepath.Abs(output)
//
// 			cmd.Args = append(cmd.Args, output)
// 			cmd = c.appendParams(cmd)
//
// 			cmd.Dir = filepath.Dir(output)
//
// 			env := os.Environ()
// 			env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
// 			env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", c.Temp))
// 			for i := range env {
// 				if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
// 					env[i] = env[i] + ";" + c.Temp
// 				}
// 			}
// 			cmd.Env = env
//
// 			cmd.Stdout = os.Stdout
// 			cmd.Stderr = os.Stderr
// 			err := cmd.Start()
// 			if err != nil {
// 				return errors.New("Cannot run ProteinProphet")
// 			}
// 			_ = cmd.Wait()
//
// 			// copy to work directory
// 			dest := fmt.Sprintf("%s%s%s", baseDir, string(filepath.Separator), filepath.Base(output))
// 			err = sys.CopyFile(output, dest)
// 			if err != nil {
// 				return err
// 			}
//
// 		}
// 	}
//
// 	return nil
// }

func (c *ProteinProphet) appendParams(cmd *exec.Cmd) *exec.Cmd {

	if c.Noplot == true {
		cmd.Args = append(cmd.Args, "NOPLOT")
	}

	if c.Nooccam == true {
		cmd.Args = append(cmd.Args, "NOOCCAM")
	}

	if c.Softoccam == true {
		cmd.Args = append(cmd.Args, "SOFTOCCAM")
	}

	if c.Icat == true {
		cmd.Args = append(cmd.Args, "ICAT")
	}

	if c.Glyc == true {
		cmd.Args = append(cmd.Args, "GLYC")
	}

	if c.Nogroupwts == true {
		cmd.Args = append(cmd.Args, "NOGROUPWTS")
	}

	if c.NonSP == true {
		cmd.Args = append(cmd.Args, "NONSP")
	}

	if c.Accuracy == true {
		cmd.Args = append(cmd.Args, "ACCURACY")
	}

	if c.Asap == true {
		cmd.Args = append(cmd.Args, "ASAP")
	}

	if c.Refresh == true {
		cmd.Args = append(cmd.Args, "REFRESH")
	}

	if c.Normprotlen == true {
		cmd.Args = append(cmd.Args, "NORMPROTLEN")
	}

	if c.Logprobs == true {
		cmd.Args = append(cmd.Args, "LOGPROBS")
	}

	if c.Confem == true {
		cmd.Args = append(cmd.Args, "CONFEM")
	}

	if c.Allpeps == true {
		cmd.Args = append(cmd.Args, "ALLPEPS")
	}

	if c.Unmapped == true {
		cmd.Args = append(cmd.Args, "UNMAPPED")
	}

	if c.Noprotlen == true {
		cmd.Args = append(cmd.Args, "NOPROTLEN")
	}

	if c.Instances == true {
		cmd.Args = append(cmd.Args, "INSTANCES")
	}

	if c.Fpkm == true {
		cmd.Args = append(cmd.Args, "FPKM")
	}

	if c.Protmw == true {
		cmd.Args = append(cmd.Args, "PROTMW")
	}

	if c.Iprophet == true {
		cmd.Args = append(cmd.Args, "IPROPHET")
	}

	if c.Asapprophet == true {
		cmd.Args = append(cmd.Args, "ASAP_PROPHET")
	}

	if c.Delude == true {
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

	if len(c.Maxppmdiff) > 0 {
		v := fmt.Sprintf("MAXPPMDIFF%s", c.Maxppmdiff)
		cmd.Args = append(cmd.Args, v)
	}

	if len(c.Minprob) > 0 {
		v := fmt.Sprintf("MINPROB=%s", c.Minprob)
		cmd.Args = append(cmd.Args, v)
	}

	if len(c.Minindep) > 0 {
		v := fmt.Sprintf("MININDEP=%s", c.Minindep)
		cmd.Args = append(cmd.Args, v)
	}

	if len(c.Mufactor) > 0 {
		v := fmt.Sprintf("MUFACTOR=%s", c.Mufactor)
		cmd.Args = append(cmd.Args, v)
	}

	return cmd
}
