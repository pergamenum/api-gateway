package core

import (
	"cloud.google.com/go/firestore"
	r "github.com/pergamenum/go-consensus-standards/repositories"
	fs "github.com/pergamenum/go-utils-firestore/dao"
	c "github.com/pergamenum/go-utils-gin/controllers"
	"go.uber.org/zap"
)

type API struct {
	User *c.Controller[User, UserDTO]
}

func NewAPI(client *firestore.Client, log *zap.SugaredLogger) *API {

	user := newUserController(client, log)

	return &API{User: user}
}

func newUserController(client *firestore.Client, log *zap.SugaredLogger) *c.Controller[User, UserDTO] {

	userManager := NewUserManager()

	dao := fs.NewDAO[UserEntity](client, "user", log)
	rConf := r.RepoConfig[User, UserEntity]{
		DAO:    dao,
		Mapper: userManager,
	}
	repo := r.NewRepo(rConf)

	sConf := ServoConfig{
		Repo:      repo,
		Validator: userManager,
		Log:       log,
	}
	servo := NewUserServo(sConf)

	cConf := c.ControllerConfig[User, UserDTO]{
		Service: servo,
		Mapper:  userManager,
	}
	controller := c.NewController[User, UserDTO](cConf)

	return controller
}
