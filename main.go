package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis"
)

var rclnt = redis.NewClient(&redis.Options{
	Addr:     redisDB(),
	Password: "",
	DB:       0,
})

type user struct {
	Email string
	Fname string
	Lname string
	Role  string
	Pword []byte
}

func (u *user) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *user) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	return nil
}

var tmpl *template.Template
var dbSessions = map[string]string{}

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	///serve css
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/admin", adminHome)
	http.HandleFunc("/admin/cumulus", cumulusPost)
	http.HandleFunc("/userhome", userHome)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	u := getUser(w, r)
	tmpl.ExecuteTemplate(w, "index.html", u)
}

func login(w http.ResponseWriter, r *http.Request) {
	if loggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		pword := r.FormValue("password")

		var u user
		res, err := rclnt.Get(email).Result()
		err = u.UnmarshalBinary([]byte(res))
		if err != nil {
			http.Error(w, "Cannot Retreive User, DB UnMarshal Error!", http.StatusInternalServerError)
		}

		if email != u.Email {
			http.Error(w, "Username/Password Error", http.StatusForbidden)
			return
		}
		err = bcrypt.CompareHashAndPassword(u.Pword, []byte(pword))
		if err != nil {
			http.Error(w, "Username/Password Error", http.StatusForbidden)
			return
		}
		//Username and Password match. Create a session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = email

		if u.Role == "admin" {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		} else if u.Role == "user" {
			http.Redirect(w, r, "/user", http.StatusSeeOther)
			return
		}
	}
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if !loggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	c, _ := r.Cookie("session")
	//remove session from db
	delete(dbSessions, c.Value)
	//delete cookie by assigning new cookie with negative age value
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if loggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var u user
	//process form to create user
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		fname := r.FormValue("firstname")
		lname := r.FormValue("lastname")
		role := r.FormValue("role")
		pword := r.FormValue("password")
		//Check if username is already in the db
		_, err := rclnt.Get(email).Result()
		if err != redis.Nil {
			fmt.Println(err)
			http.Error(w, "Username Taken", http.StatusForbidden)
		}
		//Create a session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = email
		//Store user info. First, encrpt pword
		bs, err := bcrypt.GenerateFromPassword([]byte(pword), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Server Error, Password Related", http.StatusInternalServerError)
			return
		}
		u = user{email, fname, lname, role, bs}
		dbSessions[c.Value] = email

		ru, err := u.MarshalBinary()
		if err != nil {
			fmt.Println("Marshal Err:", err)
		}

		err = rclnt.Set(email, ru, 0).Err()
		if err != nil {
			http.Error(w, "Database SET error", http.StatusInternalServerError)
		}
		if u.Role == "admin" {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		} else if u.Role == "user" {
			http.Redirect(w, r, "/user", http.StatusSeeOther)
			return
		}
	}

	tmpl.ExecuteTemplate(w, "signup.html", u)

}

func adminHome(w http.ResponseWriter, r *http.Request) {
	u := getUser(w, r)
	if !loggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if u.Role != "admin" {
		http.Error(w, "Only Admins Allowed", http.StatusForbidden)
		return
	}
	tmpl.ExecuteTemplate(w, "admin.html", u)
}

func cumulusPost(w http.ResponseWriter, r *http.Request) {
	u := getUser(w, r)
	if u.Role != "admin" {
		http.Error(w, "Only Admins Allowed", http.StatusForbidden)
		return
	}

	ipAddr := r.FormValue("ipaddr")
	cmd := r.FormValue("cmd")
	uname := r.FormValue("uname")
	pword := r.FormValue("pword")
	cmlsRes := cumulusAction(u, ipAddr, cmd, uname, pword)
	fmt.Fprint(w, string(cmlsRes))
}

func userHome(w http.ResponseWriter, r *http.Request) {
	u := getUser(w, r)
	if !loggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if u.Role != "user" {
		http.Error(w, "Users Homepage", http.StatusForbidden)
		return
	}
	tmpl.ExecuteTemplate(w, "user.html", u)
}

func getUser(w http.ResponseWriter, r *http.Request) user {
	//Check for cookie
	c, err := r.Cookie("session")
	if err != nil {
		sID, _ := uuid.NewV4() //throw away error.
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
	}
	//At this point, the user has a cookie
	var u user
	if email, ok := dbSessions[c.Value]; ok {
		res, err := rclnt.Get(email).Result()
		err = u.UnmarshalBinary([]byte(res))
		if err != nil {
			http.Error(w, "Cannot Retreive User, DB UnMarshal Error!", http.StatusInternalServerError)
		}

	}
	return u
}

func loggedIn(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	email := dbSessions[c.Value]
	_, err = rclnt.Get(email).Result()
	if err != redis.Nil {
		return true
	}
	return false
}

func cumulusAction(u user, ip, cmd, uname, pword string) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	url := "https://" + ip + ":8080/nclu/v1/rpc/"

	var jsonStr []byte
	switch cmd {
	case "Interface Stats":
		jsonStr = []byte(`{"cmd": "show counters json"}`)
	case "Version":
		jsonStr = []byte(`{"cmd": "show system"}`)
	case "Interface Config":
		jsonStr = []byte(`{"cmd": "show interface json"}`)
	case "Configuration":
		jsonStr = []byte(`{"cmd": "show configuration"}`)
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.SetBasicAuth(uname, pword)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body
}

// dynamic listening port
func redisDB() string {
	host := os.Getenv("REDIS_HN")
	if len(host) == 0 {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PT")
	if len(port) == 0 {
		port = "6379"
	}

	return host + ":" + port
}
