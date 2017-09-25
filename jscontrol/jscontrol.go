// Copyright 05-Sep-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
	"github.com/dedeme/go/libcgi/cgiio"
	"path"
	"strings"
  "sort"
)

func processFile(path string) {
	line = 0
	cgiio.Lines(path, func(l string) {
		processLine(l)
	})
}

type prLineT struct {
  pr string
  line int
}

type prLinesT []prLineT
func (p prLinesT) Len() int           { return len(p) }
func (p prLinesT) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p prLinesT) Less(i, j int) bool { return p[i].line < p[j].line }

func processPath(p string) {
	if cgiio.IsDirectory(p) {
		for _, f := range cgiio.List(p) {
			if f.Name() != "www" &&
        f.Name() != "old" &&
        f.Name() != "basic_201709" &&
        f.Name()[0:1] != "."  {
				processPath(path.Join(p, f.Name()))
			}
		}
	} else if cgiio.Exists(p) &&
		strings.HasSuffix(p, ".js") {
		props = make(map[string]int)
		processFile(p)
    if len(props) > 0 {
      fmt.Println(p)
      prLines := make(prLinesT, 0)
      for pr, line := range props {
        prLines = append(prLines, prLineT{pr, line})
      }
      sort.Sort(prLines)
      for _, prLine := range prLines {
        fmt.Printf("Line %d: Property \".%s\" is not a function\n",
          prLine.line, prLine.pr)
      }
      for _, prLine := range prLines {
        fmt.Printf("\"%s\": {},", prLine.pr)
      }
      fmt.Println()
    }
	}
}

func main() {
	processPath("/deme/dmjs17")
}
