package session
import (
  "gopkg.in/mgo.v2"
  "os"
  "bitbucket.org/tampajohn/gadget-arm/errors"
)
var (
  mgoSession *mgo.Session
)
func Get (connectionVariable string) *mgo.Session {
  if mgoSession == nil {
    cs := os.Getenv(connectionVariable)
    var err error
    mgoSession, err = mgo.Dial(cs)
    errors.Check(err)
  }
  return mgoSession.Clone()
}
