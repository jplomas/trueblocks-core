package scrapePkg

/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2016, 2021 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/
/*
 * This file was auto generated with makeClass --gocmds. DO NOT EDIT.
 */

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/globals"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
)

type ScrapeOptions struct {
	Modes        []string
	Action       string
	Sleep        float64
	Pin          bool
	Publish      bool
	BlockCnt     uint64
	BlockChanCnt uint64
	AddrChanCnt  uint64
	Globals      globals.GlobalOptions
	BadFlag      error
}

func (opts *ScrapeOptions) TestLog() {
	logger.TestLog(len(opts.Modes) > 0, "Modes: ", opts.Modes)
	logger.TestLog(len(opts.Action) > 0, "Action: ", opts.Action)
	logger.TestLog(opts.Sleep != 14, "Sleep: ", opts.Sleep)
	logger.TestLog(opts.Pin, "Pin: ", opts.Pin)
	logger.TestLog(opts.Publish, "Publish: ", opts.Publish)
	logger.TestLog(opts.BlockCnt != 2000, "BlockCnt: ", opts.BlockCnt)
	logger.TestLog(opts.BlockChanCnt != 10, "BlockChanCnt: ", opts.BlockChanCnt)
	logger.TestLog(opts.AddrChanCnt != 20, "AddrChanCnt: ", opts.AddrChanCnt)
	opts.Globals.TestLog()
}

func (opts *ScrapeOptions) ToCmdLine() string {
	options := ""
	if len(opts.Action) > 0 {
		options += " --action " + opts.Action
	}
	if opts.Sleep != 14 {
		options += (" --sleep " + fmt.Sprintf("%.1f", opts.Sleep))
	}
	if opts.Pin {
		options += " --pin"
	}
	if opts.Publish {
		options += " --publish"
	}
	if opts.BlockCnt != 2000 {
		options += (" --block_cnt " + fmt.Sprintf("%d", opts.BlockCnt))
	}
	if opts.BlockChanCnt != 10 {
		options += (" --block_chan_cnt " + fmt.Sprintf("%d", opts.BlockChanCnt))
	}
	if opts.AddrChanCnt != 20 {
		options += (" --addr_chan_cnt " + fmt.Sprintf("%d", opts.AddrChanCnt))
	}
	options += " " + strings.Join(opts.Modes, " ")
	options += fmt.Sprintf("%s", "") // silence go compiler for auto gen
	return options
}

func FromRequest(w http.ResponseWriter, r *http.Request) *ScrapeOptions {
	opts := &ScrapeOptions{}
	opts.Sleep = 14
	opts.BlockCnt = 2000
	opts.BlockChanCnt = 10
	opts.AddrChanCnt = 20
	for key, value := range r.URL.Query() {
		switch key {
		case "modes":
			for _, val := range value {
				s := strings.Split(val, " ") // may contain space separated items
				opts.Modes = append(opts.Modes, s...)
			}
		case "action":
			opts.Action = value[0]
		case "sleep":
			opts.Sleep = globals.ToFloat64(value[0])
		case "pin":
			opts.Pin = true
		case "publish":
			opts.Publish = true
		case "blockCnt":
			opts.BlockCnt = globals.ToUint64(value[0])
		case "blockChanCnt":
			opts.BlockChanCnt = globals.ToUint64(value[0])
		case "addrChanCnt":
			opts.AddrChanCnt = globals.ToUint64(value[0])
		default:
			if !globals.IsGlobalOption(key) {
				opts.BadFlag = validate.Usage("Invalid key ({0}) in {1} route.", key, "scrape")
				return opts
			}
		}
	}
	opts.Globals = *globals.FromRequest(w, r)

	return opts
}

var Options ScrapeOptions

func ScrapeFinishParse(args []string) *ScrapeOptions {
	// EXISTING_CODE
	Options.Modes = args
	// EXISTING_CODE
	return &Options
}