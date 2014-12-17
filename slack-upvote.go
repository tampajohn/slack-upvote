package main

import (
  "github.com/codegangsta/negroni"
  "github.com/gorilla/mux"
  "os"
  "fmt"
  "net/http"
  "gopkg.in/mgo.v2"
  _"gopkg.in/mgo.v2/bson"
  "strings"
)
type Team struct {
  Id string `bson:"_id"`
  Domain string `bson:"domain"`
  Father string `bson:"father"`
}
type Channel struct {
  Id string `bson:"_id"`
  Name string `bson:"name"`
}
type Mention struct {
  Id string `bson:"_id"`
  TeamId string `bson:"team_id"`
  Votes int64 `bson:"votes"`
}
func check(err error) {
  panic(err)
}
var (
  mgoSession *mgo.Session
)
func getSession () *mgo.Session {
  if mgoSession == nil {
    cs := os.Getenv("SLACKVOTE_DB")
    var err error
    mgoSession, err = mgo.Dial(cs)
    if err != nil {
      panic(err) // no, not really
    }
  }
  return mgoSession.Clone()
}
func UpDownHandler(rw http.ResponseWriter, r *http.Request) {
  r.ParseMultipartForm(5120)
  isCmd := len(r.PostForm["command"]) > 0
  isTrg := len(r.PostForm["trigger_word"]) > 0
  db := getSession().DB("slack-upvote")
  mentionId := r.PostForm["text"][0]
  sfx := ""

  if isTrg {
    mentionId = strings.Trim(mentionId, r.PostForm["trigger_word"][0])
  }
  m := Mention{
    Id: mentionId,
  }

  db.C("mentions").FindId(m.Id).One(&m)

  if isCmd {
    cmd := r.PostForm["command"][0]

    if cmd == "/up" {
      m.Votes++
    } else if cmd == "/down" {
      m.Votes--
    } else {
      rw.Write([]byte("Egad, how did you get here?!"))
      return
    }
  } else if isTrg {
    trg := r.PostForm["trigger_word"][0]
    if trg == "+" {
      m.Votes++
    } else if trg == "-" {
      m.Votes--
    }
  }
  db.C("mentions").UpsertId(m.Id, m)

  if m.Votes > 1 || m.Votes < -1 || m.Votes == 0 {
    sfx = "s"
  }
  text := fmt.Sprintf("'%s' has %v vote%s.", m.Id, m.Votes, sfx)
  if isCmd {
    rw.Write([]byte(text))
  } else {
    rw.Write([]byte(fmt.Sprintf("{\"text\":\"%s\"}", text)))
  }
}
func DownHandler(rw http.ResponseWriter, r *http.Request) {
  rw.Write([]byte("WOO"))
}
func main () {
  router := mux.NewRouter()
  router.HandleFunc("/up", UpDownHandler)
  router.HandleFunc("/down", UpDownHandler)

  n := negroni.New(
    negroni.NewLogger(),
    negroni.NewStatic(http.Dir("public")),
    negroni.Wrap(router),
  )
  n.Run(":"+os.Getenv("PORT"))
}