package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/pergamenum/api-gateway/pkg/core"
	"github.com/pergamenum/api-gateway/pkg/monitoring/logger"
	u "github.com/pergamenum/api-gateway/pkg/utilities"
)

func main() {

	log.Println("##### SERVICE STARTING #####")

	validateEnvironment()

	err := logger.Initialize()
	if err != nil {
		log.Fatal("main.logger.Initialize(): Failed:", err)
	}

	fsc, err := newFirestoreClient(context.Background(), os.Getenv("GCP_PROJECT_ID"))
	if err != nil {
		log.Fatal("main.communication.NewFirestoreClient(): Failed:", err)
	}
	defer fsc.Close()

	coreAPI := core.NewAPI(fsc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	addMiddleware(r)

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

func addMiddleware(r *gin.Engine) {

	// TODO: Auth Middleware
	recovery := func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			u.ErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		u.ErrorResponse(c, http.StatusInternalServerError)
	}

	l := logger.Get().Desugar()
	r.Use(ginzap.CustomRecoveryWithZap(l, true, recovery))
	r.Use(ginzap.Ginzap(l, time.RFC3339, true))
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
