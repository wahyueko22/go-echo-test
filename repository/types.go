// This file contains types that are used in the repository layer.
package repository

import "time"

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type UsersEntity struct {
	Id                int
	FullName          string
	Password          string
	PhoneNumber       string
	SuccessLoginCount int8
	CreatedAt         time.Time
}
