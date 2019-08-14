package peptideprophet

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/prvst/philosopher/lib/err"
	unix "github.com/prvst/philosopher/lib/ext/peptideprophet/unix"
	wPeP "github.com/prvst/philosopher/lib/ext/peptideprophet/win"
	"github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
)

// PeptideProphet is the main tool data configuration structure
type PeptideProphet struct {
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
}

// New constructor
func New(temp string) PeptideProphet {

	var self PeptideProphet

	self.UnixInteractParser = temp + string(filepath.Separator) + "InteractParser"
	self.UnixRefreshParser = temp + string(filepath.Separator) + "RefreshParser"
	self.UnixPeptideProphetParser = temp + string(filepath.Separator) + "PeptideProphetParser"
	self.WinInteractParser = temp + string(filepath.Separator) + "InteractParser.exe"
	self.WinRefreshParser = temp + string(filepath.Separator) + "RefreshParser.exe"
	self.WinPeptideProphetParser = temp + string(filepath.Separator) + "PeptideProphetParser.exe"
	self.Mv = temp + string(filepath.Separator) + "mv.exe"
	self.LibgccDLL = temp + string(filepath.Separator) + "libgcc_s_dw2-1.dll"
	self.Zlib1DLL = temp + string(filepath.Separator) + "zlib1.dll"

	return self
}

// Run is the main entry point for peptideprophet
func Run(m met.Data, args []string) met.Data {

	var pep = New(m.Temp)

	if len(m.PeptideProphet.Database) < 1 {
		err.FatalCustom(errors.New("You need to provide a protein database"))
	}

	// get the database tag from database command
	if len(m.PeptideProphet.Decoy) == 0 {
		m.PeptideProphet.Decoy = m.Database.Tag
	}

	// deploy the binaries
	pep.Deploy(m.OS, m.Distro)

	// run
	pep.Execute(m.PeptideProphet, m.Home, m.Temp, args)

	m.PeptideProphet.InputFiles = args

	return m
}

// Deploy PeptideProphet binaries on binary directory
func (p *PeptideProphet) Deploy(os, distro string) {

	if os == sys.Windows() {
		wPeP.WinInteractParser(p.WinInteractParser)
		p.DefaultInteractParser = p.WinInteractParser
		wPeP.WinRefreshParser(p.WinRefreshParser)
		p.DefaultRefreshParser = p.WinRefreshParser
		wPeP.WinPeptideProphetParser(p.WinPeptideProphetParser)
		p.DefaultPeptideProphetParser = p.WinPeptideProphetParser
		wPeP.LibgccDLL(p.LibgccDLL)
		wPeP.Zlib1DLL(p.Zlib1DLL)
		wPeP.Mv(p.Mv)
	} else {
		if strings.EqualFold(distro, sys.Debian()) {
			unix.UnixInteractParser(p.UnixInteractParser)
			p.DefaultInteractParser = p.UnixInteractParser
			unix.UnixRefreshParser(p.UnixRefreshParser)
			p.DefaultRefreshParser = p.UnixRefreshParser
			unix.UnixPeptideProphetParser(p.UnixPeptideProphetParser)
			p.DefaultPeptideProphetParser = p.UnixPeptideProphetParser
		} else if strings.EqualFold(distro, sys.Redhat()) {
			unix.UnixInteractParser(p.UnixInteractParser)
			p.DefaultInteractParser = p.UnixInteractParser
			unix.UnixRefreshParser(p.UnixRefreshParser)
			p.DefaultRefreshParser = p.UnixRefreshParser
			unix.UnixPeptideProphetParser(p.UnixPeptideProphetParser)
			p.DefaultPeptideProphetParser = p.UnixPeptideProphetParser
		} else {
			err.UnsupportedDistribution()
		}
	}

	return
}

// Execute PeptideProphet
func (p PeptideProphet) Execute(params met.PeptideProphet, home, temp string, args []string) []string {

	var output []string

	var listedArgs []string
	for _, i := range args {
		file, _ := filepath.Abs(i)
		listedArgs = append(listedArgs, file)
	}

	// run InteractParser
	files := interactParser(p, params, home, temp, listedArgs)

	for _, i := range files {
		if strings.Contains(i, params.Output) {

			// run RefreshParser
			refreshParser(p, i, params.Database, params.Output, temp)

			// run PeptideProphetParser
			output = peptideProphet(p, params, temp, i)
		}
	}

	return output
}

