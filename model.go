package main

import (
	"time"

	pb "github.com/gregory-vc/user-service/proto/user"
	"github.com/micro/protobuf/ptypes"
)

type User struct {
	ID        uint64     `db:"id"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	FirstName string     `db:"first_name"`
	LastName  string     `db:"last_name"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	Type      string     `db:"type"`
	Service   string     `db:"service"`
}

func modelFromProtobuf(p *pb.User) *User {

	createdAt, _ := ptypes.Timestamp(p.CreatedAt)
	updatedAt, _ := ptypes.Timestamp(p.UpdatedAt)

	ret := &User{
		ID:        p.Id,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Email:     p.Email,
		Password:  p.Password,
		Type:      p.Type,
		Service:   p.Service,
	}

	return ret
}

func protobufFromModel(p *User) *pb.User {

	createdAt, _ := ptypes.TimestampProto(*p.CreatedAt)
	updatedAt, _ := ptypes.TimestampProto(*p.UpdatedAt)

	ret := &pb.User{
		Id:        p.ID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Email:     p.Email,
		Password:  p.Password,
		Type:      p.Type,
		Service:   p.Service,
	}

	return ret
}

func protobufFromModelList(pl []*User) []*pb.User {
	ret := make([]*pb.User, 0, len(pl))
	for _, p := range pl {
		ret = append(ret, protobufFromModel(p))
	}

	return ret
}
