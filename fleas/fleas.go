// Copyright 27-Nov-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
	"github.com/dedeme/go/fleas/db"
	"github.com/dedeme/go/fleas/families"
	"github.com/dedeme/go/fleas/fio"
	"github.com/dedeme/go/fleas/flea"
	"github.com/dedeme/go/fleas/order"
	"github.com/dedeme/go/fleas/quote"
	"github.com/dedeme/go/fleas/stat"
	"os"
	"strconv"
	"time"
)

func help() {
	fio.WriteMsg("Usage:\n" +
		"  fleas\n" +
		"  fleas stop\n" +
		"  fleas force (remove stop lock)\n" +
		"  fleas reset\n" +
		"  fleas backup\n" +
		"  fleas restore file -> (e.g. fleas restore Fleas20171201)\n" +
		"  fleas trace <flea> -> (e.g. fleas trace 12543)\n" +
		"  fleas time <minutes> -> (e.g. fleas 5)\n" +
		"  or\n" +
		"  fleas help")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("main.Recover:\n%v", r)
			fio.WriteMsg(msg)
			os.Exit(0)
		}
	}()

	nArgs := len(os.Args)

	if nArgs == 2 {
		arg := os.Args[1]
		if arg == "force" {
			fio.Force()
			fio.WriteMsg("start.lock removed")
			return
		}
		if arg == "stop" {
			fio.Stop()
			fio.WriteMsg("Stopping...")
			return
		}
	}

	if fio.IsStarted() {
		panic("Fleas is started")
	}

	if nArgs > 3 {
		help()
		return
	}

	if nArgs > 2 {
		arg := os.Args[1]
		if arg == "restore" {
			fio.Restore(os.Args[2])
			fio.WriteMsg("Restore done")
			return
		}
		if arg == "trace" {
			flea := os.Args[2]
			if !db.ExistsFlea(flea) {
				panic("Flea " + flea + " does not exist")
			}
			families.Trace = flea
		} else {
			help()
			return
		}
	}

	minutes := 1
	if nArgs == 2 {
		arg := os.Args[1]
		if arg == "reset" {
			fio.Reset()
			fio.WriteMsg("Fleas reset")
			return
		}
		if arg == "backup" {
			fio.Backup()
			fio.WriteMsg("Backup done")
			return
		}

		its, err := strconv.ParseUint(arg, 10, 32)
		if err != nil {
			help()
			return
		}
		minutes = int(its)
	}

	fio.Ini()
	fio.Start()
	db.Restore(fio.ReadDb())

	nicks, ibex := fio.ListNicks()
	families.Ini(ibex)

	quotes := fio.Quotes(nicks)
	t1 := time.Now()
	maxCycles := 1000
	if families.Trace != "" {
		fio.WriteMsg("Tracing " + families.Trace)
		maxCycles = 1
	}
	for i := 0; i < maxCycles; i++ {
		cycle := db.GetCycle()
		fio.WriteMsg(fmt.Sprintf("%d (%d)", cycle, i+1))

		fleaData := db.GetFleas()
		for len(fleaData) < flea.Number {
			f := flea.New(cycle - 1)
			fleaData[flea.Id(f)] =
				[]string{db.FleaSerial(f.Serialize()), stat.New().String()}
		}

		bests := flea.MkBests()
		deleted := 0
		for id, data := range fleaData {
      f := flea.Restore(db.FleaRestore(data[0]))
			fSt := stat.FromString(data[1])

			if id == families.Trace {
				families.TraceTx = make([]string, 0)
				families.TraceTx = append(families.TraceTx,
					f.String()+"|"+data[1])
				families.TraceTx = append(families.TraceTx,
					"---------------------------")
			}

			for _, day := range quotes {
				flea.Process(f, day)
			}

			cash := flea.Cash(f)
			for k, v := range flea.Portfolio(f) {
				q := float32(-1)
				for i := len(quotes) - 1; i >= 0; i-- {
					q0 := quotes[i][k][quote.Close].(float32)
					if q0 > 0 {
						q = q0
						break
					}
				}
				if q < 0 {
					panic(k + ": Wrong quotes when closing")
				}
				cash += order.Income(v, q)
			}

			if cash > families.InitialCash {
				fSt.Setup(
					cycle-flea.Cycle(f), cash, flea.Nbuys(f), flea.Nsells(f))

				data[1] = fSt.String()
				flea.AddBest(bests, cycle, f, fSt.Cash, uint(fSt.Buys), uint(fSt.Sells))
			} else {
				delete(fleaData, id)
				deleted++
			}
		}

		muted := make([]flea.BestData, 0)
		for _, best := range bests {
			//      if flea.Cycle(best.F) + flea.MutationCycle <= cycle {
			muted = append(muted, best)
			//      }
		}
		mutedLen := len(muted)
		if mutedLen > 0 {
			d := 0
			for {
				for _, st := range muted {
					f := flea.Mutate(st.F, cycle)
          fleaData[flea.Id(f)] =
            []string{db.FleaSerial(f.Serialize()), stat.New().String()}
					d++
					if d >= deleted {
						break
					}
				}
				if d >= deleted {
					break
				}
			}
		}

		bestsData := make([]string, 0)
		for i := 0; i < len(bests); i++ {
			st := bests[i]
			bestsData = append(bestsData, st.F.String()+
				"|"+
				strconv.FormatFloat(st.Cash, 'f', 2, 64)+":"+
				strconv.FormatUint(uint64(st.Buys), 10)+":"+
				strconv.FormatUint(uint64(st.Sells), 10))
		}
		db.WriteBests(bestsData[0:])

		db.IncCycle()
    fio.WriteDb(db.Serialize())
		t2 := time.Now()
		min := t2.Sub(t1).Minutes()
		fio.WriteMsg(fmt.Sprintf("%.2f minutes", min))
		if min > float64(minutes) {
			break
		}
		if fio.IsStoped() {
			fio.WriteMsg("Stopped")
			break
		}
	}
	if families.Trace != "" {
		db.WriteTrace(families.TraceTx)
		fio.WriteMsg("End Trace")
	}
	fio.Force()
}
