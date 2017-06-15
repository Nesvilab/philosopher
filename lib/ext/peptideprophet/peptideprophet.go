package peptideprophet

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	unix "github.com/prvst/philosopher-source/lib/ext/peptideprophet/unix"
	wPeP "github.com/prvst/philosopher-source/lib/ext/peptideprophet/win"
	"github.com/prvst/philosopher-source/lib/meta"
	"github.com/prvst/philosopher-source/lib/sys"
)

// PeptideProphet is the tool configuration
type PeptideProphet struct {
	meta.Data
	DefaultInteractParser       string
	DefaultRefreshParser        string
	DefaultPeptideProphetParser string
	WinInteractParser           string
	WinRefreshParser            string
	WinPeptideProphetParser     string
	UnixInteractParser          string
	UnixRefreshParser           string
	UnixPeptideProphetParser    string
	LibgccDLL                   string
	Zlib1DLL                    string
	Mv                          string
	MinPepLen                   string
	Output                      string
	Clevel                      string
	Database                    string
	Minpintt                    string
	Minpiprob                   string
	Minrtntt                    string
	Minrtprob                   string
	Rtcat                       string
	Minprob                     string
	Decoy                       string
	Ignorechg                   string
	Masswidth                   string
	Combine                     bool
	Exclude                     bool
	Leave                       bool
	Perfectlib                  bool
	Icat                        bool
	Noicat                      bool
	Zero                        bool
	Accmass                     bool
	Ppm                         bool
	Nomass                      bool
	Pi                          bool
	Rt                          bool
	Glyc                        bool
	Phospho                     bool
	Maldi                       bool
	Instrwarn                   bool
	Decoyprobs                  bool
	Nontt                       bool
	Nonmc                       bool
	Expectscore                 bool
	Nonparam                    bool
	Neggamma                    bool
	Forcedistr                  bool
	Optimizefval                bool
}

// New constructor
func New() PeptideProphet {

	var o PeptideProphet
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

	o.UnixInteractParser = o.Temp + string(filepath.Separator) + "InteractParser"
	o.UnixRefreshParser = o.Temp + string(filepath.Separator) + "RefreshParser"
	o.UnixPeptideProphetParser = o.Temp + string(filepath.Separator) + "PeptideProphetParser"
	o.WinInteractParser = o.Temp + string(filepath.Separator) + "InteractParser.exe"
	o.WinRefreshParser = o.Temp + string(filepath.Separator) + "RefreshParser.exe"
	o.WinPeptideProphetParser = o.Temp + string(filepath.Separator) + "PeptideProphetParser.exe"
	o.Mv = o.Temp + string(filepath.Separator) + "mv.exe"
	o.LibgccDLL = o.Temp + string(filepath.Separator) + "libgcc_s_dw2-1.dll"
	o.Zlib1DLL = o.Temp + string(filepath.Separator) + "zlib1.dll"

	return o
}

// Deploy PeptideProphet binaries on binary directory
func (c *PeptideProphet) Deploy() error {

	if c.OS == sys.Windows() {
		wPeP.WinInteractParser(c.WinInteractParser)
		c.DefaultInteractParser = c.WinInteractParser
		wPeP.WinRefreshParser(c.WinRefreshParser)
		c.DefaultRefreshParser = c.WinRefreshParser
		wPeP.WinPeptideProphetParser(c.WinPeptideProphetParser)
		c.DefaultPeptideProphetParser = c.WinPeptideProphetParser
		wPeP.LibgccDLL(c.LibgccDLL)
		wPeP.Zlib1DLL(c.Zlib1DLL)
		wPeP.Mv(c.Mv)
	} else {
		if strings.EqualFold(c.Distro, sys.Debian()) {
			unix.UnixInteractParser(c.UnixInteractParser)
			c.DefaultInteractParser = c.UnixInteractParser
			unix.UnixRefreshParser(c.UnixRefreshParser)
			c.DefaultRefreshParser = c.UnixRefreshParser
			unix.UnixPeptideProphetParser(c.UnixPeptideProphetParser)
			c.DefaultPeptideProphetParser = c.UnixPeptideProphetParser
		} else if strings.EqualFold(c.Distro, sys.Redhat()) {
			unix.UnixInteractParser(c.UnixInteractParser)
			c.DefaultInteractParser = c.UnixInteractParser
			unix.UnixRefreshParser(c.UnixRefreshParser)
			c.DefaultRefreshParser = c.UnixRefreshParser
			unix.UnixPeptideProphetParser(c.UnixPeptideProphetParser)
			c.DefaultPeptideProphetParser = c.UnixPeptideProphetParser
		} else {
			return errors.New("Unsupported distribution for PeptideProphet")
		}
	}

	return nil
}

