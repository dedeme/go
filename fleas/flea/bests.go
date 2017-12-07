// Copyright 01-Dic-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package flea

type BestData struct {
	F    Flea
	Cash float64
	Buys  uint
	Sells uint
}

type Bests [Heritage * Number]BestData

func MkBests() *Bests {
	var bests Bests
	for i := range bests {
		bests[i].Cash = -1
	}
	return &bests
}

func AddBest(b *Bests, cycle uint64, f Flea, cash float64, Nbuys, Nsells uint) {
  fCash := cash
  if cycle > 0 {
    fCash = cash * (1 - float64(Cycle(f)) / float64(cycle))
  }
  fIsUp := func (other BestData) bool {
    if other.Cash < 0 {
      return true
    }
    otherCash := other.Cash
    if cycle > 0 {
      otherCash = otherCash * (1 - float64(Cycle(other.F)) / float64(cycle))
    }
    return fCash > otherCash
  }

	if fIsUp(b[len(b)-1]) {
		i := 0
		for ; i < len(b); i++ {
			if fIsUp(b[i]) {
				break
			}
		}
		tmp := b[i]
		b[i] = BestData{
			f,
			cash,
			Nbuys,
			Nsells}
		i++
		for ; i < len(b); i++ {
			tmp2 := b[i]
			b[i] = tmp
			tmp = tmp2
		}
	}
}
