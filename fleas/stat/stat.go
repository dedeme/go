// Copyright 01-Dic-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package stat

import (
	"strconv"
	"strings"
)

type Stat struct {
	Cash float64
	Buys  float32
	Sells float32
}

func New() *Stat {
	return &Stat{0, 0, 0}
}

func (st *Stat) Setup(cycle uint64, cash float64, Nbuys, Nsells uint) {
  st.Cash = (st.Cash * float64(cycle - 1) + cash) / float64(cycle)
  st.Buys = (st.Buys * float32(cycle - 1) + float32(Nbuys)) / float32(cycle)
  st.Sells = (st.Sells * float32(cycle - 1) + float32(Nsells)) / float32(cycle)
}

func (st *Stat) String() string {
	return strconv.FormatFloat(st.Cash, 'f', 2, 64) + ":" +
		strconv.FormatFloat(float64(st.Buys), 'f', 2, 32) + ":" +
		strconv.FormatFloat(float64(st.Sells), 'f', 2, 32)
}

func FromString(s string) *Stat {
	parts := strings.Split(s, ":")
	cash, _ := strconv.ParseFloat(parts[0], 64)
	buy, _ := strconv.ParseFloat(parts[1], 32)
	sell, _ := strconv.ParseFloat(parts[2], 32)
  return &Stat{cash, float32(buy), float32(sell)}
}

