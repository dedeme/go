// Copyright 27-Nov-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package families

import (
	"github.com/dedeme/go/fleas/order"
	"github.com/dedeme/go/fleas/quote"
)

type BuyAndHold struct {
	base *Base
}

func NewBuyAndHold(base *Base) *BuyAndHold {
	return &BuyAndHold{base}
}

func (bah *BuyAndHold) Base() *Base {
	return bah.base
}

func (bah *BuyAndHold) String() string {
	return bah.base.String()
}

func (bah *BuyAndHold) Trace(nick string) string {
	return ""
}

func (bah *BuyAndHold) Serialize() map[string]interface{} {
	return bah.base.Serialize()
}

func RestoreBuyAndHold(base *Base, s map[string]interface{}) *BuyAndHold {
	return NewBuyAndHold(RestoreBase(s))
}

func MutateBuyAndHold(base *Base, f *BuyAndHold, mrange int) *BuyAndHold {
	return &BuyAndHold{base}
}

func (bah *BuyAndHold) Process(qs map[string]quote.Quote) {
	b := bah.base
	bet := float64((b.Bet().ActualOption + 1) * 5000)
	ibex := b.Ibex().ActualOption
	canBuy := bet <= b.Cash()
	os := make([]*order.Buy, 0)
	for k, _ := range qs {
		if ibex == 0 {
			if _, ok := Ibex[k]; ok {
				continue
			}
		} else if ibex == 1 {
			if _, ok := Ibex[k]; !ok {
				continue
			}
		}

		if canBuy {
			os = append(os, order.NewBuy(k, bet))
		}
	}
	b.buys = os
}
