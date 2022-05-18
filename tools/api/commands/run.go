package commands

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/healthcheck"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/graph/generated"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/graph/resolver"
	datastorepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/datastore"
	envpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/env"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the API",
	Run:   execRunCmd,
}

var (
	deploy = envpkg.GetEnvWithDefaultValue("DEPLOY", "NO")

	httpPort = envpkg.GetEnvWithDefaultValue("HTTP_PORT", "8080")

	dbDriver   = envpkg.GetEnvWithDefaultValue("DB_DRIVER", "postgres")
	dbUser     = envpkg.GetEnvWithDefaultValue("DB_USER", "postgres")
	dbPassword = envpkg.GetEnvWithDefaultValue("DB_PASSWORD", "postgres")
	dbHost     = envpkg.GetEnvWithDefaultValue("DB_HOST", "localhost")
	dbPort     = envpkg.GetEnvWithDefaultValue("DB_PORT", "5432")
	dbName     = envpkg.GetEnvWithDefaultValue("DB_NAME", "db")
)

func execRunCmd(cmd *cobra.Command, args []string) {
	httpPort := setupHttpPort()

	dbConfig, err := setupDBConfig()
	if err != nil {
		log.Panic(err.Error())
	}

	datastore, err := datastorepkg.New(dbConfig)
	if err != nil {
		log.Panic(err.Error())
	}
	defer datastore.Close()

	db := datastore.GetInstance()
	if db == nil {
		log.Panicf("The database instance is null")
	}

	if err = db.Error; err != nil {
		log.Panicf("Got error when acessing the database instance: %s", err.Error())
	}

	healthCheckService := healthcheckservice.New(db)

	res := resolver.NewResolver(healthCheckService)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: res}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", httpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil))
}

// setupHttpPort is the function that configures the port address used by the server.
func setupHttpPort() string {
	if deploy == "YES" {
		if httpPort = os.Getenv("PORT"); httpPort == "" {
			log.Panicf("failed to read the PORT env variable to the application deployment")
		}
	}

	return httpPort
}

// setupDBConfig is the function that configures a map of parameters used to connect to the database.
func setupDBConfig() (map[string]string, error) {
	dbURL := ""

	if deploy == "YES" {
		if dbURL = os.Getenv("DATABASE_URL"); dbURL == "" {
			return nil, fmt.Errorf("failed to read the DATABASE_URL environment variable to the application deployment")
		}
	}

	dbConfig := map[string]string{
		"DRIVER":   dbDriver,
		"USER":     dbUser,
		"PASSWORD": dbPassword,
		"HOST":     dbHost,
		"PORT":     dbPort,
		"NAME":     dbName,
		"URL":      dbURL,
	}

	return dbConfig, nil
}