// interactParser executes InteractParser binary
func interactParser(p PeptideProphet, params met.PeptideProphet, home, temp string, args []string) []string {

	var files []string

	if params.Combine == false {

		for i := range args {

			bin := p.DefaultInteractParser
			cmd := exec.Command(bin)

			// remove one or two extensions
			datadir := filepath.Dir(strings.TrimSpace(args[i]))
			basename := filepath.Base(strings.TrimSpace(args[i]))
			name := strings.TrimSuffix(basename, filepath.Ext(basename))
			name = strings.TrimSuffix(name, filepath.Ext(name))

			// set the output name and sufix
			output := fmt.Sprintf("%s%s%s-%s.pep.xml", datadir, string(filepath.Separator), params.Output, name)
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
			if len(params.Database) > 0 {
				db, _ := filepath.Abs(params.Database)
				v := fmt.Sprintf("-D%s", db)
				cmd.Args = append(cmd.Args, v)
			}

			// -L<min_peptide_len (default 7)>
			if params.MinPepLen != 7 {
				v := fmt.Sprintf("-L=%d", params.MinPepLen)
				cmd.Args = append(cmd.Args, v)
			}

			if len(params.Enzyme) > 0 {
				v := fmt.Sprintf("-E%s", params.Enzyme)
				cmd.Args = append(cmd.Args, v)
			}

			cmd.Dir = filepath.Dir(output)

			env := os.Environ()
			env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
			env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", home))
			for i := range env {
				if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
					env[i] = env[i] + ";" + home
				}
			}
			cmd.Env = env

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			e := cmd.Start()
			if e != nil {
				err.ExecutingBinary(e)
			}
			_ = cmd.Wait()

		}

	} else {

		bin := p.DefaultInteractParser
		cmd := exec.Command(bin)

		datadir := filepath.Dir(strings.TrimSpace(args[0]))

		output := fmt.Sprintf("%s%s%s.pep.xml", datadir, string(filepath.Separator), params.Output)
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
		if len(params.Database) > 0 {
			db, _ := filepath.Abs(params.Database)
			v := fmt.Sprintf("-D%s", db)
			cmd.Args = append(cmd.Args, v)
		}

		// -L<min_peptide_len (default 7)>
		if params.MinPepLen != 7 {
			v := fmt.Sprintf("-L=%d", params.MinPepLen)
			cmd.Args = append(cmd.Args, v)
		}

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
			err.ExecutingBinary(e)
		}
		_ = cmd.Wait()

	}

	return files
}

// refreshParser executes RefreshParser binary
func refreshParser(p PeptideProphet, file, database, output, temp string) {

	bin := p.DefaultRefreshParser
	cmd := exec.Command(bin)

	// string of arguments to be passed as a command
	cmd.Args = append(cmd.Args, file)

	// append the database
	if len(database) > 0 {
		db, _ := filepath.Abs(database)
		v := fmt.Sprintf("%s", db)
		cmd.Args = append(cmd.Args, v)
	}

	env := os.Environ()
	env = append(env, fmt.Sprintf("XML_ONLY=%d", 1))
	env = append(env, fmt.Sprintf("WEBSERVER_ROOT=%s", temp))
	for i := range env {
		if strings.HasPrefix(strings.ToUpper(env[i]), "PATH=") {
			env[i] = env[i] + ";" + temp
		}
	}
	cmd.Env = env
	cmd.Dir = filepath.Dir(file)

	fmt.Println("\n  -", file)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	e := cmd.Start()
	if e != nil {
		err.ExecutingBinary(e)
	}
	_ = cmd.Wait()

	return
}

