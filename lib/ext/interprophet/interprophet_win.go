//go:build windows
// +build windows

package interprophet

import (
	wiPr "philosopher/lib/ext/interprophet/win"
)

// Deploy generates comet binary on workdir bin directory
func (i *InterProphet) Deploy(os, distro string) {

	wiPr.WinInterProphetParser(i.WinInterProphetParser)
	i.DefaultInterProphetParser = i.WinInterProphetParser
	wiPr.LibgccDLL(i.LibgccDLL)
	wiPr.Zlib1DLL(i.Zlib1DLL)
}
