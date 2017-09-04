// Copyright 23-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package main

import(
 "github.com/dedeme/go/libcgi"
)

const (
  dataPath = "/deme/wwwcgi/dmcgi/Hconta"
)

func readRq(rq map[string]interface{}, field, context string) interface{} {
  r, ok := rq[field]
  if !ok {
    libcgi.Err("Field " + field + " is missing in" + context)
  }
  return r
}

func main () {
  libcgi.Init(dataPath)

  rq := libcgi.ReadRequest()

  rp := make(map[string]interface{})
  page := readRq(rq, libcgi.Page, "main").(string)
  switch page {
    case libcgi.PageConnection:
      sessionId := readRq(rq, libcgi.SessionId, page).(string)
      pageId := readRq(rq, libcgi.PageId, page).(string)
      rp[libcgi.SessionOk] = libcgi.Connect(sessionId, pageId)
      libcgi.Ok(rp)
    default:
      libcgi.Err("Page " + page + " is unknown")
  }

}
