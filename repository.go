package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/gregory-vc/user-service/proto/user"
	"gopkg.in/couchbase/gocb.v1"
)

const primaryKey = "user_%d"
const table = "user"

type Repository interface {
	GetAll() ([]*pb.User, error)
	GetByIDs(ids []uint64) ([]*pb.User, error)
	Get(id uint64) (*pb.User, error)
	Delete(id uint64) (*pb.User, error)
	Create(user *pb.User) error
	Update(userUpdate *pb.User) (*pb.User, error)
	GetByEmail(email string) (*pb.User, error)
}

type UserRepository struct {
	bucket *gocb.Bucket
}

func (repo *UserRepository) Delete(id uint64) (*pb.User, error) {
	user := &pb.User{}
	var users []*pb.User
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
		user = &pb.User{}
	}

	if len(users) <= 0 {
		return nil, errors.New("Not found user")
	}

	return users[0], nil
}

func (repo *UserRepository) Update(userUpdate *pb.User) (*pb.User, error) {
	user := &pb.User{}
	var users []*pb.User

	queryStr := fmt.Sprintf("UPDATE `%s` SET "+
		"first_name=$first_name, "+
		"last_name=$last_name, "+
		"email=$email "+
		"updated_at=$updated_at "+
		"WHERE type=$type and id=$id RETURNING id, first_name, last_name, email, service", couchbaseBucket)

	params := make(map[string]interface{})

	t, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}

	params["first_name"] = userUpdate.FirstName
	params["last_name"] = userUpdate.LastName
	params["email"] = userUpdate.Email
	params["type"] = table
	params["id"] = userUpdate.Id
	params["updated_at"] = t

	rows, err := repo.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(queryStr), params)

	if err != nil {
		return nil, err
	}

	for rows.Next(&user) {
		users = append(users, user)
		user = &pb.User{}
	}

	if len(users) <= 0 {
		return nil, errors.New("Not found user")
	}

	return users[0], nil
}

func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	user := &pb.User{}
	var users []*pb.User
	queryStr := fmt.Sprintf("SELECT id, first_name, last_name, email, service FROM `%s` WHERE type=$type", couchbaseBucket)

	params := make(map[string]interface{})
	params["type"] = table

	rows, err := repo.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(queryStr), params)

	if err != nil {
		return nil, err
	}

	for rows.Next(&user) {
		users = append(users, user)
		user = &pb.User{}
	}

	return users, nil
}

func (repo *UserRepository) GetByIDs(ids []uint64) ([]*pb.User, error) {
	user := &pb.User{}
	var users []*pb.User
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
		user = &pb.User{}
	}

	return users, nil
}

func (repo *UserRepository) Get(id uint64) (*pb.User, error) {

	user := &pb.User{}
	var users []*pb.User
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
		user = &pb.User{}
	}

	if len(users) <= 0 {
		return nil, errors.New("Not found user")
	}

	return users[0], nil
}

func (repo *UserRepository) GetByEmail(email string) (*pb.User, error) {
	user := &pb.User{}
	queryStr := fmt.Sprintf("SELECT id, first_name, last_name, email, service, `password` FROM `%s` WHERE email=$email and type=$type", couchbaseBucket)

	params := make(map[string]interface{})
	params["email"] = email
	params["type"] = table

	rows, err := repo.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(queryStr), params)

	if err != nil {
		return nil, err
	}

	users := []*pb.User{}

	for rows.Next(&user) {
		users = append(users, user)
	}

	if len(users) <= 0 {
		return nil, errors.New("Not found user")
	}

	return users[0], nil
}

func (repo *UserRepository) Create(user *pb.User) error {

	_, err := repo.GetByEmail(user.Email)

	if err == nil {
		return errors.New("User already exist")
	}

	if err != nil && err.Error() != "Not found user" {
		return err
	}

	initialValue, _, err := repo.bucket.Counter("user_type", 1, 1, 0)

	if err != nil {
		return err
	}

	userKey := fmt.Sprintf(primaryKey, initialValue)
	user.Id = initialValue
	user.Type = table
	t, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return err
	}
	user.CreatedAt = t
	user.UpdatedAt = t

	_, err = repo.bucket.Insert(userKey, user, 0)

	if err != nil {
		return err
	}

	return nil
}
