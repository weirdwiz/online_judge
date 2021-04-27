package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	AccountType string `json:"accounttype"`
}

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Result bool   `json:"result"`
}

type Assignment struct {
	ID        string `json:"id"`
	Question  string `json:"question"`
	TestCases []TestCase
}

type Batch struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Students    []string `json:"students"`
	Assignments []string `json:"assignments"`
}

type Submission struct {
	ID           string     `json:"id"`
	AssignmentID string     `json:"assignmentId`
	Code         string     `json:"code"`
	Language     string     `json:"lang"`
	Result       []TestCase `json:"result"`
}

type Student struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Batches     []string `json:"batches"`
	Submissions []string `json:"submissions"`
}

type Teacher struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Batches []string `json:"classes"`
}

func (u *User) HashPassword() error {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}
	u.Password = string(password)
	return nil
}

func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
