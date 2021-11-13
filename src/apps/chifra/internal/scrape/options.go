package scrape

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
 * The file was auto generated with makeClass --gocmds. DO NOT EDIT.
 */

import (
	"net/http"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/cmd/root"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

type ScrapeOptionsType struct {
	Modes        []string
	Action       string
	Sleep        float64
	Pin          bool
	Publish      bool
	BlockCnt     uint64
	BlockChanCnt uint64
	AddrChanCnt  uint64
	Globals      root.RootOptionsType
}

var Options ScrapeOptionsType

func (opts *ScrapeOptionsType) TestLog() {
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

func FromRequest(r *http.Request) *ScrapeOptionsType {
	opts := &ScrapeOptionsType{}
	for key, value := range r.URL.Query() {
		switch key {
		case "modes":
			opts.Modes = append(opts.Modes, value...)
		case "action":
			opts.Action = value[0]
		case "sleep":
			opts.Sleep = root.ToFloat(value[0])
		case "pin":
			opts.Pin = true
		case "publish":
			opts.Publish = true
		case "blockcnt":
			opts.BlockCnt = root.ToUint(value[0])
		case "blockchancnt":
			opts.BlockChanCnt = root.ToUint(value[0])
		case "addrchancnt":
			opts.AddrChanCnt = root.ToUint(value[0])
		}
	}
	opts.Globals = *root.FromRequest(r)

	return opts
}
