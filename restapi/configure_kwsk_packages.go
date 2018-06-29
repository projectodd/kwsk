package restapi

import (
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/packages"

	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

func configurePackages(api *operations.KwskAPI, knativeClient *knative.Clientset) {
	api.PackagesDeletePackageHandler = packages.DeletePackageHandlerFunc(func(params packages.DeletePackageParams) middleware.Responder {
		return middleware.NotImplemented("operation packages.DeletePackage has not yet been implemented")
	})
	api.PackagesGetAlPackagesHandler = packages.GetAlPackagesHandlerFunc(func(params packages.GetAlPackagesParams) middleware.Responder {
		return packages.NewGetAlPackagesOK()
	})
	api.PackagesGetPackageByNameHandler = packages.GetPackageByNameHandlerFunc(func(params packages.GetPackageByNameParams) middleware.Responder {
		return middleware.NotImplemented("operation packages.GetPackageByName has not yet been implemented")
	})
	api.PackagesUpdatePackageHandler = packages.UpdatePackageHandlerFunc(func(params packages.UpdatePackageParams) middleware.Responder {
		return middleware.NotImplemented("operation packages.UpdatePackage has not yet been implemented")
	})
}