// Run PeptideProphet
func (c *PeptideProphet) Run(args []string) error {

	var listedArgs []string
	for _, i := range args {
		file, _ := filepath.Abs(i)
		listedArgs = append(listedArgs, file)
	}

	// run InteractParser
	files, err := interactParser(*c, listedArgs)
	if err != nil {
		return err
	}

	for _, i := range files {
		if strings.Contains(i, c.Output) {

			// run RefreshParser
			err = refreshParser(*c, i)
			if err != nil {
				return err
			}

			// run PeptideProphetParser
			err = peptideProphet(*c, i)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

// interactParser executes InteractParser binary
func interactParser(p PeptideProphet, args []string) ([]string, error) {

	var files []string

	if p.Combine == false {

		for i := range args {

			bin := p.DefaultInteractParser
			cmd := exec.Command(bin)

			// remove one or two extensions
			datadir := filepath.Dir(strings.TrimSpace(args[i]))
			basename := filepath.Base(strings.TrimSpace(args[i]))
			name := strings.TrimSuffix(basename, filepath.Ext(basename))
			name = strings.TrimSuffix(name, filepath.Ext(name))

			// set the output name and sufix
			output := fmt.Sprintf("%s%s%s-%s.pep.xml", datadir, string(filepath.Separator), p.Output, name)
			cmd.Args = append(cmd.Args, output)
			files = append(files, output)

			pepfile, _ := filepath.Abs(args[i])
			cmd.Args = append(cmd.Args, pepfile)
			files = append(files, pepfile)

			// append the directory with the mz files
			datadir, _ = filepath.Abs(datadir)
			mzfile := fmt.Sprintf("-a%s", datadir)
			cmd.Args = append(cmd.Args, mzfile)

			// -D<path_to_database>
			if len(p.Database) > 0 {
				db, _ := filepath.Abs(p.Database)
				v := fmt.Sprintf("-D%s", db)
				cmd.Args = append(cmd.Args, v)
			}

			// -L<min_peptide_len (default 7)>
			if len(p.MinPepLen) > 0 {
				v := fmt.Sprintf("-L=%s", p.MinPepLen)
				cmd.Args = append(cmd.Args, v)
			}

			cmd.Dir = filepath.Dir(output)

			env := os.Environ()
			env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
			env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", p.Home))
			for i := range env {
				if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
					env[i] = env[i] + ";" + p.Home
				}
			}
			cmd.Env = env

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Start()
			if err != nil {
				return nil, errors.New("Could not run InteractParser")
			}
			_ = cmd.Wait()

		}

	} else {

		bin := p.DefaultInteractParser
		cmd := exec.Command(bin)

		datadir := filepath.Dir(strings.TrimSpace(args[0]))

		output := fmt.Sprintf("%s%s%s.pep.xml", datadir, string(filepath.Separator), p.Output)
		cmd.Args = append(cmd.Args, output)
		files = append(files, output)

		for i := range args {
			file, _ := filepath.Abs(args[i])
			cmd.Args = append(cmd.Args, file)
		}

		// append the directory with the mz files
		datadir, _ = filepath.Abs(datadir)
		mzfile := fmt.Sprintf("-a%s", datadir)
		cmd.Args = append(cmd.Args, mzfile)

		// -D<path_to_database>
		if len(p.Database) > 0 {
			db, _ := filepath.Abs(p.Database)
			v := fmt.Sprintf("-D%s", db)
			cmd.Args = append(cmd.Args, v)
		}

		// -L<min_peptide_len (default 7)>
		if len(p.MinPepLen) > 0 {
			v := fmt.Sprintf("-L=%s", p.MinPepLen)
			cmd.Args = append(cmd.Args, v)
		}

		cmd.Dir = filepath.Dir(output)

		env := os.Environ()
		env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
		env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", p.Temp))
		for i := range env {
			if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
				env[i] = env[i] + ";" + p.Temp
			}
		}
		cmd.Env = env

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			return nil, errors.New("Could not run InteractParser")
		}
		_ = cmd.Wait()

	}

	return files, nil
}

// refreshParser executes RefreshParser binary
func refreshParser(p PeptideProphet, file string) error {

	bin := p.DefaultRefreshParser
	cmd := exec.Command(bin)

	// string of arguments to be passed as a command
	cmd.Args = append(cmd.Args, file)

	// append the database
	if len(p.Database) > 0 {
		db, _ := filepath.Abs(p.Database)
		v := fmt.Sprintf("%s", db)
		cmd.Args = append(cmd.Args, v)
	}

	env := os.Environ()
	env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
	env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", p.Temp))
	for i := range env {
		if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
			env[i] = env[i] + ";" + p.Temp
		}
	}
	cmd.Env = env
	cmd.Dir = filepath.Dir(file)

	fmt.Println("\n  -", file)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return errors.New("Cannot run RefreshParser")
	}
	_ = cmd.Wait()

	return nil
}

