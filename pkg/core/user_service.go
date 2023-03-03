package core

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pergamenum/go-utils-gin/logger"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(ctx context.Context, user User) error
	ReadUser(ctx context.Context, id string) (User, error)
	UpdateUser(ctx context.Context, update UserUpdate) error
	DeleteUser(ctx context.Context, id string) error
	SearchUsers(ctx context.Context, query []UserQuery) ([]User, error)
}

type User struct {
	ID      string
	Name    string
	Created time.Time
	Updated time.Time
}

type UserUpdate struct {
	ID   string
	Name *string
}

type UserQuery struct {
	// Key should match the name of a User field.
	Key string
	// Operator represents relational operator, ex: EQ, LT, ...
	Operator string
	// Value should be a value of Key's underlying User field type, ex: string for Name.
	Value any
}

type UserServo struct {
	repo UserRepository
	log  *zap.SugaredLogger
}

func NewUserServo(repo UserRepository) *UserServo {
	return &UserServo{
		repo: repo,
		log:  logger.Get().Named("Core.UserServo"),
	}
}

func (s *UserServo) CreateUser(ctx context.Context, user User) error {

	err := user.isValid()
	if err != nil {
		return err
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServo) ReadUser(ctx context.Context, id string) (User, error) {

	log := s.log.Named("ReadUser")

	if id == "" {
		return User{}, errUserIDEmpty
	}

	user, err := s.repo.ReadUser(ctx, id)
	if err != nil {
		return User{}, err
	}

	err = user.isValid()
	if err != nil {
		log.With("user", user).
			Error(errUserFoundInvalid)
		return User{}, errUserFoundInvalid
	}

	return user, nil
}

func (s *UserServo) UpdateUser(ctx context.Context, update UserUpdate) error {

	err := update.isValid()
	if err != nil {
		return err
	}

	if !s.userFound(ctx, update.ID) {
		return errUserNotFound
	}

	err = s.repo.UpdateUser(ctx, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServo) DeleteUser(ctx context.Context, id string) error {

	if id == "" {
		return errUserIDEmpty
	}

	if !s.userFound(ctx, id) {
		return errUserNotFound
	}

	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServo) SearchUsers(ctx context.Context, query []UserQuery) ([]User, error) {

	log := s.log.Named("SearchUsers")

	for _, q := range query {
		err := q.isValid()
		if err != nil {
			return nil, err
		}
	}

	users, err := s.repo.SearchUsers(ctx, query)
	if err != nil {
		return nil, err
	}

	var vus []User
	for _, vu := range users {
		err = vu.isValid()
		if err != nil {
			log.With("user", vu).
				Error(errUserFoundInvalid)
			continue
		}
		vus = append(vus, vu)
	}
	// Empty slice instead of nil slice. We don't want to send null body responses...?
	if vus == nil {
		vus = []User{}
	}
	return vus, nil
}

func (s *UserServo) userFound(ctx context.Context, id string) bool {

	_, err := s.ReadUser(ctx, id)
	return err == nil
}

const UserQueryTimeFormat = "2006-01-02_15:04"

var (
	errUserAlreadyExists = errors.New("user already exists")
	errUserNotFound      = errors.New("user not found")
	errUserFoundInvalid  = errors.New("user found, but invalid")
	errUserIDEmpty       = errors.New("user id empty")
)

func (u User) isValid() error {

	var sb strings.Builder

	if u.ID == "" {
		sb.WriteString("(id: empty) ")
	}

	if len(u.Name) > 100 {
		sb.WriteString("(name: max 100 chars) ")
	}

	if len(sb.String()) > 0 {
		cause := fmt.Sprint("invalid user: ", strings.TrimSpace(sb.String()))
		return errors.New(cause)
	} else {
		return nil
	}
}

func (u UserUpdate) isValid() error {

	user := User{
		ID: u.ID,
	}

	if u.Name != nil {
		user.Name = *u.Name
	}

	return user.isValid()
}

func (q UserQuery) isValid() error {

	var sb strings.Builder

	validKeys := []string{"id", "name", "created"}
	keyValid := false
	for _, vk := range validKeys {
		if strings.EqualFold(vk, q.Key) {
			keyValid = true
			break
		}
	}
	if !keyValid {
		sb.WriteString(fmt.Sprintf("(key[%v]: invalid) ", q.Key))
	}

	validOperators := []string{"EQ", "NE", "LT", "GT", "LE", "GE"}
	operatorValid := false
	for _, vo := range validOperators {
		if strings.EqualFold(vo, q.Operator) {
			operatorValid = true
			break
		}
	}
	if !operatorValid {
		sb.WriteString(fmt.Sprintf("(operator[%v]: invalid) ", q.Operator))
	}

	valueValid := false
	switch strings.ToLower(q.Key) {
	case "id", "name":
		if _, ok := q.Value.(string); ok {
			valueValid = true
		}
	case "created":
		if s, ok := q.Value.(string); ok {
			_, err := time.Parse(UserQueryTimeFormat, s)
			if err != nil {
				sb.WriteString("time must be: YYYY-MM-DD_hh:mm ")
				break
			}
			valueValid = true
		}
	}
	if !valueValid && keyValid {
		sb.WriteString(fmt.Sprintf("(value[%v]: invalid) ", q.Value))
	}

	if len(sb.String()) > 0 {
		cause := fmt.Sprint("invalid query: ", strings.TrimSpace(sb.String()))
		return errors.New(cause)
	} else {
		return nil
	}
}
