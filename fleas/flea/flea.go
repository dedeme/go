// Copyright 27-Nov-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package flea

import (
	"github.com/dedeme/go/fleas/families"
	"github.com/dedeme/go/fleas/gen"
	"github.com/dedeme/go/fleas/order"
	"github.com/dedeme/go/fleas/quote"
	"math/rand"
	"strconv"
)

const (
	Number   = 10000
	Heritage = 0.01
)

const (
	BuyAndHold     = families.BuyAndHoldId
	MovingAverage  = families.MovingAverageId
	WmovingAverage = families.WmovingAverageId
	EndFamilies    = families.EndFamilies
)

type Flea interface {
	Base() *families.Base

	Serialize() map[string]interface{}
	String() string
	Trace(nick string) string
	Process(map[string]quote.Quote)
}

func Id(f Flea) string {
	return f.Base().Id()
}

func Cycle(f Flea) uint64 {
	return f.Base().Cycle()
}

func Family(f Flea) *gen.Gen {
	return f.Base().Family()
}

func Ibex(f Flea) *gen.Gen {
	return f.Base().Ibex()
}

func Bet(f Flea) *gen.Gen {
	return f.Base().Bet()
}

func Mutability(f Flea) *gen.Gen {
	return f.Base().Mutability()
}

func Buys(f Flea) []*order.Buy {
	return f.Base().Buys()
}

func Sells(f Flea) []*order.Sell {
	return f.Base().Sells()
}

func CleanOrders(f Flea) {
	f.Base().CleanOrders()
}

func Cash(f Flea) float64 {
	return f.Base().Cash()
}

func AddCash(f Flea, money float64) {
	f.Base().AddCash(money)
}

func Nbuys(f Flea) uint {
	return f.Base().Nbuys()
}

func IncNbuys(f Flea) {
	f.Base().IncNbuys()
}

func Nsells(f Flea) uint {
	return f.Base().Nsells()
}

func IncNsells(f Flea) {
	f.Base().IncNsells()
}

func Portfolio(f Flea) map[string]uint64 {
	return f.Base().Portfolio()
}

func new(base *families.Base) Flea {
	switch base.Family().ActualOption {
	case BuyAndHold:
		return families.NewBuyAndHold(base)
	case MovingAverage:
		return families.NewMovingAverage(base)
	case WmovingAverage:
		return families.NewWmovingAverage(base)
	}
	panic("flea.new: Bad family value")
}

func New(cycle uint64) Flea {
	base := families.NewBase(cycle)
	return new(base)
}

func Restore(s map[string]interface{}) Flea {
	base := families.RestoreBase(s)
	switch base.Family().ActualOption {
	case BuyAndHold:
		return families.RestoreBuyAndHold(base, s)
	case MovingAverage:
		return families.RestoreMovingAverage(base, s)
	case WmovingAverage:
		return families.RestoreWmovingAverage(base, s)
	}
	panic("flea.Restore: Bad family value")
}

func Mutate(f Flea, cycle uint64) Flea {
	mrange := (Mutability(f).ActualOption + 10)
	base := families.MutateBase(f.Base(), mrange, cycle)

	family := base.Family().ActualOption
	if family == Family(f).ActualOption {
		switch family {
		case BuyAndHold:
			return families.MutateBuyAndHold(
				base, f.(*families.BuyAndHold), mrange)
		case MovingAverage:
			return families.MutateMovingAverage(
				base, f.(*families.MovingAverage), mrange)
		case WmovingAverage:
			return families.MutateWmovingAverage(
				base, f.(*families.WmovingAverage), mrange)
		}
		panic("flea.Mutate: Bad family value")
	}
	return new(base)
}

func traceIncomes(first bool, f Flea) {
	tracePortfolio := func() {
    portfolioTx := "Result portfolio:"
    if first {
      portfolioTx = "Initial porfolio:"
    }
		families.TraceTx = append(families.TraceTx, portfolioTx)
		for k, v := range Portfolio(f) {
			families.TraceTx = append(families.TraceTx,
				"  "+k+": "+strconv.FormatUint(v, 10))
		}
	}
	tracePortfolio()
	families.TraceTx = append(families.TraceTx,
		"Cash: "+strconv.FormatFloat(Cash(f), 'f', 2, 64))

}

func Process(f Flea, qs map[string]quote.Quote) {
	buys := Buys(f)
	for i := range buys {
		j := rand.Intn(i + 1)
		buys[i], buys[j] = buys[j], buys[i]
	}
	for _, b := range buys {
		if b.Money <= Cash(f) {
			params := ""
			if families.Trace == Id(f) {
				params = "Parameters:\n" + f.Trace(b.Nick)
			}

			stocks, cash := b.DoBuy(qs)
			if stocks > 0 {

				if families.Trace == Id(f) {
					traceIncomes(true, f)
					families.TraceTx = append(families.TraceTx, params)
					families.TraceTx = append(families.TraceTx,
						"Buy: "+b.String()+"\n     "+qs[b.Nick].String())
				}

				if ss, ok := Portfolio(f)[b.Nick]; ok {
					Portfolio(f)[b.Nick] = ss + stocks
				} else {
					Portfolio(f)[b.Nick] = stocks
				}

				AddCash(f, -cash)
				IncNbuys(f)

				if families.Trace == Id(f) {
					traceIncomes(false, f)
				}
			}
		}
	}

	for _, s := range Sells(f) {
		cash := s.DoSell(qs)
		if cash > 0 {
			if families.Trace == Id(f) {
				traceIncomes(true, f)
				families.TraceTx = append(families.TraceTx,
					"Parameters:\n"+f.Trace(s.Nick))
				families.TraceTx = append(families.TraceTx,
					"Sell: "+s.String()+"\n      "+qs[s.Nick].String())
			}

			delete(Portfolio(f), s.Nick)
			AddCash(f, cash)
			IncNsells(f)

			if families.Trace == Id(f) {
				traceIncomes(false, f)
			}
		}
	}

	if families.Trace == Id(f) {
		families.TraceTx = append(families.TraceTx,
			"---------------------------")
	}

	CleanOrders(f)
	f.Process(qs)
}
