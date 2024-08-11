package main

import (
	"fmt"
	"html/template"
	"mymod/pkg/postgres"
	"net/http"

	"github.com/gorilla/sessions"
	_ "golang.org/x/crypto/bcrypt"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key          = []byte("super-secret-key")   // in env or hold in config
	store        = sessions.NewCookieStore(key) // ключ для шифрования сессий
	account_temp = make(map[string]string)
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "not path", http.StatusNotFound)
		return
	}

	session, _ := store.Get(r, "cookie-name")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "not auth or not login", http.StatusForbidden)
		return // нельзя ему продолжать
	}
	tmp, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	tmp, err := template.ParseFiles("./ui/html/login.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	//authentication
	if r.Method == http.MethodPost {
		// Authentication goes here
		temp_login := r.FormValue("login")
		temp_passw := r.FormValue("password")

		//examination avalible
		check_account := postgres.CheckAvailibleUsers(temp_login, temp_passw)
		if check_account {
			// Set user as authenticated
			session.Values["authenticated"] = true
			err := session.Save(r, w)
			if err != nil {
				http.Error(w, "Error saving session", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	tmp.Execute(w, nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	// Set user as authenticated
	session.Values["authenticated"] = false
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func registration(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		templogin := r.FormValue("login")
		temppassw := r.FormValue("password")
		tempconfirmpassw := r.FormValue("password_confirm")
		if temppassw == tempconfirmpassw {
			err := postgres.InsertDb(templogin, temppassw)
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	tmp, err := template.ParseFiles("./ui/html/registration.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, nil)
}

func main() {
	account_temp["wefd"] = "qwert"
	mux := http.NewServeMux()

	mux.HandleFunc("/", MainPage)
	mux.HandleFunc("/registration", registration)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)

	http.ListenAndServe(":8080", mux)
}
