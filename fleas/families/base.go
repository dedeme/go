// Copyright 6-Dec-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package families

import (
	"github.com/dedeme/go/fleas/db"
	"github.com/dedeme/go/fleas/gen"
	"github.com/dedeme/go/fleas/order"
	"strconv"
)

type Base struct {
	id         string
	cycle      uint64
	family     *gen.Gen
	bet        *gen.Gen
	ibex       *gen.Gen
	mutability *gen.Gen

	buys      []*order.Buy
	sells     []*order.Sell
	cash      float64
	portfolio map[string]uint64
	nbuys     uint
	nsells    uint
}

func NewBase(cycle uint64) *Base {
	return &Base{
		id:         db.NewFleaId(),
		cycle:      cycle,
		family:     gen.NewRandom(EndFamilies),
		ibex:       gen.NewRandom(3),
		bet:        gen.NewRandom(10),
		mutability: gen.NewRandom(40),

		buys:      make([]*order.Buy, 0),
		sells:     make([]*order.Sell, 0),
		cash:      InitialCash,
		portfolio: make(map[string]uint64),
		nbuys:     0,
		nsells:    0}

}

func (base *Base) Id() string {
	return base.id
}

func (base *Base) Cycle() uint64 {
	return base.cycle
}

func (base *Base) Family() *gen.Gen {
	return base.family
}

// Ibex returns a Gen with three values: 0 -> No Ibex, 1-> Only Ibex, 2 -> All
func (base *Base) Ibex() *gen.Gen {
	return base.ibex
}

func (base *Base) Bet() *gen.Gen {
	return base.bet
}

func (base *Base) Mutability() *gen.Gen {
	return base.mutability
}

func (base *Base) Buys() []*order.Buy {
	return base.buys
}

func (base *Base) Sells() []*order.Sell {
	return base.sells
}

func (base *Base) CleanOrders() {
	base.buys = make([]*order.Buy, 0)
	base.sells = make([]*order.Sell, 0)
}

func (base *Base) Cash() float64 {
	return base.cash
}

func (base *Base) AddCash(money float64) {
	base.cash += money
}

func (base *Base) Nbuys() uint {
	return base.nbuys
}

func (base *Base) IncNbuys() {
	base.nbuys++
}

func (base *Base) Nsells() uint {
	return base.nsells
}

func (base *Base) IncNsells() {
	base.nsells++
}

func (base *Base) Portfolio() map[string]uint64 {
	return base.portfolio
}

func (base *Base) String() string {
	return "id " + base.id + ":" +
		Names[base.family.ActualOption] + ":" +
		"cycle " +
		strconv.FormatUint(base.cycle, 10) + ":" +
		"bet " +
		strconv.FormatInt(int64((base.bet.ActualOption+1)*5000), 10) + ":" +
		"ibex " +
		strconv.FormatInt(int64(base.ibex.ActualOption), 10) + ":" +
		"mutability " +
		strconv.FormatInt(int64(base.mutability.ActualOption), 10)
}

func (base *Base) Serialize() map[string]interface{} {
	serial := make(map[string]interface{})
	serial["id"] = base.id
	serial["cycle"] = base.cycle
	serial["family"] = base.family.Serialize()
	serial["bet"] = base.bet.Serialize()
	serial["ibex"] = base.ibex.Serialize()
	serial["mutability"] = base.mutability.Serialize()
	return serial
}

func RestoreBase(s map[string]interface{}) *Base {
	return &Base{
		id:         s["id"].(string),
		cycle:      uint64(s["cycle"].(float64)),
		family:     gen.Restore(s["family"].([]interface{})),
		bet:        gen.Restore(s["bet"].([]interface{})),
		ibex:       gen.Restore(s["ibex"].([]interface{})),
		mutability: gen.Restore(s["mutability"].([]interface{})),

		buys:      make([]*order.Buy, 0),
		sells:     make([]*order.Sell, 0),
		cash:      InitialCash,
		portfolio: make(map[string]uint64),
		nbuys:     0,
		nsells:    0}
}

func MutateBase(base *Base, mrange int, cycle uint64) *Base {
	return &Base{
		id:         db.NewFleaId(),
		cycle:      cycle,
		family:     base.family.Mutate(mrange),
		bet:        base.bet.Mutate(mrange),
		ibex:       base.ibex.Mutate(mrange),
		mutability: base.mutability.Mutate(mrange),

		buys:      make([]*order.Buy, 0),
		sells:     make([]*order.Sell, 0),
		cash:      InitialCash,
		portfolio: make(map[string]uint64),
		nbuys:     0,
		nsells:    0}
}
