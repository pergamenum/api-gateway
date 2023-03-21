package core

import (
	"context"
	"errors"

	e "github.com/pergamenum/go-consensus-standards/ehandler"
	i "github.com/pergamenum/go-consensus-standards/interfaces"
	t "github.com/pergamenum/go-consensus-standards/types"
	"go.uber.org/zap"
)

type UserServo struct {
	repo      i.Repository[User]
	validator i.Validator[User]
	log       *zap.SugaredLogger
}

type ServoConfig struct {
	Repo      i.Repository[User]
	Validator i.Validator[User]
	Log       *zap.SugaredLogger
}

func NewUserServo(config ServoConfig) *UserServo {

	namedLogger := config.Log.Named("UserServo")

	return &UserServo{
		log:       namedLogger,
		validator: config.Validator,
		repo:      config.Repo,
	}
}

func (s *UserServo) Create(ctx context.Context, user User) error {

	err := s.validator.ValidateModel(user)
	if err != nil {
		return e.Wrap(err, e.ErrBadRequest)
	}

	err = s.repo.Create(ctx, user.ID, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServo) Read(ctx context.Context, id string) (User, error) {

	log := s.log.Named("Read")

	if id == "" {
		return User{}, e.Wrap("(id missing)", e.ErrBadRequest)
	}

	user, err := s.repo.Read(ctx, id)
	if err != nil {
		return User{}, err
	}

	err = s.validator.ValidateModel(user)
	if err != nil {
		log.With("user", user).
			Error(e.ErrCorrupt)
		return User{}, e.Wrap(err, e.ErrCorrupt)
	}

	return user, nil

}

func (s *UserServo) Update(ctx context.Context, update t.Update) error {

	err := s.validator.ValidateUpdate(update)
	if err != nil {
		return e.Wrap(err, e.ErrBadRequest)
	}

	id := update["id"].(string)

	err = s.userFound(ctx, id)
	if err != nil {
		// Allow updating a corrupt resource.
		if !errors.Is(err, e.ErrCorrupt) {
			return err
		}
	}

	err = s.repo.Update(ctx, id, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServo) Delete(ctx context.Context, id string) error {

	if id == "" {
		return e.Wrap("(id missing)", e.ErrBadRequest)
	}

	err := s.userFound(ctx, id)
	if err != nil {
		// Allow deleting a corrupt resource.
		if !errors.Is(err, e.ErrCorrupt) {
			return err
		}
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServo) Search(ctx context.Context, query []t.Query) ([]User, error) {

	log := s.log.Named("Search")

	err := s.validator.ValidateQuery(query)
	if err != nil {
		err = e.Wrap(err, e.ErrBadRequest)
		return nil, err
	}

	users, err := s.repo.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	var vus []User
	for _, vu := range users {
		err = s.validator.ValidateModel(vu)
		if err != nil {
			log.With("user", vu).
				Error(e.ErrCorrupt)
			continue
		}
		vus = append(vus, vu)
	}
	return vus, nil
}

func (s *UserServo) userFound(ctx context.Context, id string) error {

	_, err := s.Read(ctx, id)
	return err
}
