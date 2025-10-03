package main

import (
	"context"

	userv1 "github.com/Daniel-Sogbey/micro-weekend/proto/user/v1"
)

type UserServer struct {
	userv1.UnimplementedUserServiceServer
	repo *Repo
}

func (u *UserServer) CreateUser(ctx context.Context, in *userv1.CreateUserRequest) (*userv1.User, error) {
	user := &User{
		Name:  in.Name,
		Email: in.Email,
	}

	err := u.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &userv1.User{Id: user.Id, Name: user.Name, Email: user.Email, CreatedUnix: user.CreatedUnix}, nil
}

func (u *UserServer) GetUser(ctx context.Context, in *userv1.GetUserRequest) (*userv1.User, error) {
	user, err := u.repo.Get(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &userv1.User{Id: user.Id, Name: user.Name, Email: user.Email, CreatedUnix: user.CreatedUnix}, nil
}
