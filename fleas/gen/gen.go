// Copyright 30-Nov-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package gen

import "math/rand"

const (
  rndMax = 100
)

type Gen struct {
  MaxOptions int
  ActualOption int
}

func New (maxOptions, actualOption int) *Gen {
  return &Gen {maxOptions, actualOption}
}

func NewRandom (maxOptions int) *Gen {
  return &Gen {
    maxOptions,
    rand.Intn(maxOptions)}
}

// Mutate mutates 'g'. Probablity to mutate is rangeTrue / rndMax (1000000)
func (g *Gen) Mutate (rangeTrue int) *Gen {
  if rand.Intn(rndMax) < rangeTrue {
    return NewRandom(g.MaxOptions)
  }
  return New(g.MaxOptions, g.ActualOption)
}

func (g *Gen) Serialize() []interface{} {
  return []interface{}{g.MaxOptions, g.ActualOption}
}

func Restore (serial []interface{}) *Gen {
  return New(int(serial[0].(float64)), int(serial[1].(float64)));
}
