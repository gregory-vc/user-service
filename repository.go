package main

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/couchbase/gocb.v1"
)

const primaryKey = "user_%d"
const table = "user"

type Repository interface {
	GetAll() ([]*User, error)
	GetByIDs(ids []uint64) ([]*User, error)
	Get(id uint64) (*User, error)
	Delete(id uint64) (*User, error)
	Create(user *User) (*User, error)
	Update(userUpdate *User) (*User, error)
	GetByEmail(email string) (*User, error)
}

type UserRepository struct {
	bucket *gocb.Bucket
}

func (repo *UserRepository) Delete(id uint64) (*User, error) {
	user := &User{}
	var users []*User
	queryStr := fmt.Sprintf("DELETE FROM `%s` "+
		"USE KEYS $ids RETURNING id, first_name, last_name, email, service", couchbaseBucket)

	idsString := make([]string, 1)
	idsString[0] = fmt.Sprintf(primaryKey, id)

	params := make(map[string]interface{})
	params["ids"] = idsString

	rows, err := repo.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(queryStr), params)

	if err != nil {
		return nil, err
	}

	for rows.Next(&user) {
		users = append(users, user)
		user = &User{}
	}

	if len(users) <= 0 {
		return nil, errors.New("Not found user")
	}

	return users[0], nil
}

func (repo *UserRepository) Update(userUpdate *User) (*User, error) {
	user := &User{}
	var users []*User

	queryStr := fmt.Sprintf("UPDATE `%s` SET "+
		"first_name=$first_name, "+
		"last_name=$last_name, "+
		"email=$email, "+
		"updated_at=$updated_at "+
		"WHERE type=$type and id=$id RETURNING id, first_name, last_name, email, service", couchbaseBucket)

	params := make(map[string]interface{})

	params["first_name"] = userUpdate.FirstName
	params["last_name"] = userUpdate.LastName
	params["email"] = userUpdate.Email
	params["type"] = table
	params["id"] = userUpdate.ID
	params["updated_at"] = time.Now()

	rows, err := repo.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(queryStr), params)

	if err != nil {
		return nil, err
	}

	for rows.Next(&user) {
		users = append(users, user)
		user = &User{}
	}

	if len(users) <= 0 {
		return nil, errors.New("Not found user")
	}

	return users[0], nil
}

func (repo *UserRepository) GetAll() ([]*User, error) {
	user := &User{}
	var users []*User
	queryStr := fmt.Sprintf("SELECT id, first_name, last_name, email, service FROM `%s` WHERE type=$type", couchbaseBucket)

	params := make(map[string]interface{})
	params["type"] = table

	rows, err := repo.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(queryStr), params)

	if err != nil {
		return nil, err
	}

	for rows.Next(&user) {
		users = append(users, user)
		user = &User{}
	}

	return users, nil
}

func (repo *UserRepository) GetByIDs(ids []uint64) ([]*User, error) {
	user := &User{}
	var users []*User
	queryStr := fmt.Sprintf("SELECT id, first_name, last_name, email, service FROM `%s` USE KEYS $ids", couchbaseBucket)
	idsString := make([]string, len(ids))

	for i, j := range ids {
		idsString[i] = fmt.Sprintf(primaryKey, j)
	}

	params := make(map[string]interface{})
	params["ids"] = idsString

	rows, err := repo.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(queryStr), params)

	if err != nil {
		return nil, err
	}

	for rows.Next(&user) {
		users = append(users, user)
		user = &User{}
	}

	return users, nil
}

func (repo *UserRepository) Get(id uint64) (*User, error) {

	user := &User{}
	var users []*User
	queryStr := fmt.Sprintf("SELECT id, first_name, last_name, email, service FROM `%s` USE KEYS $ids", couchbaseBucket)
	idsString := make([]string, 1)
	idsString[0] = fmt.Sprintf(primaryKey, id)

	params := make(map[string]interface{})
	params["ids"] = idsString

	rows, err := repo.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(queryStr), params)

	if err != nil {
		return nil, err
	}

	for rows.Next(&user) {
		users = append(users, user)
		user = &User{}
	}

	if len(users) <= 0 {
		return nil, errors.New("Not found user")
	}

	return users[0], nil
}

func (repo *UserRepository) GetByEmail(email string) (*User, error) {
	user := &User{}
	queryStr := fmt.Sprintf("SELECT id, first_name, last_name, email, service, `password` FROM `%s` WHERE email=$email and type=$type", couchbaseBucket)

	params := make(map[string]interface{})
	params["email"] = email
	params["type"] = table

	rows, err := repo.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(queryStr), params)

	if err != nil {
		return nil, err
	}

	users := []*User{}

	for rows.Next(&user) {
		users = append(users, user)
	}

	if len(users) <= 0 {
		return nil, errors.New("Not found user")
	}

	return users[0], nil
}

func (repo *UserRepository) Create(userCreate *User) (*User, error) {

	_, err := repo.GetByEmail(userCreate.Email)

	if err == nil {
		return nil, errors.New("User already exist")
	}

	if err != nil && err.Error() != "Not found user" {
		return nil, err
	}

	initialValue, _, err := repo.bucket.Counter("user_type", 1, 1, 0)

	if err != nil {
		return nil, err
	}

	userKey := fmt.Sprintf(primaryKey, initialValue)
	now := time.Now()
	userCreate.ID = initialValue
	userCreate.Type = table
	userCreate.CreatedAt = &now
	userCreate.UpdatedAt = &now

	user := &User{}
	var users []*User

	queryStr := fmt.Sprintf("INSERT INTO `%s` (KEY, VALUE) "+
		"VALUES ($key, $user) "+
		"RETURNING id, first_name, last_name, email, service, created_at, updated_at", couchbaseBucket)

	params := make(map[string]interface{})

	params["key"] = userKey
	params["user"] = userCreate

	rows, err := repo.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(queryStr), params)

	if err != nil {
		return nil, err
	}

	for rows.Next(&user) {
		users = append(users, user)
		user = &User{}
	}

	if len(users) <= 0 {
		return nil, errors.New("Not found user")
	}

	return users[0], nil
}
