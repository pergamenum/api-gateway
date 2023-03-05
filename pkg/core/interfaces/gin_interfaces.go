package interfaces

import "context"

/*
type Controller[DTO, Update, Query any] interface {
	Create(ctx *gin.Context)
	//Read(ctx context.Context, id string) (DTO, error)
	//Update(ctx context.Context, update Update) error
	//Delete(ctx context.Context, id string) error
	//Search(ctx context.Context, query Query) ([]DTO, error)
}
*/

type ControllerMapper[D, M, U any] interface {
	ToDTO(M) D
	ToModel(D) M
	ToUpdate(D) U
}

// type Service[Model, Update, Query any] interface {

type Service[Model any] interface {
	Create(ctx context.Context, model Model) error
	//Read(ctx context.Context, id string) (Model, error)
	//Update(ctx context.Context, update Update) error
	//Delete(ctx context.Context, id string) error
	//Search(ctx context.Context, query Query) ([]Model, error)
}

//type Repository[Entity, Update, Query any] interface {

type Repository[Model any] interface {
	Create(ctx context.Context, id string, model Model) error
	//Read(ctx context.Context, id string) (DTO, error)
	//Update(ctx context.Context, update Update) error
	//Delete(ctx context.Context, id string) error
	//Search(ctx context.Context, query Query) ([]DTO, error)
}

type RepositoryMapper[E, M any] interface {
	ToEntity(M) E
	FromEntity(E) M
}
