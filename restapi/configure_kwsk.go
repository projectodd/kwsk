// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	graceful "github.com/tylerb/graceful"

	"github.com/projectodd/kwsk/models"
	"github.com/projectodd/kwsk/restapi/operations"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

//go:generate swagger generate server --target .. --name Kwsk --spec ../apiv1swagger.json --principal models.Principal

var kwskFlags = struct {
	Master     string `long:"master" description:"Kubernetes Master URL"`
	Kubeconfig string `long:"kubeconfig" description:"Absolute path to the kubeconfig"`
	Istio      string `long:"istio" description:"Host and port of Istio Ingress service"`
}{}

func configureFlags(api *operations.KwskAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{
			ShortDescription: "Kubernetes Options",
			Options:          &kwskFlags,
		},
	}
}

func configureAPI(api *operations.KwskAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the Authorization header is set with the Basic scheme
	api.BasicAuthAuth = func(user string, pass string) (*models.Principal, error) {
		principal := models.Principal("someuser")
		return &principal, nil
	}

	knativeClient, err := knativeClient()
	if err != nil {
		log.Fatalf("Error creating Knative client: %s\n", err.Error())
	}

	activationCache := cache.NewTTLStore(activationKeyFunc, 30*time.Minute)

	configureActions(api, knativeClient, activationCache)
	configureActivations(api, knativeClient, activationCache)
	configurePackages(api, knativeClient)
	configureRules(api, knativeClient)
	configureTriggers(api, knativeClient)
	configureNamespaces(api, knativeClient)

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

func activationKeyFunc(obj interface{}) (string, error) {
	if activation, ok := obj.(models.Activation); ok {
		return *activation.ActivationID, nil
	}
	return "", fmt.Errorf("object is not an activation: %v", obj)
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

func knativeClient() (*knative.Clientset, error) {
	var kubeconfig string
	if kwskFlags.Kubeconfig != "" {
		kubeconfig = kwskFlags.Kubeconfig
	} else if home := homeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = ""
	}

	var config *rest.Config
	var err error
	if kwskFlags.Master == "" && kwskFlags.Kubeconfig == "" && os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err = rest.InClusterConfig()
	} else {
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags(kwskFlags.Master, kubeconfig)
	}
	if err != nil {
		return nil, err
	}

	knativeClient, err := knative.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return knativeClient, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
