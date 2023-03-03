package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/pergamenum/api-gateway/pkg/core"
	"github.com/pergamenum/go-utils-gin/logger"
	"github.com/pergamenum/go-utils-gin/middleware"
)

func main() {

	log.Println("##### SERVICE STARTING #####")

	validateEnvironment()

	err := logger.Initialize()
	if err != nil {
		log.Fatal("main.logger.Initialize(): Failed:", err)
	}
	l := logger.Get()

	fsc, err := newFirestoreClient(context.Background(), os.Getenv("GCP_PROJECT_ID"))
	if err != nil {
		log.Fatal("main.communication.NewFirestoreClient(): Failed:", err)
	}
	defer fsc.Close()

	coreAPI := core.NewAPI(fsc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	middleware.AddRecovery(r, l.Desugar())
	middleware.AddRequestLogger(r, l.Desugar())

	r.POST("api/v1/user", coreAPI.User.CreateUser)
	r.GET("api/v1/user/:id", coreAPI.User.ReadUser)
	r.PATCH("api/v1/user", coreAPI.User.UpdateUser)
	r.DELETE("api/v1/user/:id", coreAPI.User.DeleteUser)
	r.GET("api/v1/user", coreAPI.User.SearchUsers)

	r.GET("/", func(c *gin.Context) {
		// TODO: Serve Swagger or other helpful information.
		c.String(http.StatusOK, "Pong!")
	})

	addr := fmt.Sprint("0.0.0.0:", os.Getenv("PORT"))
	// TODO: Graceful Shutdown
	log.Fatal(r.Run(addr))
}

func validateEnvironment() {

	var missing []string
	checkEnv := func(env string) {
		if os.Getenv(env) == "" {
			missing = append(missing, env)
		}
	}

	checkEnv("PORT")
	checkEnv("GCP_PROJECT_ID")

	if len(missing) > 0 {
		m := strings.Join(missing, ", ")
		log.Fatal(fmt.Sprint("Missing Environment Variables: ", m))
	}
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