// peptideProphet executes peptideprophet
func peptideProphet(p PeptideProphet, file string) error {

	bin := p.DefaultPeptideProphetParser
	cmd := exec.Command(bin)

	// string of arguments to be passed as a command
	cmd.Args = append(cmd.Args, file)

	if p.Exclude == true {
		cmd.Args = append(cmd.Args, "EXCLUDE")
	}

	if p.Leave == true {
		cmd.Args = append(cmd.Args, "LEAVE")
	}

	if p.Perfectlib == true {
		cmd.Args = append(cmd.Args, "PERFECTLIB")
	}

	if p.Icat == true {
		cmd.Args = append(cmd.Args, "ICAT")
	}

	if p.Noicat == true {
		cmd.Args = append(cmd.Args, "NOICAT")
	}

	if p.Zero == true {
		cmd.Args = append(cmd.Args, "ZERO")
	}

	if p.Accmass == true {
		cmd.Args = append(cmd.Args, "ACCMASS")
	}

	if p.Ppm == true {
		cmd.Args = append(cmd.Args, "PPM")
	}

	if p.Nomass == true {
		cmd.Args = append(cmd.Args, "NOMASS")
	}

	if p.Pi == true {
		cmd.Args = append(cmd.Args, "PI")
	}

	if p.Rt == true {
		cmd.Args = append(cmd.Args, "RT")
	}

	if p.Glyc == true {
		cmd.Args = append(cmd.Args, "GLYC")
	}

	if p.Phospho == true {
		cmd.Args = append(cmd.Args, "PHOSPHO")
	}

	if p.Maldi == true {
		cmd.Args = append(cmd.Args, "MALDI")
	}

	if p.Instrwarn == true {
		cmd.Args = append(cmd.Args, "INSTRWARN")
	}

	if p.Decoyprobs == true {
		cmd.Args = append(cmd.Args, "DECOYPROBS")
	}

	if p.Nontt == true {
		cmd.Args = append(cmd.Args, "NONTT")
	}

	if p.Nonmc == true {
		cmd.Args = append(cmd.Args, "NONMC")
	}

	if p.Expectscore == true {
		cmd.Args = append(cmd.Args, "EXPECTSCORE")
	}

	if p.Nonparam == true {
		cmd.Args = append(cmd.Args, "NONPARAM")
	}

	if p.Neggamma == true {
		cmd.Args = append(cmd.Args, "NEGGAMMA")
	}

	if p.Forcedistr == true {
		cmd.Args = append(cmd.Args, "FORCEDISTR")
	}

	if p.Nonparam == true {
		cmd.Args = append(cmd.Args, "NONPARAM")
	}

	if len(p.Masswidth) > 0 {
		v := fmt.Sprintf("MASSWIDTH=%s", p.Masswidth)
		cmd.Args = append(cmd.Args, v)
	}

	if len(p.Clevel) > 0 {
		v := fmt.Sprintf("CLEVEL=%s", p.Clevel)
		cmd.Args = append(cmd.Args, v)
	}

	if len(p.Minpintt) > 0 {
		v := fmt.Sprintf("MINPINTT=%s", p.Minpintt)
		cmd.Args = append(cmd.Args, v)
	}

	if len(p.Minpiprob) > 0 {
		v := fmt.Sprintf("MINPIPROB=%s", p.Minpiprob)
		cmd.Args = append(cmd.Args, v)
	}

	if len(p.Minrtntt) > 0 {
		v := fmt.Sprintf("MINRTNTT=%s", p.Minrtntt)
		cmd.Args = append(cmd.Args, v)
	}

	if len(p.Minrtprob) > 0 {
		v := fmt.Sprintf("MINRTPROB=%s", p.Minrtprob)
		cmd.Args = append(cmd.Args, v)
	}

	if len(p.Rtcat) > 0 {
		v := fmt.Sprintf("RTCAT=%s", p.Rtcat)
		cmd.Args = append(cmd.Args, v)
	}

	if len(p.Minprob) > 0 {
		v := fmt.Sprintf("MINPROB=%s", p.Minprob)
		cmd.Args = append(cmd.Args, v)
	}

	if len(p.Decoy) > 0 {
		v := fmt.Sprintf("DECOY=%s", p.Decoy)
		cmd.Args = append(cmd.Args, v)
	}

	cmd.Dir = filepath.Dir(file)

	env := os.Environ()
	env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
	env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", p.Temp))
	for i := range env {
		if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
			env[i] = env[i] + ";" + p.Temp
		}
	}
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return errors.New("Cannot run PeptidePophet")
	}
	_ = cmd.Wait()

	return nil
}