// peptideProphet executes peptideprophet
func peptideProphet(p PeptideProphet, params met.PeptideProphet, temp, file string) []string {
	bin := p.DefaultPeptideProphetParser
	cmd := exec.Command(bin)

	// string of arguments to be passed as a command
	cmd.Args = append(cmd.Args, file)

	if params.Exclude == true {
		cmd.Args = append(cmd.Args, "EXCLUDE")
	}

	if params.Leave == true {
		cmd.Args = append(cmd.Args, "LEAVE")
	}

	if params.Perfectlib == true {
		cmd.Args = append(cmd.Args, "PERFECTLIB")
	}

	if params.Icat == true {
		cmd.Args = append(cmd.Args, "ICAT")
	}

	if params.Noicat == true {
		cmd.Args = append(cmd.Args, "NOICAT")
	}

	if params.Zero == true {
		cmd.Args = append(cmd.Args, "ZERO")
	}

	if params.Accmass == true {
		cmd.Args = append(cmd.Args, "ACCMASS")
	}

	if params.Ppm == true {
		cmd.Args = append(cmd.Args, "PPM")
	}

	if params.Nomass == true {
		cmd.Args = append(cmd.Args, "NOMASS")
	}

	if params.Pi == true {
		cmd.Args = append(cmd.Args, "PI")
	}

	if params.Rt == true {
		cmd.Args = append(cmd.Args, "RT")
	}

	if params.Glyc == true {
		cmd.Args = append(cmd.Args, "GLYC")
	}

	if params.Phospho == true {
		cmd.Args = append(cmd.Args, "PHOSPHO")
	}

	if params.Maldi == true {
		cmd.Args = append(cmd.Args, "MALDI")
	}

	if params.Instrwarn == true {
		cmd.Args = append(cmd.Args, "INSTRWARN")
	}

	if params.Decoyprobs == true {
		cmd.Args = append(cmd.Args, "DECOYPROBS")
	}

	if params.Nontt == true {
		cmd.Args = append(cmd.Args, "NONTT")
	}

	if params.Nonmc == true {
		cmd.Args = append(cmd.Args, "NONMC")
	}

	if params.Expectscore == true {
		cmd.Args = append(cmd.Args, "EXPECTSCORE")
	}

	if params.Nonparam == true {
		cmd.Args = append(cmd.Args, "NONPARAM")
	}

	if params.Neggamma == true {
		cmd.Args = append(cmd.Args, "NEGGAMMA")
	}

	if params.Forcedistr == true {
		cmd.Args = append(cmd.Args, "FORCEDISTR")
	}

	if params.Nonparam == true {
		cmd.Args = append(cmd.Args, "NONPARAM")
	}

	if params.Masswidth != 5.0 {
		v := fmt.Sprintf("MASSWIDTH=%.4f", params.Masswidth)
		cmd.Args = append(cmd.Args, v)
	}

	if params.Clevel != 0 {
		v := fmt.Sprintf("CLEVEL=%d", params.Clevel)
		cmd.Args = append(cmd.Args, v)
	}

	if params.Minpintt != 2 {
		v := fmt.Sprintf("MINPINTT=%d", params.Minpintt)
		cmd.Args = append(cmd.Args, v)
	}

	if params.Minpiprob != 0.9 {
		v := fmt.Sprintf("MINPIPROB=%.4f", params.Minpiprob)
		cmd.Args = append(cmd.Args, v)
	}

	if params.Minrtntt != 2 {
		v := fmt.Sprintf("MINRTNTT=%d", params.Minrtntt)
		cmd.Args = append(cmd.Args, v)
	}

	if params.Minrtprob != 0.9 {
		v := fmt.Sprintf("MINRTPROB=%.4f", params.Minrtprob)
		cmd.Args = append(cmd.Args, v)
	}

	if len(params.Rtcat) > 0 {
		v := fmt.Sprintf("RTCAT=%s", params.Rtcat)
		cmd.Args = append(cmd.Args, v)
	}

	if params.Minprob != 0.05 {
		v := fmt.Sprintf("MINPROB=%.4f", params.Minprob)
		cmd.Args = append(cmd.Args, v)
	}

	if len(params.Decoy) > 0 {
		v := fmt.Sprintf("DECOY=%s", params.Decoy)
		cmd.Args = append(cmd.Args, v)
	}

	cmd.Dir = filepath.Dir(file)

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
		err.ExecutingBinary(e)
	}
	_ = cmd.Wait()

	// collect all resulting files
	var output []string
	for _, i := range cmd.Args {
		if strings.Contains(i, "interact") {
			output = append(output, i)
		}
	}

	return output
}
