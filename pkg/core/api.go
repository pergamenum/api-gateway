package core

import (
	"cloud.google.com/go/firestore"
)

type API struct {
	User *UserController
}

func NewAPI(fsc *firestore.Client) *API {

	userRepo := NewUserFirestoreRepo(fsc)
	userServo := NewUserServo(userRepo)
	userController := NewUserController(userServo)

	return &API{User: userController}
}
