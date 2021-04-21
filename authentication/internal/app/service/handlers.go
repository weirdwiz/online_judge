package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/weirdwiz/online_judge/authentication/internal/app/dbclient"
	"github.com/weirdwiz/online_judge/authentication/internal/app/model"
)

var mySigningKey = []byte("signingKey")

type JWTstring struct {
	Token string `json:"token"`
}

func GenerateJWT(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error")
				}
				return mySigningKey, nil
			})
			if err != nil {
				WriteError(w, http.StatusUnauthorized, err)
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			WriteError(w, http.StatusUnauthorized, nil)
		}
	})
}

func HandleRegistration(w http.ResponseWriter, r *http.Request) {
	var user model.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("error decoding user"))
	}
	success, err := DBClient.CreateUser(user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Fprintf(w, "Status: %t", success)
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, err.Error())
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if r.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			WriteError(w, http.StatusBadRequest, fmt.Errorf("error decoding user"))
		}
	} else {
		email := r.FormValue("email")
		password := r.FormValue("password")
		accounttype := r.FormValue("accounttype")

		user.Email = email
		user.Password = password
		user.AccountType = accounttype
	}

	_, err := DBClient.Login(user.Email, user.Password, user.AccountType)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	tokenString, err := GenerateJWT(user.Email)
	if err != nil {
		fmt.Fprintf(w, "Error generating token string")
	}
	token := &JWTstring{
		Token: tokenString,
	}
	data, _ := json.Marshal(token)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

<<<<<<< HEAD
func HandleAddBatch(w http.ResponseWriter, r *http.Request) {

	tokenClaims, valid := extractClaims(r.Header.Get("Token"))
	if !valid {
		WriteError(w, http.StatusUnauthorized, nil)
		return
	}

	teacherEmail := fmt.Sprintf("%v", tokenClaims["email"])

=======

func HandleAddBatch(w http.ResponseWriter, r *http.Request) {
>>>>>>> 905bc4d (add Batch)
	var batch model.Batch
	if r.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&batch)
		if err != nil {
			WriteError(w, http.StatusBadRequest, fmt.Errorf("Error Decoding Batch"))
		}
<<<<<<< HEAD
	}

	_, err := DBClient.AddBatch(batch, teacherEmail)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HandleGetBatches(w http.ResponseWriter, r *http.Request) {
	tokenClaims, valid := extractClaims(r.Header.Get("Token"))
	if !valid {
		WriteError(w, http.StatusUnauthorized, nil)
		return
	}

	email := fmt.Sprintf("%v", tokenClaims["email"])

	user, err := DBClient.GetUser(email)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
	}

	batches, _ := DBClient.GetBatches(user)

	data, _ := json.Marshal(batches)
=======
	} else {
		name := r.FormValue("name")
		students := r.FormValue("students")

		batch.Name = name
		batch.Students = students
	}

	_, err := DBClient.AddBatch(batch.Name, batch.Students)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
>>>>>>> 905bc4d (add Batch)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
<<<<<<< HEAD
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := "signingKey"
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
=======

>>>>>>> 905bc4d (add Batch)
}

var DBClient dbclient.IDBClient
