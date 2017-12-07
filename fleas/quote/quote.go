// Copyright 28-Nov-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package quote

import "strconv"

// Index of Quote
const (
  Date = iota
  Open = iota
  Close = iota
  Max = iota
  Min = iota
  Vol = iota
)

// Quote is a group of values: Date, Open, Close, Max, Min, Vol. These can be
// accessed in format type:
//   q[quote.Open]
type Quote [6]interface{}

func (q Quote)String() string {
  r:= q[Date].(string)
  for i := 1; i < len(q); i++ {
    r += ":" + strconv.FormatFloat(float64(q[i].(float32)), 'f', 4, 32)
  }
  return r
}
