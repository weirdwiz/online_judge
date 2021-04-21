package dbclient

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/weirdwiz/online_judge/authentication/internal/app/model"
)

// Interface for the DB Client
type IDBClient interface {
	Initialize(filepath string)
	Open()
	Close()
	CreateUser(u model.User) (bool, error)
	Login(email, password, accounttype string) (string, error)
	CheckUserIsNew(email string) (bool, error)
	AddBatch(batch, []students string) (string, error)
}

// Struct to handle the DB Connection
type DBClient struct {
	client   *bolt.DB
	filepath string
}

const (
	usersBucketName string = "Users"
	batchBucketName string = "Batch"
)

func (db *DBClient) Initialize(filepath string) {
	db.filepath = filepath
	db.Open()
	defer db.Close()
	err := db.client.Update(func(txn *bolt.Tx) error {
		// Initialize Users Bucket
		_, err := txn.CreateBucketIfNotExists([]byte(usersBucketName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (db *DBClient) Open() {
	if db.filepath == "" {
		log.Fatal("Filepath required for Database")
	}
	d, err := bolt.Open(db.filepath, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	db.client = d
}

func (db *DBClient) Close() {
	db.client.Close()
}

/* func (db *DBClient) CreateStudent(s model.Student) (bool, error) {

}
*/

func (db *DBClient) CreateUser(u model.User) (bool, error) {
	isNew, err := db.CheckUserIsNew(u.Email)
	if err != nil || isNew == false {
		return false, err
	}
	db.Open()
	defer db.Close()

	err = db.client.Update(func(txn *bolt.Tx) error {
		b := txn.Bucket([]byte(usersBucketName))
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		u.ID = strconv.Itoa(int(id))
		u.HashPassword()
		userBytes, err := json.Marshal(u)
		if err != nil {
			return err
		}
		err = b.Put([]byte(u.Email), userBytes)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, nil
	}

	switch u.AccountType {
	case "student":
		err = db.client.Update(func(txn *bolt.Tx) error {
			//b := txn.Bucket([]byte(studentList))
			return nil
		})
	case "teacher":
	}

	return true, nil
}

func (db *DBClient) CheckUserIsNew(email string) (bool, error) {
	db.Open()
	defer db.Close()
	err := db.client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucketName))
		userBytes := b.Get([]byte(email))
		if userBytes != nil {
			return fmt.Errorf("User with email %s already exists", email)
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DBClient) Login(email, password, accounttype string) (string, error) {
	db.Open()
	defer db.Close()
	err := db.client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucketName))
		userBytes := b.Get([]byte(email))
		u := model.User{}
		json.Unmarshal(userBytes, &u)
		err := u.CheckPassword(password)
		if err != nil {
			log.Printf("ERROR: %s", err)
			log.Printf("Invalid Login Attempt for User: %s", email)
			return fmt.Errorf("Invalid Email/Password Combnation")
		}
		return nil
	})
	if err != nil {
		return "nope", err
	}
	return "yep", nil
}


func (db *DBClient) AddBatch(ub model.Batch) (bool, error) {
	db.Open()
	defer db.Close()

	err = db.client.Update(func(txn *bolt.Tx) error {
		b := txn.Bucket([]byte(batchBucketName))
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		ub.ID = strconv.Itoa(int(id))
		ub.HashPassword()
		userBytes, err := json.Marshal(u)
		if err != nil {
			return err
		}
		err = b.Put([]byte(ub.Name), userBytes)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, nil
	}

	return true, nil
}
