package configPkg

import (
	"fmt"

	"github.com/bykof/gostradamus"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/config"
)

// HandlePaths handles the paths command for the command line. Returns error only as per cobra.
func (opts *ConfigOptions) HandlePaths() error {
	chain := opts.Globals.Chain

	// TODO: This needs to be a SimpleType and use StreamMany
	dateStr := gostradamus.Now().Format("02-01|15:04:05.000")
	fmt.Printf("\nchifra status --paths:\n")
	fmt.Println(dateStr, colors.Green+"Config Path: "+colors.Off, config.PathToRootConfig())
	fmt.Println(dateStr, colors.Green+"Cache Path:  "+colors.Off, config.PathToCache(chain))
	fmt.Println(dateStr, colors.Green+"Index Path:  "+colors.Off, config.PathToIndex(chain))
	fmt.Println()
	return nil
}
