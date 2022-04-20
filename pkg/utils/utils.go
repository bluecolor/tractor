package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
)

// RespondwithJSON write json response format
func RespondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	fmt.Println(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func ErrorWithJSON(w http.ResponseWriter, code int, err error) {
	RespondwithJSON(w, code, map[string]string{"error": err.Error()})
}
func MapToStruct(m map[string]interface{}, s interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s)
}
func TwoToOneDim(data [][]interface{}) []driver.Value {
	var result []driver.Value
	for _, row := range data {
		for _, col := range row {
			result = append(result, col)
		}
	}
	return result
}
func ToChunksStr(items []string, chunkCount int) [][]string {
	cc := int(math.Min(float64(int(len(items)/chunkCount)), 1))
	chunks := [][]string{}
	chunk := []string{}
	for i := 0; i < len(items); i++ {
		chunk = append(chunk, items[i])
		if (i%cc == 0 && i != 0) || i == len(items)-1 {
			chunks = append(chunks, chunk)
			chunk = []string{}
		}
	}
	return chunks
}

func ToChunks[T comparable](items []T, chunkCount int) [][]T {
	cc := int(math.Min(float64(int(len(items)/chunkCount)), 1))
	chunks := [][]T{}
	chunk := []T{}
	for i := 0; i < len(items); i++ {
		chunk = append(chunk, items[i])
		if (i%cc == 0 && i != 0) || i == len(items)-1 {
			chunks = append(chunks, chunk)
			chunk = []T{}
		}
	}
	return chunks
}

// Taken from https://github.com/lithammer/dedent
// Dedent removes any common leading whitespace from every line in text.
//
// This can be used to make multiline strings to line up with the left edge of
// the display, while still presenting them in the source code in indented
// form.
func Dedent(text string) string {
	var margin string
	var (
		whitespaceOnly    = regexp.MustCompile("(?m)^[ \t]+$")
		leadingWhitespace = regexp.MustCompile("(?m)(^[ \t]*)(?:[^ \t\n])")
	)

	text = whitespaceOnly.ReplaceAllString(text, "")
	indents := leadingWhitespace.FindAllStringSubmatch(text, -1)

	// Look for the longest leading string of spaces and tabs common to all
	// lines.
	for i, indent := range indents {
		if i == 0 {
			margin = indent[1]
		} else if strings.HasPrefix(indent[1], margin) {
			// Current line more deeply indented than previous winner:
			// no change (previous winner is still on top).
			continue
		} else if strings.HasPrefix(margin, indent[1]) {
			// Current line consistent with and no deeper than previous winner:
			// it's the new winner.
			margin = indent[1]
		} else {
			// Current line and previous winner have no common whitespace:
			// there is no margin.
			margin = ""
			break
		}
	}

	if margin != "" {
		text = regexp.MustCompile("(?m)^"+margin).ReplaceAllString(text, "")
	}
	return text
}
func GetFileExtension(filename string) string {
	tokens := strings.Split(filename, ".")
	if len(tokens) == 1 {
		return ""
	}
	return tokens[len(tokens)-1]
}
func SplitExt(filename string) (name string, ext string) {
	tokens := strings.Split(filename, ".")
	if len(tokens) == 1 {
		return filename, ""
	}
	return strings.Join(tokens[:len(tokens)-1], "."), tokens[len(tokens)-1]
}
