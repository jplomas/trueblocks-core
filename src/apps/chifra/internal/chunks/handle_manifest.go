// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package chunksPkg

import (
	"sort"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/manifest"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

func (opts *ChunksOptions) HandleManifest(blockNums []uint64) error {
	maxTestItems = 10

	source := manifest.FromCache
	if opts.Remote {
		source = manifest.FromContract
	}

	man, err := manifest.ReadManifest(opts.Globals.Chain, source)
	if err != nil {
		return err
	}

	sort.Slice(man.Chunks, func(i, j int) bool {
		iPin := man.Chunks[i]
		jPin := man.Chunks[j]
		return iPin.Range < jPin.Range
	})

	if opts.Globals.Format == "txt" || opts.Globals.Format == "csv" {
		defer opts.Globals.RenderFooter()
		err := opts.Globals.RenderHeader(types.SimpleChunkRecord{}, &opts.Globals.Writer, opts.Globals.Format, opts.Globals.ApiMode, opts.Globals.NoHeader, true)
		if err != nil {
			return err
		}

		pinList := make([]types.SimpleChunkRecord, len(man.Chunks))
		for i := range man.Chunks {
			pinList[i].Range = man.Chunks[i].Range
			pinList[i].BloomHash = string(man.Chunks[i].BloomHash)
			pinList[i].IndexHash = string(man.Chunks[i].IndexHash)
		}

		for i, r := range pinList {
			if opts.Globals.TestMode && i > maxTestItems {
				continue
			}
			err := opts.Globals.RenderObject(r, i == 0)
			if err != nil {
				return err
			}
		}

		return nil

	} else {
		defer opts.Globals.RenderFooter()
		err = opts.Globals.RenderHeader(types.SimpleManifest{}, &opts.Globals.Writer, opts.Globals.Format, opts.Globals.ApiMode, opts.Globals.NoHeader, true)
		if err != nil {
			return err
		}

		if opts.Globals.TestMode {
			if len(man.Chunks) > maxTestItems {
				s := len(man.Chunks) - maxTestItems - 1
				e := len(man.Chunks) - 1
				man.Chunks = man.Chunks[s:e]
			}
		}

		return opts.Globals.RenderManifest(opts.Globals.Writer, man)
	}
}
