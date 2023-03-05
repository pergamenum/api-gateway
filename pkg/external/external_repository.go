package external

import (
	"context"

	"cloud.google.com/go/firestore"
	i "github.com/pergamenum/api-gateway/pkg/core/interfaces"
	fsu "github.com/pergamenum/go-utils-firestore"
	"github.com/pergamenum/go-utils-gin/logger"
	"go.uber.org/zap"
)

type RepoConfig[E, M any] struct {
	Client  *firestore.Client
	Mapper  i.RepositoryMapper[E, M]
	LogName string
	Path    string
}

type Repo[E, M any] struct {
	log    *zap.SugaredLogger
	cruds  fsu.FirestoreCRUDS[E]
	mapper i.RepositoryMapper[E, M]
}

func NewRepo[E, M any](conf RepoConfig[E, M]) *Repo[E, M] {

	l := logger.Get().Named(conf.LogName)

	fc := fsu.NewCRUDS[E](conf.Client, conf.Path)

	return &Repo[E, M]{
		log:    l,
		mapper: conf.Mapper,
		cruds:  fc,
	}
}

func (r Repo[E, M]) Create(ctx context.Context, id string, model M) error {

	entity := r.mapper.ToEntity(model)
	err := r.cruds.Create(ctx, id, entity)
	switch err {
	case nil:
		// NO-OP
	case fsu.ErrDocumentAlreadyExists:
		return ErrAlreadyExists

	default:
		return err
	}

	return nil
}
