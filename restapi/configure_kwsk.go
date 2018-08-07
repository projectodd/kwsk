// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	graceful "github.com/tylerb/graceful"

	"github.com/projectodd/kwsk/models"
	"github.com/projectodd/kwsk/restapi/operations"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	eventing "github.com/knative/eventing/pkg/client/clientset/versioned"
	serving "github.com/knative/serving/pkg/client/clientset/versioned"
)

//go:generate swagger generate server --target .. --name Kwsk --spec ../apiv1swagger.json --principal models.Principal

var kwskFlags = struct {
	Master      string `long:"master" description:"Kubernetes Master URL"`
	Kubeconfig  string `long:"kubeconfig" description:"Absolute path to the kubeconfig"`
	Istio       string `long:"istio" description:"Host and port of Istio Ingress service"`
	ImagePrefix string `long:"image-prefix" description:"Image prefix for action runtime images"`
	ImageTag    string `long:"image-tag" description:"Image tag for action runtime images"`
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

	config, err := kubeConfig()
	if err != nil {
		log.Fatalf("Error creating Kubernetes client config: %s\n", err.Error())
	}
	servingClient, err := servingClient(config)
	if err != nil {
		log.Fatalf("Error creating Knative serving client: %s\n", err.Error())
	}
	eventingClient, err := eventingClient(config)
	if err != nil {
		log.Fatalf("Error creating Knative eventing client: %s\n", err.Error())
	}

	activationCache := cache.NewTTLStore(activationKeyFunc, 30*time.Minute)

	configureActions(api, servingClient, activationCache)
	configureActivations(api, servingClient, activationCache)
	configurePackages(api, servingClient)
	configureRules(api, eventingClient)
	configureTriggers(api, servingClient, eventingClient)
	configureNamespaces(api, servingClient)

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

func kubeConfig() (*rest.Config, error) {
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
	return config, nil
}

func servingClient(config *rest.Config) (*serving.Clientset, error) {
	client, err := serving.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func eventingClient(config *rest.Config) (*eventing.Clientset, error) {
	client, err := eventing.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func namespaceOrDefault(namespace string) string {
	if namespace == "_" {
		// TODO: In OpenWhisk land, the "_" namespace means the
		// default namespace of the authenticated user. Because we're
		// not dealing with any auth yet, just hardcode this to the
		// "default" namespace for now.
		namespace = "default"
	}
	return namespace
}

func sanitizeObjectName(name string) string {
	return strings.Replace(strings.ToLower(name), " ", "-", -1)
}

func errorMessageFromErr(err error) *models.ErrorMessage {
	msg := err.Error()
	return &models.ErrorMessage{
		Error: &msg,
	}
}

func newActivationId() (string, error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return strings.Replace(newUuid.String(), "-", "", -1), nil
}

func istioHostAndPort() string {
	// If we're running in-cluster this needs to be an internal
	// hostname. If we're running outside the cluster, this needs
	// to be the exposed route and/or nodeport. For now, don't
	// worry about magic and expect it to be explicitly configured
	// via a flag.
	//
	// host := "istio-ingress.istio-system.svc.cluster.local"
	istioHostAndPort := kwskFlags.Istio
	if istioHostAndPort == "" {
		panic("Istio host and port must be provided via --istio flag to invoke actions")
	}
	return istioHostAndPort
}

func imagePrefix() string {
	prefix := kwskFlags.ImagePrefix
	if prefix == "" {
		prefix = "projectodd"
	}
	return prefix
}

func imageTag() string {
	tag := kwskFlags.ImageTag
	if tag == "" {
		tag = "latest"
	}
	return tag
}

const (
	KwskName    string = "kwsk_name"
	KwskVersion string = "kwsk_version"
	KwskKind    string = "kwsk_kind"
	KwskImage   string = "kwsk_image"
)
