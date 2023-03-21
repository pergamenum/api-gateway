package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/pergamenum/api-gateway/pkg/core"
	"github.com/pergamenum/go-consensus-standards/setup"
	"github.com/pergamenum/go-utils-gin/logger"
	"github.com/pergamenum/go-utils-gin/middleware"
)

var RequiredEnvironment = []string{
	"PORT",
	"GCP_PROJECT_ID",
}

func main() {

	err := logger.Initialize()
	if err != nil {
		cause := fmt.Sprintln("main.logger.Initialize(): Failed: ", err)
		panic(cause)
	}
	log := logger.Get()

	log.Info("##### SERVICE STARTING #####")

	err = setup.ValidateEnvironment(RequiredEnvironment)
	if err != nil {
		log.Fatal(err)
	}

	fsc, err := newFirestoreClient(context.Background(), os.Getenv("GCP_PROJECT_ID"))
	if err != nil {
		log.With("error", err).Fatal("main.newFirestoreClient(): Failed")
	}
	defer fsc.Close()

	coreAPI := core.NewAPI(fsc, log)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	middleware.AddRecovery(r, log.Desugar())
	middleware.AddRequestLogger(r, log.Desugar())

	r.POST("api/v1/user", coreAPI.User.Create)
	r.GET("api/v1/user/:id", coreAPI.User.Read)
	r.PATCH("api/v1/user", coreAPI.User.Update)
	r.DELETE("api/v1/user/:id", coreAPI.User.Delete)
	r.GET("api/v1/user", coreAPI.User.Search)

	r.GET("/", func(c *gin.Context) {
		// TODO: Serve Swagger or other helpful information.
		// TODO: Serve User:Password HTML for manipulating remote config.
		c.String(http.StatusOK, "Pong!")
	})

	addr := fmt.Sprint("0.0.0.0:", os.Getenv("PORT"))
	// TODO: Graceful Shutdown
	log.Fatal(r.Run(addr))
}

func newFirestoreClient(ctx context.Context, projectID string) (*firestore.Client, error) {

	conf := &firebase.Config{
		ProjectID: projectID,
	}
	fba, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, err
	}

	fsc, err := fba.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return fsc, nil
}
