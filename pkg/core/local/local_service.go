package local

import (
	"context"
	"errors"
	"fmt"
	"strings"

	i "github.com/pergamenum/api-gateway/pkg/core/interfaces"
	"github.com/pergamenum/go-utils-gin/logger"
	"go.uber.org/zap"
)

type UserServo struct {
	log  *zap.SugaredLogger
	repo i.Repository[User]
}

func NewUserServo(repo i.Repository[User]) *UserServo {
	return &UserServo{
		log:  logger.Get().Named("UserServo"),
		repo: repo,
	}
}

func (s *UserServo) Create(ctx context.Context, user User) error {

	err := user.isValid()
	if err != nil {
		return err
	}

	err = s.repo.Create(ctx, user.ID, user)
	if err != nil {
		return err
	}

	return nil

}

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
