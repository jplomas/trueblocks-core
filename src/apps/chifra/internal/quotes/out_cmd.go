package quotesPkg

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

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/cmd/globals"
	"github.com/spf13/cobra"
)

var Options QuotesOptions

func RunQuotes(cmd *cobra.Command, args []string) error {
	err := Validate(cmd, args)
	if err != nil {
		return err
	}

	options := ""
	if Options.Freshen {
		options += " --freshen"
	}
	if len(Options.Period) > 0 {
		options += " --period " + Options.Period
	}
	if len(Options.Pair) > 0 {
		options += " --pair " + Options.Pair
	}
	if len(Options.Feed) > 0 {
		options += " --feed " + Options.Feed
	}
	arguments := ""
	for _, arg := range args {
		arguments += " " + arg
	}

	return globals.PassItOn("getQuotes", &Options.Globals, options, arguments)
}
