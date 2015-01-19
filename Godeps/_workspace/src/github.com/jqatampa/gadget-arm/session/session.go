package session

import (
	"os"
	"strings"

	"github.com/jqatampa/gadget-arm/errors"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

func Get(connectionVariable string) *mgo.Session {
	if session == nil {

		var cs string

		if strings.HasPrefix(connectionVariable, "mongodb://") {
			cs = connectionVariable
		} else {
			cs = os.Getenv(connectionVariable)
		}

		var err error

		session, err = mgo.Dial(cs)

		errors.Check(err)

		// http://godoc.org/labix.org/v2/mgo#Session.SetMode
		session.SetMode(mgo.Monotonic, true)
	}
	return session.Copy()
}
