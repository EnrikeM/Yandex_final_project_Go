package httpsrv

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
	"github.com/golang-jwt/jwt/v5"
)

type PassRequest struct {
	Password string `json:"password"`
}

type Claims struct {
	Password string `json:"password"`
	jwt.RegisteredClaims
}

var secretKey = []byte("super-secret-key-for-test")

func GenerateJWT(password string) (string, error) {

	claims := Claims{
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (a *API) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err == nil {
			jwtToken := cookie.Value

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})

			if err == nil && token.Valid {
				next.ServeHTTP(w, r)
				return
			}

		}

		pass := a.config.TODO_PASSWORD
		if len(pass) > 0 {
			var passRequest PassRequest
			err := json.NewDecoder(r.Body).Decode(&passRequest)
			if err != nil {
				rErr := apierrors.New(err.Error())
				rErr.Error(w, http.StatusUnauthorized)
				return
			}

			if passRequest.Password == pass {
				token, err := GenerateJWT(passRequest.Password)
				if err != nil {
					rErr := apierrors.New(err.Error())
					rErr.Error(w, http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "token",
					Value:    token,
					Path:     "/",
					HttpOnly: true,
					Secure:   false,
				})

				WriteResponse("token", token, w, http.StatusOK)
				return

			}

			rErr := apierrors.New("invalid password")
			rErr.Error(w, http.StatusUnauthorized)
			return

		}

		next.ServeHTTP(w, r)
	})
}

func (a *API) signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	WriteResponse("success", "you have successfully signed in", w, http.StatusOK)
}
