package main

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret"))

const MainSession = "session"

type Context struct {
	User    *User
	Session *sessions.Session

	w http.ResponseWriter
	r *http.Request
}

func (c *Context) AddFlash(s string) {
	c.Session.AddFlash(s)
}

func (c *Context) GetFlashes() []interface{} {
	return c.Session.Flashes()
}

func (c *Context) Close() {
	c.Session.Save(c.r, c.w)
}

func GetContext(w http.ResponseWriter, r *http.Request) Context {

	res := Context{w: w, r: r}

	// Get session
	sess, err := store.Get(r, MainSession)
	if err != nil {
		log.Fatal(err)
	}
	res.Session = sess

	// Get logged user
	id, ok := sess.Values["User.ID"]
	if ok {
		u := User{}
		db.First(&u, id.(int))
		res.User = &u
	}

	return res
}
