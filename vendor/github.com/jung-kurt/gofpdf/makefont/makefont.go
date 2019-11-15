package main

import (
	"flag"
	"fmt"
	"github.com/jung-kurt/gofpdf/v2"
	"os"
)

func errPrintf(fmtStr string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, fmtStr, args...)
}

func showHelp() {
	errPrintf("Usage: %s [options] font_file [font_file...]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\n"+
		"font_file is the name of the TrueType file (extension .ttf), OpenType file\n"+
		"(extension .otf) or binary Type1 file (extension .pfb) from which to\n"+
		"generate a definition file. If an OpenType file is specified, it must be one\n"+
		"that is based on TrueType outlines, not PostScript outlines; this cannot be\n"+
		"determined from the file extension alone. If a Type1 file is specified, a\n"+
		"metric file with the same pathname except with the extension .afm must be\n"+
		"present.")
	errPrintf("\nExample: %s --embed --enc=../font/cp1252.map --dst=../font calligra.ttf /opt/font/symbol.pfb\n", os.Args[0])
}

func tutorialSummary(f *gofpdf.Fpdf, fileStr string) {
	if f.Ok() {
		fl, err := os.Create(fileStr)
		defer fl.Close()
		if err == nil {
			f.Output(fl)
		} else {
			f.SetError(err)
		}
	}
	if f.Ok() {
		fmt.Printf("Successfully generated %s\n", fileStr)
	} else {
		errPrintf("%s\n", f.Error())
	}
}

func main() {
	var dstDirStr, encodingFileStr string
	var err error
	var help, embed bool
	flag.StringVar(&dstDirStr, "dst", ".", "directory for output files (*.z, *.json)")
	flag.StringVar(&encodingFileStr, "enc", "cp1252.map", "code page file")
	flag.BoolVar(&embed, "embed", false, "embed font into PDF")
	flag.BoolVar(&help, "help", false, "command line usage")
	flag.Parse()
	if help {
		showHelp()
	} else {
		args := flag.Args()
		if len(args) > 0 {
			for _, fileStr := range args {
				err = gofpdf.MakeFont(fileStr, encodingFileStr, dstDirStr, os.Stderr, embed)
				if err != nil {
					errPrintf("%s\n", err)
				}
				// errPrintf("Font file [%s], Encoding file [%s], Embed [%v]\n", fileStr, encodingFileStr, embed)
			}
		} else {
			errPrintf("At least one Type1 or TrueType font must be specified\n")
			showHelp()
		}
	}
}
