package commands

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/graph"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/graph/generated"
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
)

func execRunCmd(cmd *cobra.Command, args []string) {
	httpPort := setupHttpPort()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

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
