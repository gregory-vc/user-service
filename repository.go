package main

import (
	"fmt"

	pb "github.com/gregory-vc/user-service/proto/user"
	"gopkg.in/couchbase/gocb.v1"
)

type Repository interface {
	GetAll() ([]*pb.User, error)
	Get(id uint32) (*pb.User, error)
	Create(user *pb.User) error
	GetByEmailAndPassword(user *pb.User) (*pb.User, error)
	GetByEmail(email string) (*pb.User, error)
}

type UserRepository struct {
	bucket *gocb.Bucket
}

func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	var users []*pb.User
	// if err := repo.db.Find(&users).Error; err != nil {
	// 	return nil, err
	// }
	return users, nil
}

func (repo *UserRepository) Get(id uint32) (*pb.User, error) {
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
	// if err := repo.db.Where("email = ?", email).
	// 	First(&user).Error; err != nil {
	// 	return nil, err
	// }
	return user, nil
}

func (repo *UserRepository) Create(user *pb.User) error {
	userKey := fmt.Sprintf("user::%s", user.Email)
	user.Id = 1
	user.Type = "user"
	_, err := repo.bucket.Insert(userKey, user, 0)
	if err != nil {
		return err
	}
	return nil
}
