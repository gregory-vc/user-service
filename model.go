package main

import (
	"time"

	pb "github.com/gregory-vc/user-service/proto/user"
	"github.com/micro/protobuf/ptypes"
)

type User struct {
	ID        uint64     `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	FirstName string     `json:"first_name,omitempty"`
	LastName  string     `json:"last_name,omitempty"`
	Email     string     `json:"email,omitempty"`
	Password  string     `json:"password,omitempty"`
	Type      string     `json:"type,omitempty"`
	Service   string     `json:"service,omitempty"`
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
