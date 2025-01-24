package services

import (
	"nevacarwash.com/main/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}
