// Copyright 3-Dec-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package families

import (
	"encoding/json"
	"github.com/dedeme/go/fleas/gen"
	"github.com/dedeme/go/fleas/order"
	"github.com/dedeme/go/fleas/quote"
	"strconv"
)

type wmovingAverageData struct {
	closes    []float32
	canBuy    bool
	canSell   bool
	total     float32
	numerator float32
}

type WmovingAverage struct {
	base *Base

	length    *gen.Gen // (v + 1) * 5
	buyStrip  *gen.Gen // v * 0.005
	sellStrip *gen.Gen // v * 0.005

	data map[string]*wmovingAverageData
}

func (wa *WmovingAverage) Base() *Base {
	return wa.base
}

func NewWmovingAverage(base *Base) *WmovingAverage {
	return &WmovingAverage{
		base: base,

		length:    gen.NewRandom(30),
		buyStrip:  gen.NewRandom(11),
		sellStrip: gen.NewRandom(11),

		data: make(map[string]*wmovingAverageData),
	}
}

func (wa *WmovingAverage) String() string {
	return wa.base.String() + ":" +
		"length " +
		strconv.FormatInt(int64((wa.length.ActualOption+1)*5), 10) + ":" +
		"buyStrip " +
		strconv.FormatFloat(float64(wa.buyStrip.ActualOption)/2.0, 'f', 2, 32) + "%:" +
		"sellStrip " +
		strconv.FormatFloat(float64(wa.sellStrip.ActualOption)/2.0, 'f', 2, 32) + "%"
}

func (wa *WmovingAverage) Trace(nick string) string {
	data := wa.data[nick]
	closes, _ := json.Marshal(data.closes)
	return "  Closes: " + string(closes) + "\n" +
		"  Totals: " +
		strconv.FormatFloat(float64(data.total), 'f', 2, 32) + "\n" +
		"  Numerators: " +
		strconv.FormatFloat(float64(data.numerator), 'f', 2, 32)
}

func (wa *WmovingAverage) Serialize() map[string]interface{} {
	s := wa.base.Serialize()
	s["length"] = wa.length.Serialize()
	s["buyStrip"] = wa.buyStrip.Serialize()
	s["sellStrip"] = wa.sellStrip.Serialize()
	return s
}

func RestoreWmovingAverage(
	base *Base, s map[string]interface{}) *WmovingAverage {

	return &WmovingAverage{
		base: base,

		length:    gen.Restore(s["length"].([]interface{})),
		buyStrip:  gen.Restore(s["buyStrip"].([]interface{})),
		sellStrip: gen.Restore(s["sellStrip"].([]interface{})),

		data: make(map[string]*wmovingAverageData),
	}
}

func MutateWmovingAverage(
	base *Base, f *WmovingAverage, mrange int) *WmovingAverage {

	return &WmovingAverage{
		base: base,

		length:    f.length.Mutate(mrange),
		buyStrip:  f.buyStrip.Mutate(mrange),
		sellStrip: f.sellStrip.Mutate(mrange),

		data: make(map[string]*wmovingAverageData),
	}
}

func (wa *WmovingAverage) Process(qs map[string]quote.Quote) {
	b := wa.base
	lg := uint((wa.length.ActualOption + 1) * 5)
	denominator := float32(lg * (lg + 1) / 2)
	bet := float64((b.Bet().ActualOption + 1) * 5000)
	ibex := b.Ibex().ActualOption
	canBuy := bet <= b.Cash()
	bs := make([]*order.Buy, 0)
	ss := make([]*order.Sell, 0)
	for k, v := range qs {
		if ibex == 0 {
			if _, ok := Ibex[k]; ok {
				continue
			}
		} else if ibex == 1 {
			if _, ok := Ibex[k]; !ok {
				continue
			}
		}

		close := v[quote.Close].(float32)
		if close < 0 {
			continue
		}

		if data, ok := wa.data[k]; ok {
			closes := append(data.closes, close)
			if uint(len(closes)) >= lg {
				if uint(len(closes)) == lg {
					total := float32(0)
					numerator := float32(0)
					for i, v := range closes {
						total += v
						numerator += (v * float32(i+1))
					}
					data.total = total
					data.numerator = numerator
          data.closes = closes
				} else {
					data.numerator = data.numerator - data.total + float32(lg)*close
					data.total = data.total - closes[0] + close
					data.closes = closes[1:]
				}

				prevCanBuy := data.canBuy
				prevCanSell := data.canSell
				avg := data.numerator / denominator
				buyCross := avg * (1 + float32(wa.buyStrip.ActualOption)*0.005)
				sellCross := avg * (1 - float32(wa.sellStrip.ActualOption)*0.005)
				data.canBuy = close <= buyCross
				data.canSell = close >= sellCross

				if canBuy && prevCanBuy && close > buyCross {
					bs = append(bs, order.NewBuy(k, bet))
				}
				if prevCanSell && close < sellCross {
					if num, ok := b.portfolio[k]; ok {
						ss = append(ss, order.NewSell(k, num))
					}
				}
			} else {
        data.closes = closes
      }
		} else {
			wa.data[k] = &wmovingAverageData{
				closes:  []float32{close},
				canBuy:  false,
				canSell: false,
				total:     0.0,
				numerator:     0.0,
			}
		}
	}
	b.buys = bs
	b.sells = ss
}
