package service

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	sharedGraphQL "github.com/azahir21/go-backend-boilerplate/internal/shared/graphql"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

func NewGraphQLServer(log *logrus.Logger, cfg config.GraphQLServerConfig, modules []module.GraphQLModule) (*http.Server, error) {
	// Gin mode is set in cmd/app/app.go based on environment.
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

	// Collect schema builders from all modules
	var schemaBuilders []sharedGraphQL.SchemaBuilder
	for _, m := range modules {
		log.Infof("Registering GraphQL schema for module: %s", m.Name())
		schemaBuilders = append(schemaBuilders, m.GraphQLSchemaBuilder())
	}

	// Create GraphQL schema
	schema, err := sharedGraphQL.NewRootSchema(schemaBuilders)
	if err != nil {
		return nil, fmt.Errorf("failed to create GraphQL schema: %w", err)
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

	// GraphQL Playground endpoint
	router.GET("/graphql/playground", func(c *gin.Context) {
		playgroundHTML, err := os.ReadFile("web/playground.html")
		if err != nil {
			log.Errorf("failed to read playground.html: %v", err)
			c.String(http.StatusInternalServerError, "Failed to load GraphQL Playground")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", playgroundHTML)
	})

	// Serve static assets for GraphQL Playground
	router.GET("/graphql/playground/css/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join("web/css", filename)
		c.File(filePath)
	})
	router.GET("/graphql/playground/js/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join("web/js", filename)
		c.File(filePath)
	})

	readTimeout, err := parseDuration(cfg.ReadTimeout, "read timeout")
	if err != nil {
		return nil, err
	}

	writeTimeout, err := parseDuration(cfg.WriteTimeout, "write timeout")
	if err != nil {
		return nil, err
	}

	idleTimeout, err := parseDuration(cfg.IdleTimeout, "idle timeout")
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	if cfg.StartupBanner {
		log.Infof("🚀 Starting GraphQL server on :%s", cfg.Port)
	}

	return server, nil
}
