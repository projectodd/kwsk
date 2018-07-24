package restapi

import (
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/projectodd/kwsk/models"
	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/packages"

	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

func configurePackages(api *operations.KwskAPI, knativeClient *knative.Clientset) {
	api.PackagesDeletePackageHandler = packages.DeletePackageHandlerFunc(func(params packages.DeletePackageParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation packages.DeletePackage has not yet been implemented")
	})
	api.PackagesGetAllPackagesHandler = packages.GetAllPackagesHandlerFunc(func(params packages.GetAllPackagesParams, principal *models.Principal) middleware.Responder {
		return packages.NewGetAllPackagesOK()
	})
	api.PackagesGetPackageByNameHandler = packages.GetPackageByNameHandlerFunc(func(params packages.GetPackageByNameParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation packages.GetPackageByName has not yet been implemented")
	})
	api.PackagesUpdatePackageHandler = packages.UpdatePackageHandlerFunc(func(params packages.UpdatePackageParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation packages.UpdatePackage has not yet been implemented")
	})
}
