package local

import (
	"cloud.google.com/go/firestore"
	e "github.com/pergamenum/api-gateway/pkg/external"
)

type API struct {
	User *e.Controller[UserDTO, User, UserUpdate]
}

func NewAPI(client *firestore.Client) *API {

	user := newUserController(client)

	return &API{User: user}
}

func newUserController(client *firestore.Client) *e.Controller[UserDTO, User, UserUpdate] {

	mapper := UserMapper{}

	rConf := e.RepoConfig[UserEntity, User]{
		Client:  client,
		Mapper:  mapper,
		LogName: "UserRepo",
		Path:    "user",
	}

	repo := e.NewRepo(rConf)

	servo := NewUserServo(repo)

	cConf := e.ControllerConfig[UserDTO, User, UserUpdate]{
		Service: servo,
		Mapper:  mapper,
		LogName: "UserController",
	}
	c := e.NewController[UserDTO, User, UserUpdate](cConf)

	return c
}
