// Copyright 4-Dec-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package families

const (
	InitialCash = 125000
)

const (
	BuyAndHoldId     = 0
	MovingAverageId  = 1
	WmovingAverageId = 2
	EndFamilies      = 3
)

var Trace = ""

var Names [EndFamilies]string

var Ibex map[string]struct{}

var TraceTx []string

func Ini(ibex map[string]struct{}) {
  Ibex = ibex
	Names[BuyAndHoldId] = "BuyAndHold"
	Names[MovingAverageId] = "MovingAverage"
	Names[WmovingAverageId] = "WmovingAverage"
}
