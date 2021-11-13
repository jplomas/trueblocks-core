package abis

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

type AbisOptionsType struct {
	Addrs   []string
	Known   bool
	Sol     bool
	Find    []string
	Source  bool
	Classes bool
	Globals root.RootOptionsType
}

var Options AbisOptionsType

func (opts *AbisOptionsType) TestLog() {
	logger.TestLog(len(opts.Addrs) > 0, "Addrs: ", opts.Addrs)
	logger.TestLog(opts.Known, "Known: ", opts.Known)
	logger.TestLog(opts.Sol, "Sol: ", opts.Sol)
	logger.TestLog(len(opts.Find) > 0, "Find: ", opts.Find)
	logger.TestLog(opts.Source, "Source: ", opts.Source)
	logger.TestLog(opts.Classes, "Classes: ", opts.Classes)
	opts.Globals.TestLog()
}

func FromRequest(r *http.Request) *AbisOptionsType {
	opts := &AbisOptionsType{}
	for key, value := range r.URL.Query() {
		switch key {
		case "addrs":
			opts.Addrs = append(opts.Addrs, value...)
		case "known":
			opts.Known = true
		case "sol":
			opts.Sol = true
		case "find":
			opts.Find = append(opts.Find, value...)
		case "source":
			opts.Source = true
		case "classes":
			opts.Classes = true
		}
	}
	opts.Globals = *root.FromRequest(r)

	return opts
}
