// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package globals

import (
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/output"
)

func (opts *GlobalOptions) OutputOpts() output.OutputOptions {
	extra := map[string]interface{}{
		"ether": opts.Ether,
	}

	return output.OutputOptions{
		Writer:     opts.Writer,
		Chain:      opts.Chain,
		TestMode:   opts.TestMode,
		NoHeader:   opts.NoHeader,
		ShowRaw:    opts.ShowRaw,
		Verbose:    opts.Verbose,
		Format:     opts.Format,
		OutputFn:   opts.OutputFn,
		Append:     opts.Append,
		JsonIndent: "  ",
		Extra:      extra,
	}
}

func (opts *GlobalOptions) OutputOptsWithExtra(extra map[string]interface{}) output.OutputOptions {
	if extra != nil {
		extra["ether"] = opts.Ether
	}

	return output.OutputOptions{
		Writer:     opts.Writer,
		Chain:      opts.Chain,
		TestMode:   opts.TestMode,
		NoHeader:   opts.NoHeader,
		ShowRaw:    opts.ShowRaw,
		Verbose:    opts.Verbose,
		Format:     opts.Format,
		OutputFn:   opts.OutputFn,
		Append:     opts.Append,
		JsonIndent: "  ",
		Extra:      extra,
	}
}
