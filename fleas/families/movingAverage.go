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

type movingAverageData struct {
	closes  []float32
	canBuy  bool
	canSell bool
	sum     float32
}

type MovingAverage struct {
	base *Base

	length    *gen.Gen // (v + 1) * 5
	buyStrip  *gen.Gen // v * 0.005
	sellStrip *gen.Gen // v * 0.005

	data map[string]*movingAverageData
}

func (ma *MovingAverage) Base() *Base {
	return ma.base
}

func NewMovingAverage(base *Base) *MovingAverage {
	return &MovingAverage{
		base: base,

		length:    gen.NewRandom(30),
		buyStrip:  gen.NewRandom(11),
		sellStrip: gen.NewRandom(11),

		data: make(map[string]*movingAverageData),
	}
}

func (ma *MovingAverage) String() string {
	return ma.base.String() + ":" +
		"length " +
		strconv.FormatInt(int64((ma.length.ActualOption+1)*5), 10) + ":" +
		"buyStrip " +
		strconv.FormatFloat(float64(ma.buyStrip.ActualOption)/2.0, 'f', 2, 32) + "%:" +
		"sellStrip " +
		strconv.FormatFloat(float64(ma.sellStrip.ActualOption)/2.0, 'f', 2, 32) + "%"
}

func (ma *MovingAverage) Trace(nick string) string {
	data := ma.data[nick]
	closes, _ := json.Marshal(data.closes)
	return "  Closes: " + string(closes) + "\n" +
		"  Sum: " + strconv.FormatFloat(float64(data.sum), 'f', 2, 32)
}

func (ma *MovingAverage) Serialize() map[string]interface{} {
	s := ma.base.Serialize()
	s["length"] = ma.length.Serialize()
	s["buyStrip"] = ma.buyStrip.Serialize()
	s["sellStrip"] = ma.sellStrip.Serialize()
	return s
}

func RestoreMovingAverage(
	base *Base, s map[string]interface{}) *MovingAverage {

	return &MovingAverage{
		base: base,

		length:    gen.Restore(s["length"].([]interface{})),
		buyStrip:  gen.Restore(s["buyStrip"].([]interface{})),
		sellStrip: gen.Restore(s["sellStrip"].([]interface{})),

		data: make(map[string]*movingAverageData),
	}
}

func MutateMovingAverage(
	base *Base, f *MovingAverage, mrange int) *MovingAverage {

	return &MovingAverage{
		base: base,

		length:    f.length.Mutate(mrange),
		buyStrip:  f.buyStrip.Mutate(mrange),
		sellStrip: f.sellStrip.Mutate(mrange),

		data: make(map[string]*movingAverageData),
	}
}

func (ma *MovingAverage) Process(qs map[string]quote.Quote) {
	b := ma.base
	lg := uint((ma.length.ActualOption + 1) * 5)
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

		if data, ok := ma.data[k]; ok {
			closes := append(data.closes, close)
			if uint(len(closes)) >= lg {
				if uint(len(closes)) == lg {
					sum := float32(0)
					for _, v := range closes {
						sum += v
					}
					data.sum = sum
					data.closes = closes
				} else {
					data.sum = data.sum - closes[0] + close
					data.closes = closes[1:]
				}

				prevCanBuy := data.canBuy
				prevCanSell := data.canSell
				avg := data.sum / float32(lg)
				buyCross := avg * (1 + float32(ma.buyStrip.ActualOption)*0.005)
				sellCross := avg * (1 - float32(ma.sellStrip.ActualOption)*0.005)
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
			ma.data[k] = &movingAverageData{
				closes:  []float32{close},
				canBuy:  false,
				canSell: false,
				sum:     0.0,
			}
		}
	}
	b.buys = bs
	b.sells = ss
}
