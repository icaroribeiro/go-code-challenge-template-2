package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	authservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/auth"
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/healthcheck"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/application/service/user"
	authdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/auth"
	logindatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/login"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/infrastructure/storage/datastore/repository/user"
	graphqlhandler "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/presentation/handler/graphql"
	graphqlrouter "github.com/icaroribeiro/new-go-code-challenge-template-2/internal/transport/router/graphql"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/auth"
	datastorepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/datastore"
	envpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/env"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/adapter"
	handlerhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/handler"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/httputil/route"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/middleware/auth"
	securitypkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/security"
	serverpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/server"
	validatorpkg "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/validator"
	passwordvalidator "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/validator/password"
	usernamevalidator "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/validator/username"
	uuidvalidator "github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/validator/uuid"
	"github.com/spf13/cobra"
	validatorv2 "gopkg.in/validator.v2"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the API",
	Run:   execRunCmd,
}

var (
	deploy = envpkg.GetEnvWithDefaultValue("DEPLOY", "NO")

	httpPort = envpkg.GetEnvWithDefaultValue("HTTP_PORT", "8080")

	publicKeyPath                  = envpkg.GetEnvWithDefaultValue("RSA_PUBLIC_KEY_PATH", "./configs/auth/rsa_keys/rsa.public")
	privateKeyPath                 = envpkg.GetEnvWithDefaultValue("RSA_PRIVATE_KEY_PATH", "./configs/auth/rsa_keys/rsa.private")
	tokenExpTimeInSecStr           = envpkg.GetEnvWithDefaultValue("TOKEN_EXP_TIME_IN_SEC", "120")
	timeBeforeTokenExpTimeInSecStr = envpkg.GetEnvWithDefaultValue("TIME_BEFORE_TOKEN_EXP_TIME_IN_SEC", "30")

	dbDriver   = envpkg.GetEnvWithDefaultValue("DB_DRIVER", "postgres")
	dbUser     = envpkg.GetEnvWithDefaultValue("DB_USER", "postgres")
	dbPassword = envpkg.GetEnvWithDefaultValue("DB_PASSWORD", "postgres")
	dbHost     = envpkg.GetEnvWithDefaultValue("DB_HOST", "localhost")
	dbPort     = envpkg.GetEnvWithDefaultValue("DB_PORT", "5432")
	dbName     = envpkg.GetEnvWithDefaultValue("DB_NAME", "db")
)

func execRunCmd(cmd *cobra.Command, args []string) {
	httpPort := setupHttpPort()

	rsaKeys, err := setupRSAKeys()
	if err != nil {
		log.Panic(err.Error())
	}

	authN := authpkg.New(rsaKeys)

	tokenExpTimeInSec, err := strconv.Atoi(tokenExpTimeInSecStr)
	if err != nil {
		log.Panic(err.Error())
	}

	timeBeforeTokenExpTimeInSec, err := strconv.Atoi(timeBeforeTokenExpTimeInSecStr)
	if err != nil {
		log.Panic(err.Error())
	}

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

	authDatastoreRepository := authdatastorerepository.New(db)
	loginDatastoreRepository := logindatastorerepository.New(db)
	userDatastoreRepository := userdatastorerepository.New(db)

	validationFuncs := map[string]validatorv2.ValidationFunc{
		"uuid":     uuidvalidator.Validate,
		"username": usernamevalidator.Validate,
		"password": passwordvalidator.Validate,
	}

	validator, err := validatorpkg.New(validationFuncs)
	if err != nil {
		log.Panic(err.Error())
	}

	security := securitypkg.New()

	healthCheckService := healthcheckservice.New(db)
	authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
		authN, security, validator, tokenExpTimeInSec)
	userService := userservice.New(userDatastoreRepository, validator)

	graphqlHandler := graphqlhandler.New(healthCheckService, authService, userService, db, authN, timeBeforeTokenExpTimeInSec)

	adapters := map[string]adapterhttputilpkg.Adapter{
		"authMiddleware": authmiddlewarepkg.Auth(),
	}

	routes := make(routehttputilpkg.Routes, 0)
	routes = append(routes, graphqlrouter.ConfigureRoutes(graphqlHandler, adapters)...)

	router := setupRouter(routes)

	server := serverpkg.New(fmt.Sprintf(":%s", httpPort), router)

	idleChan := make(chan struct{})

	go func() {
		waitForShutdown(*server)
		close(idleChan)
	}()

	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		log.Panicf("%s", err.Error())
	}

	<-idleChan
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

// setupRSAKeys is the function that configures the RSA keys.
func setupRSAKeys() (authpkg.RSAKeys, error) {
	publicKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return authpkg.RSAKeys{}, fmt.Errorf("failed to read the RSA public key file: %s", err.Error())
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return authpkg.RSAKeys{}, fmt.Errorf("failed to parse the RSA public key: %s", err.Error())
	}

	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return authpkg.RSAKeys{}, fmt.Errorf("failed to read the RSA private key file: %s", err.Error())
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return authpkg.RSAKeys{}, fmt.Errorf("failed to parse the RSA private key: %s", err.Error())
	}

	return authpkg.RSAKeys{
		PublicKey:  rsaPublicKey,
		PrivateKey: rsaPrivateKey,
	}, nil
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

// setupRouter is the function that builds the router by arranging API routes.
func setupRouter(apiRoutes routehttputilpkg.Routes) *mux.Router {
	router := mux.NewRouter()

	methodNotAllowedHandler := handlerhttputilpkg.GetMethodNotAllowedHandler()
	router.MethodNotAllowedHandler = methodNotAllowedHandler

	notFoundHandler := handlerhttputilpkg.GetNotFoundHandler()
	router.NotFoundHandler = notFoundHandler

	for _, apiRoute := range apiRoutes {
		route := router.NewRoute()
		route.Name(apiRoute.Name)
		route.Methods(apiRoute.Method)
		route.Path(apiRoute.Path)
		route.HandlerFunc(apiRoute.HandlerFunc)
	}

	return router
}

// waitForShutdown is the function that waits for a signal to shutdown the server.
func waitForShutdown(server serverpkg.Server) {
	interruptChan := make(chan os.Signal, 1)

	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	<-interruptChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := server.Stop(ctx); err != nil && err != context.DeadlineExceeded {
		log.Panicf("%s", err.Error())
	}
}
