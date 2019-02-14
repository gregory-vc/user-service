package main

import (
	"errors"
	"fmt"

	pb "github.com/gregory-vc/user-service/proto/user"
	"gopkg.in/couchbase/gocb.v1"
)

type Repository interface {
	GetAll() ([]*pb.User, error)
	Get(id uint64) (*pb.User, error)
	Create(user *pb.User) error
	GetByEmailAndPassword(user *pb.User) (*pb.User, error)
	GetByEmail(email string) (*pb.User, error)
}

type UserRepository struct {
	bucket *gocb.Bucket
}

func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	user := &pb.User{}
	var users []*pb.User
	queryStr := fmt.Sprintf("SELECT id, first_name, last_name, email FROM `%s` WHERE type=$type", couchbaseBucket)

	params := make(map[string]interface{})
	params["type"] = "user"

	rows, err := repo.bucket.ExecuteN1qlQuery(gocb.NewN1qlQuery(queryStr), params)

	if err != nil {
		return nil, err
	}

	for rows.Next(&user) {
		fmt.Println("1", user)
		users = append(users, user)
		fmt.Println(users)
	}

	return users, nil
}

func (repo *UserRepository) Get(id uint64) (*pb.User, error) {
	var user *pb.User
	// user.Id = id
	// if err := repo.db.First(&user).Error; err != nil {
	// 	return nil, err
	// }
	return user, nil
}

func (repo *UserRepository) GetByEmailAndPassword(user *pb.User) (*pb.User, error) {
	// if err := repo.db.First(&user).Error; err != nil {
	// 	return nil, err
	// }
	return user, nil
}

func (repo *UserRepository) GetByEmail(email string) (*pb.User, error) {
	user := &pb.User{}
	queryStr := fmt.Sprintf("SELECT id, first_name, last_name, email, `password` FROM `%s` WHERE email=$email and type=$type", couchbaseBucket)

	params := make(map[string]interface{})
	params["email"] = email
	params["type"] = "user"

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
	} else {
		return users[0], nil
	}

	return user, nil
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

	userKey := fmt.Sprintf("user_%d", initialValue)
	user.Id = initialValue
	user.Type = "user"

	_, err = repo.bucket.Insert(userKey, user, 0)

	if err != nil {
		return err
	}

	return nil
}
