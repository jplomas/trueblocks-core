package names

import (
	"bytes"
	"encoding/csv"
	"testing"

	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/types"
)

func TestNewNameWriter(t *testing.T) {
	buf := bytes.Buffer{}
	w := NewNameWriter(&buf)

	if w.WriteHeader != true {
		t.Fatal("expected WriteHeader to be set")
	}
	if len(w.Header) != 13 {
		t.Fatal("expected Header to have 13 items")
	}
	if w.csvWriter == nil {
		t.Fatal("expected csvWriter to be initialized")
	}
}

func TestNameWriter_writeHeader(t *testing.T) {
	buf := bytes.Buffer{}
	w := NewNameWriter(&buf)

	var err error
	if err = w.writeHeader(); err != nil {
		t.Fatal(err)
	}
	w.Flush()

	if err = w.Error(); err != nil {
		t.Fatal(err)
	}

	bufReader := bytes.NewReader(buf.Bytes())
	r := csv.NewReader(bufReader)
	r.Comma = '\t'
	result, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if lr := len(result[0]); lr != len(defaultHeader) {
		t.Fatal("wrong header item count:", lr)
	}

	for index := range defaultHeader {
		if result[0][index] != defaultHeader[index] {
			t.Fatal("expected", defaultHeader[index], "but got", result[0][index])
		}
	}
}

func TestNameWriterWrite(t *testing.T) {
	buf := bytes.Buffer{}
	w := NewNameWriter(&buf)

	var err error
	name := &types.SimpleName{
		Address: base.HexToAddress("0x1f9090aae28b8a3dceadf281b0f12828e676c326"),
		Name:    "John",
	}
	if err = w.Write(name); err != nil {
		t.Fatal(err)
	}
	w.Flush()

	if err = w.Error(); err != nil {
		t.Fatal(err)
	}

	bufReader := bytes.NewReader(buf.Bytes())
	r := csv.NewReader(bufReader)
	r.Comma = '\t'
	result, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if lr := len(result); lr != 2 {
		t.Fatal("wrong line count:", lr)
	}

	addressIndex := -1
	nameIndex := -1
	for index, headerItem := range result[0] {
		if headerItem == "address" {
			addressIndex = index
		}
		if headerItem == "name" {
			nameIndex = index
		}
	}
	if nameIndex == -1 || addressIndex == -1 {
		t.Fatal("cannot read header")
	}

	for index, item := range result[1] {
		if index == addressIndex {
			if item != name.Address.Hex() {
				t.Fatal("wrong address:", item)
			}
		}
		if index == nameIndex {
			if item != name.Name {
				t.Fatal("wrong name:", name)
			}
		}
	}
}

func TestNameWriterWriteNoHeader(t *testing.T) {
	buf := bytes.Buffer{}
	w := NewNameWriter(&buf)
	w.WriteHeader = false

	var err error
	name := &types.SimpleName{
		Address: base.HexToAddress("0x1f9090aae28b8a3dceadf281b0f12828e676c326"),
		Name:    "John",
	}
	if err = w.Write(name); err != nil {
		t.Fatal(err)
	}
	w.Flush()

	if err = w.Error(); err != nil {
		t.Fatal(err)
	}

	bufReader := bytes.NewReader(buf.Bytes())
	r := csv.NewReader(bufReader)
	r.Comma = '\t'
	result, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if lr := len(result); lr != 1 {
		t.Fatal("expected only 1 line of output, but got:", lr)
	}
}
