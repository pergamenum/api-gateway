package core

import (
	"context"
	"strings"
	"time"

	"cloud.google.com/go/firestore"

	fsu "github.com/pergamenum/go-utils-firestore"
	"github.com/pergamenum/go-utils-gin/logger"
	"go.uber.org/zap"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	ReadUser(ctx context.Context, id string) (User, error)
	UpdateUser(ctx context.Context, update UserUpdate) error
	DeleteUser(ctx context.Context, id string) error
	SearchUsers(ctx context.Context, query []UserQuery) ([]User, error)
}

// See User
type firestoreUser struct {
	ID      string    `firestore:"id"`
	Name    string    `firestore:"name"`
	Created time.Time `firestore:"created"`
	Updated time.Time `firestore:"updated"`
}

type UserFirestoreRepo struct {
	fc  fsu.FirestoreCRUDS[firestoreUser]
	log *zap.SugaredLogger
}

func NewUserFirestoreRepo(c *firestore.Client) *UserFirestoreRepo {
	fc := fsu.NewCRUDS[firestoreUser](c, "user")
	return &UserFirestoreRepo{
		fc:  fc,
		log: logger.Get().Named("core.UserFirestoreRepo"),
	}
}

func (r *UserFirestoreRepo) CreateUser(ctx context.Context, user User) error {

	now := time.Now()
	fu := r.fromUser(user)
	fu.Created = now
	fu.Updated = now

	err := r.fc.Create(ctx, fu.ID, fu)
	switch err {
	case nil:
		// NO-OP
	case fsu.ErrDocumentAlreadyExists:
		return errUserAlreadyExists

	default:
		return err
	}

	return nil
}

func (r *UserFirestoreRepo) ReadUser(ctx context.Context, id string) (User, error) {

	fu, err := r.fc.Read(ctx, id)
	switch err {
	case nil:
		// NO-OP
	case fsu.ErrDocumentNotFound:
		return User{}, errUserNotFound

	default:
		return User{}, err
	}

	user := r.toUser(fu)

	return user, nil
}

func (r *UserFirestoreRepo) UpdateUser(ctx context.Context, update UserUpdate) error {

	fus := r.fromUserUpdate(update)
	if len(fus) == 0 {
		return nil
	}

	fus = append(fus, firestore.Update{
		Path:  "updated",
		Value: time.Now(),
	})

	err := r.fc.Update(ctx, update.ID, fus)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserFirestoreRepo) DeleteUser(ctx context.Context, id string) error {

	err := r.fc.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserFirestoreRepo) SearchUsers(ctx context.Context, query []UserQuery) ([]User, error) {

	log := r.log.Named("SearchUsers")

	fqs := r.fromUserQuery(query)
	fus, err := r.fc.Search(ctx, fqs)
	if err != nil {
		if err == fsu.ErrDocumentSkipped {
			// Log the error, then carry on.
			log.Error(err)
		}
		return nil, err
	}

	var us []User
	for _, fu := range fus {
		u := r.toUser(fu)
		us = append(us, u)
	}

	return us, nil
}

func (r *UserFirestoreRepo) fromUser(input User) firestoreUser {
	return firestoreUser{
		ID:   input.ID,
		Name: input.Name,
	}
}

func (r *UserFirestoreRepo) toUser(input firestoreUser) User {
	return User(input)
}

func (r *UserFirestoreRepo) fromUserUpdate(input UserUpdate) []firestore.Update {

	var fus []firestore.Update

	if input.Name != nil {
		fus = append(fus, firestore.Update{
			Path:  "name",
			Value: *input.Name,
		})
	}

	return fus
}

func (r *UserFirestoreRepo) fromUserQuery(input []UserQuery) []fsu.Query {

	// fro = Firestore Relational Operator
	fro := func(input string) string {
		switch strings.ToUpper(input) {
		case "EQ":
			return "=="
		case "NE":
			return "!="
		case "LT":
			return "<"
		case "GT":
			return ">"
		case "LE":
			return "<="
		case "GE":
			return ">="
		default:
			return "UNKNOWN"
		}
	}

	var qs []fsu.Query
	for _, uq := range input {

		if uq.Key == "created" {
			if s, ok := uq.Value.(string); ok {
				t, _ := time.Parse(UserQueryTimeFormat, s)
				uq.Value = t
			}
		}

		fq := fsu.Query{
			Path:     uq.Key,
			Operator: fro(uq.Operator),
			Value:    uq.Value,
		}
		qs = append(qs, fq)
	}
	return qs
}
