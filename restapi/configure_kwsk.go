// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"path/filepath"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	graceful "github.com/tylerb/graceful"

	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/activations"
	"github.com/projectodd/kwsk/restapi/operations/namespaces"
	"github.com/projectodd/kwsk/restapi/operations/packages"
	"github.com/projectodd/kwsk/restapi/operations/rules"
	"github.com/projectodd/kwsk/restapi/operations/triggers"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

//go:generate swagger generate server --target .. --name Kwsk --spec ../../../../../../openwhisk/core/controller/build/resources/main/apiv1swagger.json --principal models.Principal

var kwskFlags = struct {
	Master     string `long:"master" description:"Kubernetes Master URL"`
	Kubeconfig string `long:"kubeconfig" description:"Absolute path to the kubeconfig"`
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

	knativeClient, err := knativeClient()
	if err != nil {
		log.Fatalf("Error creating Knative client: %s\n", err.Error())
	}

	configureActions(api, knativeClient)

	// TODO: This is super messy. Separate out all these generated
	// handlers by type at a minimum

	api.ActivationsGetNamespacesNamespaceActivationsActivationidLogsHandler = activations.GetNamespacesNamespaceActivationsActivationidLogsHandlerFunc(func(params activations.GetNamespacesNamespaceActivationsActivationidLogsParams) middleware.Responder {
		return middleware.NotImplemented("operation activations.GetNamespacesNamespaceActivationsActivationidLogs has not yet been implemented")
	})
	api.ActivationsGetNamespacesNamespaceActivationsActivationidResultHandler = activations.GetNamespacesNamespaceActivationsActivationidResultHandlerFunc(func(params activations.GetNamespacesNamespaceActivationsActivationidResultParams) middleware.Responder {
		return middleware.NotImplemented("operation activations.GetNamespacesNamespaceActivationsActivationidResult has not yet been implemented")
	})
	api.PackagesDeletePackageHandler = packages.DeletePackageHandlerFunc(func(params packages.DeletePackageParams) middleware.Responder {
		return middleware.NotImplemented("operation packages.DeletePackage has not yet been implemented")
	})
	api.RulesDeleteRuleHandler = rules.DeleteRuleHandlerFunc(func(params rules.DeleteRuleParams) middleware.Responder {
		return middleware.NotImplemented("operation rules.DeleteRule has not yet been implemented")
	})
	api.TriggersDeleteTriggerHandler = triggers.DeleteTriggerHandlerFunc(func(params triggers.DeleteTriggerParams) middleware.Responder {
		return middleware.NotImplemented("operation triggers.DeleteTrigger has not yet been implemented")
	})
	api.TriggersFireTriggerHandler = triggers.FireTriggerHandlerFunc(func(params triggers.FireTriggerParams) middleware.Responder {
		return middleware.NotImplemented("operation triggers.FireTrigger has not yet been implemented")
	})
	api.ActivationsGetActivationByIDHandler = activations.GetActivationByIDHandlerFunc(func(params activations.GetActivationByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation activations.GetActivationByID has not yet been implemented")
	})
	api.ActivationsGetActivationsHandler = activations.GetActivationsHandlerFunc(func(params activations.GetActivationsParams) middleware.Responder {
		return middleware.NotImplemented("operation activations.GetActivations has not yet been implemented")
	})
	api.PackagesGetAlPackagesHandler = packages.GetAlPackagesHandlerFunc(func(params packages.GetAlPackagesParams) middleware.Responder {
		return packages.NewGetAlPackagesOK()
	})
	api.NamespacesGetAllNamespacesHandler = namespaces.GetAllNamespacesHandlerFunc(func(params namespaces.GetAllNamespacesParams) middleware.Responder {
		return middleware.NotImplemented("operation namespaces.GetAllNamespaces has not yet been implemented")
	})
	api.RulesGetAllRulesHandler = rules.GetAllRulesHandlerFunc(func(params rules.GetAllRulesParams) middleware.Responder {
		return rules.NewGetAllRulesOK()
	})
	api.TriggersGetAllTriggersHandler = triggers.GetAllTriggersHandlerFunc(func(params triggers.GetAllTriggersParams) middleware.Responder {
		return triggers.NewGetAllTriggersOK()
	})
	api.PackagesGetPackageByNameHandler = packages.GetPackageByNameHandlerFunc(func(params packages.GetPackageByNameParams) middleware.Responder {
		return middleware.NotImplemented("operation packages.GetPackageByName has not yet been implemented")
	})
	api.RulesGetRuleByNameHandler = rules.GetRuleByNameHandlerFunc(func(params rules.GetRuleByNameParams) middleware.Responder {
		return middleware.NotImplemented("operation rules.GetRuleByName has not yet been implemented")
	})
	api.TriggersGetTriggerByNameHandler = triggers.GetTriggerByNameHandlerFunc(func(params triggers.GetTriggerByNameParams) middleware.Responder {
		return middleware.NotImplemented("operation triggers.GetTriggerByName has not yet been implemented")
	})
	api.RulesSetStateHandler = rules.SetStateHandlerFunc(func(params rules.SetStateParams) middleware.Responder {
		return middleware.NotImplemented("operation rules.SetState has not yet been implemented")
	})
	api.PackagesUpdatePackageHandler = packages.UpdatePackageHandlerFunc(func(params packages.UpdatePackageParams) middleware.Responder {
		return middleware.NotImplemented("operation packages.UpdatePackage has not yet been implemented")
	})
	api.RulesUpdateRuleHandler = rules.UpdateRuleHandlerFunc(func(params rules.UpdateRuleParams) middleware.Responder {
		return middleware.NotImplemented("operation rules.UpdateRule has not yet been implemented")
	})
	api.TriggersUpdateTriggerHandler = triggers.UpdateTriggerHandlerFunc(func(params triggers.UpdateTriggerParams) middleware.Responder {
		return middleware.NotImplemented("operation triggers.UpdateTrigger has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
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
