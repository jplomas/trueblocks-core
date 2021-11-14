package output

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/tabwriter"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/cmd/globals"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

// Formatter type represents a function that can be used to format data
type Formatter func(data interface{}, opts *globals.GlobalOptionsType) ([]byte, error)

// The map below maps format string to Formatter
var formatStringToFormatter = map[string]Formatter{
	"api":  JsonFormatter,
	"json": JsonFormatter,
	"csv":  CsvFormatter,
	"txt":  TxtFormatter,
	"tab":  TabFormatter,
}

// JsonFormatter turns data into JSON
func JsonFormatter(data interface{}, opts *globals.GlobalOptionsType) ([]byte, error) {
	formatted := &JsonFormatted{}
	err, ok := data.(error)
	if ok {
		formatted.Errors = []string{
			err.Error(),
		}
	} else {
		formatted.Data = data
	}

	return AsJsonBytes(formatted, opts)
}

// TxtFormatter turns data into TSV string
func TxtFormatter(data interface{}, opts *globals.GlobalOptionsType) ([]byte, error) {
	out := bytes.Buffer{}
	tsv, err := AsTsv(data)
	if err != nil {
		return nil, err
	}
	_, err = out.Write(tsv)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// TabFormatter turns data into a table (string)
func TabFormatter(data interface{}, opts *globals.GlobalOptionsType) ([]byte, error) {
	tabOutput := &bytes.Buffer{}
	tab := tabwriter.NewWriter(tabOutput, 0, 0, 2, ' ', 0)
	tsv, err := AsTsv(data)
	if err != nil {
		return nil, err
	}
	tab.Write(tsv)
	err = tab.Flush()

	return tabOutput.Bytes(), err
}

// This type is used to carry CSV layout information
type CsvFormatted struct {
	Header  []string
	Content [][]string
}

// CsvFormatter turns a type into CSV string. It uses custom code instead of
// Go's encoding/csv to maintain compatibility with C++ output, which
// quotes each item. encoding/csv would double-quote a quoted string...
func CsvFormatter(i interface{}, opts *globals.GlobalOptionsType) ([]byte, error) {
	records, err := ToStringRecords(i, true)
	if err != nil {
		return nil, err
	}
	result := []string{}
	// We have to join records in one row with a ","
	for _, row := range records {
		result = append(result, strings.Join(row, ","))
	}

	// Now we need to join all rows with a newline. We also add one
	// final newline to be concise with both Go's encoding/csv and C++
	// version
	return []byte(
		strings.Join(result, "\n") + "\n",
	), nil
}

// Output converts data into the given format and writes to where writer
func Output(opts *globals.GlobalOptionsType, where io.Writer, format string, data interface{}) error {
	nonEmptyFormat := format
	if format == "" || format == "none" {
		if utils.IsApiMode() {
			nonEmptyFormat = "api"
		} else {
			nonEmptyFormat = "txt"
		}
	}
	if nonEmptyFormat == "txt" {
		_, ok := where.(http.ResponseWriter)
		// We would never want to use tab format in server environment
		if utils.IsTerminal() && !ok {
			nonEmptyFormat = "tab"
		}
	}
	formatter, ok := formatStringToFormatter[nonEmptyFormat]
	if !ok {
		return fmt.Errorf("unsupported format %s", format)
	}

	output, err := formatter(data, opts)
	if err != nil {
		return err
	}

	where.Write(output)
	// Maintain newline compatibility with C++ version
	if nonEmptyFormat == "json" || nonEmptyFormat == "api" {
		where.Write([]byte{'\n'})
	}
	return nil
}
