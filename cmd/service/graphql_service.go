package service

import (
	"fmt"
	"net/http"
	"time"

	graphqlDelivery "github.com/azahir21/go-backend-boilerplate/internal/user/delivery/graphql"
	userUsecase "github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

func NewGraphQLServer(log *logrus.Logger, cfg config.GraphQLServerConfig, userUsecase userUsecase.UserUsecase) *http.Server {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	if cfg.StartupBanner {
		fmt.Println("🚀 Starting GraphQL server...")
	}

	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CorsOrigins,
		AllowMethods:     []string{"POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create GraphQL schema
	schema, err := graphqlDelivery.NewGraphQLSchema(log, userUsecase)
	if err != nil {
		log.Fatalf("failed to create GraphQL schema: %v", err)
	}

	// GraphQL endpoint
	router.POST("/graphql", func(c *gin.Context) {
		var r struct {
			Query     string                 `json:"query"`
			Operation string                 `json:"operationName"`
			Variables map[string]interface{} `json:"variables"`
		}
		if err := c.ShouldBindJSON(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  r.Query,
			VariableValues: r.Variables,
			OperationName:  r.Operation,
			Context:        c.Request.Context(),
		})

		c.JSON(http.StatusOK, result)
	})

	readTimeout, err := time.ParseDuration(cfg.ReadTimeout)
	if err != nil {
		log.Fatalf("invalid read timeout duration: %v", err)
	}

	writeTimeout, err := time.ParseDuration(cfg.WriteTimeout)
	if err != nil {
		log.Fatalf("invalid write timeout duration: %v", err)
	}

	idleTimeout, err := time.ParseDuration(cfg.IdleTimeout)
	if err != nil {
		log.Fatalf("invalid idle timeout duration: %v", err)
	}

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	return server
}
