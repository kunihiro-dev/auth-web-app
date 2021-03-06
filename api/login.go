package api

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/kunihiro-dev/auth-web-app/model/entity"
	"github.com/kunihiro-dev/auth-web-app/session"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, r)
		return
	}

	name := r.PostFormValue("name")
	password := r.PostFormValue("password")

	if password == "" || name == "" {
		Error(w, r)
		return
	}

	_, err := generatePassword(password, 32)
	if err != nil {
		fmt.Println("Password generate error.")
		Error(w, r)
		return
	}

	u := entity.UserInfo{Name: name, Password: password}

	n, err := session.Generate()
	if err != nil {
		Error(w, r)
		return
	}

	cookie := http.Cookie {Name: SESSION_KEY, Value: n, HttpOnly: true}
	http.SetCookie(w, &cookie)
	session.Add(u.Name, n)

	http.Redirect(w, r, "/top", http.StatusTemporaryRedirect)
}

func generatePassword(password string, saltLen int) ([32]byte, error) {
	salt := make([]byte, saltLen)
	var result [32]byte
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return result, err
	}

	passbytes := append(salt, []byte(password)...)
	result = sha256.Sum256(passbytes)
	return result, nil
}
