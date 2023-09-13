// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error)
	InsertNewUser(ctx context.Context, input UsersEntity) (output UsersEntity, err error)
	Login(ctx context.Context, input UsersEntity) (output UsersEntity, err error)
	GetMyMyProfile(ctx context.Context, input UsersEntity) (output UsersEntity, err error)
	UpdateUser(ctx context.Context, input UsersEntity) (output UsersEntity, err error)
	UpdateLoginCount(ctx context.Context, input UsersEntity) (output UsersEntity, err error)
}
