// Copyright 01-Dic-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package db

import (
	"encoding/json"
  "strconv"
)

type db struct {
	newFlea uint64
	cycle   uint64
	newBest uint64
	bests   map[uint64][]string
	trace   []string
	fleas   map[string][]string // Array with two string: data and stat
}

var data *db = nil

func Serialize() string {
	serial := make(map[string]interface{})
	serial["newFlea"] = data.newFlea
	serial["cycle"] = data.cycle
	serial["newBest"] = data.newBest
	serial["bests"] = data.bests
	serial["trace"] = data.trace
	serial["fleas"] = data.fleas

	s, err := json.Marshal(serial)
	if err != nil {
		panic("db.Serialize: " + err.Error())
	}
	return string(s)
}

func Restore(serial string) {
	if serial == "" {
		data = &db{
			newFlea: 1,
			cycle:   0,
			newBest: 0,
			bests:   make(map[uint64][]string, 0),
			trace:   make([]string, 0),
			fleas:   make(map[string][]string),
		}
		return
	}

	var r map[string]interface{}
	err := json.Unmarshal([]byte(serial), &r)
	if err != nil {
		panic("db.Restore: " + err.Error())
	}

  bests := make(map[uint64][]string)
  for k, v := range r["bests"].(map[string]interface{}) {
    i, _ := strconv.ParseUint(k, 10, 64)
    vis := v.([]interface{})
    vs := make([]string, len(vis))
    for ix, vi := range vis {
      vs[ix] = vi.(string)
    }
    bests[i] = vs
  }

  traceIs := r["trace"].([]interface{})
  trace := make([]string, len(traceIs))
  for ix, traceI := range traceIs {
    trace[ix] = traceI.(string)
  }

  fleas:= make(map[string][]string)
  for k, v := range r["fleas"].(map[string]interface{}) {
    vis := v.([]interface{})
    vs := make([]string, len(vis))
    for ix, vi := range vis {
      vs[ix] = vi.(string)
    }
    fleas[k] = vs
  }

	data = &db{
		newFlea: uint64(r["newFlea"].(float64)),
		cycle:   uint64(r["cycle"].(float64)),
		newBest: uint64(r["newBest"].(float64)),
		bests:   bests,
		trace:   trace,
		fleas:   fleas,
	}
}

func GetCycle() uint64 {
	return data.cycle
}

func IncCycle() {
	data.cycle++
}

func ExistsFlea(f string) bool {
	if data == nil {
		return false
	}
	if _, ok := data.fleas[f]; ok {
		return true
	}
	return false
}

// GetFleas returns a map which values are Arrays with two string: data and stat
func GetFleas() map[string][]string {
	return data.fleas
}

// New FleasId returns a new id and increments the id counter.
func NewFleaId() string {
	id := strconv.FormatUint(data.newFlea, 10)
  data.newFlea++
	return id
}


func FleaSerial(fs map[string]interface{}) string {
	s, err := json.Marshal(fs)
	if err != nil {
		panic("db.FleaSerial: " + err.Error())
	}
	return string(s)
}

func FleaRestore(serial string) map[string]interface{} {
	var r map[string]interface{}
	err := json.Unmarshal([]byte(serial), &r)
	if err != nil {
		panic("db.FleaRestore: " + err.Error())
	}
	return r
}

// Writes the 100 first fleas data in format: FleaData|cash:buy:sell
func WriteBests(bestsData []string) {
	bests := data.bests
	ix := data.newBest
	bests[ix] = bestsData
	data.newBest++
	for k, _ := range bests {
		dif := ix - k
		if (dif > 10000 && k%10000 != 0) ||
			(dif > 10000 && k%10000 != 0) ||
			(dif > 100 && k%100 != 0) ||
			(dif > 10 && k%10 != 0) {

			delete(bests, k)
		}
	}
}

func WriteTrace(trace []string) {
	data.trace = trace
}
