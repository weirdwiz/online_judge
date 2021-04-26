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
	AddBatch(ub model.Batch, teacherEmail string) (bool, error)
	GetBatches(u model.User) ([]model.Batch, error)
	GetUser(email string) (model.User, error)
	AddQuestion(ud model.Assignment, teacherEmail string) (bool, error)
}

// Struct to handle the DB Connection
type DBClient struct {
	client   *bolt.DB
	filepath string
}

const (
	usersBucketName       string = "Users"
	studentListBucketName string = "Students"
	teacherListBucketName string = "Teachers"
	batchBucketName       string = "Batches"
	questionBucketName    string = "Questions"
)

var bucketList []string = []string{usersBucketName, studentListBucketName, teacherListBucketName, batchBucketName}

func (db *DBClient) Initialize(filepath string) {
	db.filepath = filepath
	db.Open()
	defer db.Close()
	err := db.client.Update(func(txn *bolt.Tx) error {
		// Initialize Bucket
		for _, bucket := range bucketList {
			_, err := txn.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
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

func (db *DBClient) CreateUser(u model.User) (bool, error) {
	isNew, err := db.CheckUserIsNew(u.Email)
	if err != nil || !isNew {
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
		s := model.Student{
			Email: u.Email,
			Name:  u.Name,
		}

		err = db.client.Update(func(txn *bolt.Tx) error {
			b := txn.Bucket([]byte(studentListBucketName))
			id, err := b.NextSequence()
			if err != nil {
				return err
			}
			s.ID = strconv.Itoa(int(id))

			studentBytes, err := json.Marshal(s)
			if err != nil {
				return err
			}

			err = b.Put([]byte(u.Email), studentBytes)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return false, nil
		}

	case "teacher":
		t := model.Teacher{
			Email: u.Email,
			Name:  u.Name,
		}

		err := db.client.Update(func(txn *bolt.Tx) error {
			b := txn.Bucket([]byte(teacherListBucketName))
			id, err := b.NextSequence()
			if err != nil {
				return err
			}
			t.ID = strconv.Itoa(int(id))

			teacherBytes, err := json.Marshal(t)
			if err != nil {
				return err
			}

			err = b.Put([]byte(u.Email), teacherBytes)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return false, nil
		}
	}

	return true, nil
}

func (db *DBClient) GetUser(email string) (model.User, error) {
	db.Open()
	defer db.Close()
	user := model.User{}
	err := db.client.Update(func(txn *bolt.Tx) error {
		user, _ = getUser(txn, email)
		return nil
	})
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func getUser(txn *bolt.Tx, email string) (model.User, error) {
	b := txn.Bucket([]byte(usersBucketName))
	userBytes := b.Get([]byte(email))
	u := model.User{}
	err := json.Unmarshal(userBytes, &u)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (db *DBClient) CheckUserIsNew(email string) (bool, error) {
	db.Open()
	defer db.Close()
	err := db.client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucketName))
		userBytes := b.Get([]byte(email))
		if userBytes != nil {
			return fmt.Errorf("user with email %s already exists", email)
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

	err := db.client.View(func(txn *bolt.Tx) error {
		u, err := getUser(txn, email)
		if err != nil {
			return err
		}

		err = u.CheckPassword(password)

		if err != nil {
			log.Printf("ERROR: %s", err)
			log.Printf("Invalid Login Attempt for User: %s", email)
			return fmt.Errorf("invalid email/password combnation")
		}
		return nil
	})
	if err != nil {
		return "nope", err
	}
	return "yep", nil
}

func (db *DBClient) GetBatches(u model.User) ([]model.Batch, error) {
	db.Open()
	defer db.Close()

	var res []model.Batch

	err := db.client.Update(func(txn *bolt.Tx) error {
		var batchList []string
		switch u.AccountType {
		case "student":
			s, err := getStudent(txn, u.Email)
			if err != nil {
				return err
			}
			batchList = s.Batches
		case "teacher":
			t, err := getTeacher(txn, u.Email)
			if err != nil {
				return err
			}
			batchList = t.Batches
		}
		for _, batchID := range batchList {
			batch, err := FetchBatch(txn, batchID)
			if err != nil {
				return err
			}
			res = append(res, batch)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FetchBatch(txn *bolt.Tx, ID string) (model.Batch, error) {
	b := txn.Bucket([]byte(batchBucketName))
	batchByte := b.Get([]byte(ID))
	batch := model.Batch{}

	err := json.Unmarshal(batchByte, &batch)
	if err != nil {
		return model.Batch{}, err
	}
	return batch, nil
}

func getStudent(txn *bolt.Tx, email string) (model.Student, error) {
	b := txn.Bucket([]byte(studentListBucketName))
	studentByte := b.Get([]byte(email))
	s := model.Student{}

	err := json.Unmarshal(studentByte, &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func getTeacher(txn *bolt.Tx, email string) (model.Teacher, error) {
	b := txn.Bucket([]byte(teacherListBucketName))
	teacherByte := b.Get([]byte(email))
	t := model.Teacher{}

	err := json.Unmarshal(teacherByte, &t)
	if err != nil {
		return t, err
	}

	return t, nil
}

func (db *DBClient) AddBatch(ub model.Batch, teacherEmail string) (bool, error) {
	db.Open()
	defer db.Close()

	err := db.client.Update(func(txn *bolt.Tx) error {
		b := txn.Bucket([]byte(batchBucketName))
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		ub.ID = strconv.Itoa(int(id))
		batchBytes, err := json.Marshal(ub)
		if err != nil {
			return err
		}
		err = b.Put([]byte(ub.ID), batchBytes)
		if err != nil {
			return err
		}

		b = txn.Bucket([]byte(studentListBucketName))
		for _, student := range ub.Students {
			s, err := getStudent(txn, student)
			if err != nil {
				return err
			}

			s.Batches = append(s.Batches, ub.ID)

			studentByte, err := json.Marshal(s)
			if err != nil {
				return err
			}

			err = b.Put([]byte(s.Email), studentByte)
			if err != nil {
				return err
			}
		}

		t, err := getTeacher(txn, teacherEmail)
		if err != nil {
			return err
		}

		t.Batches = append(t.Batches, ub.ID)

		teacherByte, err := json.Marshal(t)
		if err != nil {
			return err
		}

		err = b.Put([]byte(t.Email), teacherByte)
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

func (db *DBClient) AddAssignment(ud model.Assignment, teacherEmail string) (bool, error) {
	db.Open()
	defer db.Close()

	err := db.client.Update(func(txn *bolt.Tx) error {
		b := txn.Bucket([]byte(questionBucketName))
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		ud.ID = strconv.Itoa(int(id))
		assignmentBytes, err := json.Marshal(ud)
		if err != nil {
			return err
		}
		err = b.Put([]byte(ud.ID), assignmentBytes)
		if err != nil {
			return err
		}

			s.Batches = append(s.Batches, ud.ID)

			studentByte, err := json.Marshal(s)
			if err != nil {
				return err
			}

			err = b.Put([]byte(s.Email), studentByte)
			if err != nil {
				return err
			}
		}

		t, err := getTeacher(txn, teacherEmail)
		if err != nil {
			return err
		}

		t.Batches = append(t.Batches, ud.ID)

		teacherByte, err := json.Marshal(t)
		if err != nil {
			return err
		}

		err = b.Put([]byte(t.Email), teacherByte)
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

func AddAssignment(bID string, a model.Assignment) (bool, error) {

	return true, nil
}
