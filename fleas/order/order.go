// Copyright 2-Dec-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package order

import (
	"github.com/dedeme/go/fleas/quote"
	"strconv"
)

type Buy struct {
	Nick  string
	Money float64
}

func NewBuy(nick string, money float64) *Buy {
	return &Buy{Nick: nick, Money: money}
}

func (b *Buy) DoBuy(qs map[string]quote.Quote) (uint64, float64) {
	q := qs[b.Nick]
	price := q[quote.Open].(float32)

	if price > 0 {
		omoney := b.Money
		broker := 9.75
		if omoney > 25000 {
			broker = omoney * 0.001
		}
		bolsa := 4.65 + omoney*0.00012
		if omoney > 140000 {
			bolsa = 13.4
		} else if omoney > 70000 {
			bolsa = 9.2 + omoney*0.00003
		} else if omoney > 35000 {
			bolsa = 6.4 + omoney*0.00007
		}

		money := omoney - broker - bolsa
		stocks := uint64(money / float64(price))
		for {
			cost := float64(stocks) * float64(price)
			broker := 9.75
			if cost > 25000 {
				broker = cost * 0.001
			}
			bolsa := 4.65 + cost*0.00012
			if cost > 140000 {
				bolsa = 13.4
			} else if cost > 70000 {
				bolsa = 9.2 + cost*0.00003
			} else if cost > 35000 {
				bolsa = 6.4 + cost*0.00007
			}
			cost += broker + bolsa
			if cost <= omoney {
				return stocks, cost
			}
			stocks--
		}
	}
	return 0, 0
}

func (b *Buy) String() string {
	return b.Nick + ":" +
		strconv.FormatFloat(b.Money, 'f', 2, 64)
}

type Sell struct {
	Nick   string
	stocks uint64
}

func NewSell(nick string, stocks uint64) *Sell {
	return &Sell{Nick: nick, stocks: stocks}
}

func (s *Sell) DoSell(qs map[string]quote.Quote) float64 {
	q := qs[s.Nick]
	price := q[quote.Open].(float32)

	if price > 0 {
		return Income(s.stocks, price)
	}
	return 0
}

func (s *Sell) String() string {
	return s.Nick + ":" +
		strconv.FormatUint(s.stocks, 10)
}

func Income(stocks uint64, price float32) float64 {
	income := float64(stocks) * float64(price)
	broker := 9.75
	if income > 25000 {
		broker = income * 0.001
	}
	bolsa := 4.65 + income*0.00012
	if income > 140000 {
		bolsa = 13.4
	} else if income > 70000 {
		bolsa = 9.2 + income*0.00003
	} else if income > 35000 {
		bolsa = 6.4 + income*0.00007
	}
	return income - broker - bolsa
}

