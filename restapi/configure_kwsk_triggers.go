package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/projectodd/kwsk/models"
	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/triggers"

	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/knative/eventing/pkg/apis/channels/v1alpha1"
	eventing "github.com/knative/eventing/pkg/client/clientset/versioned"
	serving "github.com/knative/serving/pkg/client/clientset/versioned"
)

func configureTriggers(api *operations.KwskAPI, servingClient *serving.Clientset, eventingClient *eventing.Clientset) {
	api.TriggersDeleteTriggerHandler = triggers.DeleteTriggerHandlerFunc(deleteTriggerFunc(eventingClient))
	api.TriggersFireTriggerHandler = triggers.FireTriggerHandlerFunc(fireTriggerFunc(servingClient, eventingClient))
	api.TriggersGetAllTriggersHandler = triggers.GetAllTriggersHandlerFunc(getAllTriggersFunc(eventingClient))
	api.TriggersGetTriggerByNameHandler = triggers.GetTriggerByNameHandlerFunc(getTriggerByNameFunc(eventingClient))
	api.TriggersUpdateTriggerHandler = triggers.UpdateTriggerHandlerFunc(updateTriggerFunc(eventingClient))
}

func deleteTriggerFunc(eventingClient *eventing.Clientset) triggers.DeleteTriggerHandlerFunc {
	return func(params triggers.DeleteTriggerParams, principal *models.Principal) middleware.Responder {
		channelName := sanitizeObjectName(params.TriggerName)
		namespace := namespaceOrDefault(params.Namespace)
		err := eventingClient.ChannelsV1alpha1().Channels(namespace).Delete(channelName, &metav1.DeleteOptions{})
		if err != nil {
			errorMessage := errorMessageFromErr(err)
			if k8sErrors.IsNotFound(err) {
				return triggers.NewDeleteTriggerNotFound().WithPayload(errorMessage)
			}
			return triggers.NewDeleteTriggerInternalServerError().WithPayload(errorMessage)
		}

		return triggers.NewDeleteTriggerOK()
	}
}

func fireTriggerFunc(servingClient *serving.Clientset, eventingClient *eventing.Clientset) triggers.FireTriggerHandlerFunc {
	return func(params triggers.FireTriggerParams, principal *models.Principal) middleware.Responder {
		channelName := sanitizeObjectName(params.TriggerName)
		namespace := namespaceOrDefault(params.Namespace)
		channel, err := eventingClient.ChannelsV1alpha1().Channels(namespace).Get(channelName, metav1.GetOptions{})
		if err != nil {
			errorMessage := errorMessageFromErr(err)
			if k8sErrors.IsNotFound(err) {
				return triggers.NewFireTriggerNotFound().WithPayload(errorMessage)
			}
			return triggers.NewFireTriggerInternalServerError().WithPayload(errorMessage)
		}
		channel, err = withChannelReady(eventingClient, channel)
		if err != nil {
			fmt.Printf("Error waiting for channel readiness: %s\n", err)
			msg := "Error waiting for channel readiness\n"
			errorMessage := &models.ErrorMessage{
				Error: &msg,
			}
			return triggers.NewFireTriggerInternalServerError().WithPayload(errorMessage)
		}

		err = exposeChannelViaIstio(servingClient, channel)
		if err != nil {
			fmt.Printf("Error exposing channel via Istio: %s\n", err)
			return triggers.NewFireTriggerInternalServerError().WithPayload(errorMessageFromErr(err))
		}
		triggerHost := channel.Status.DomainInternal
		// trigger := channelToTrigger(channel)
		istioHostAndPort := istioHostAndPort()

		url := fmt.Sprintf("http://%s/", istioHostAndPort)
		fmt.Printf("Sending POST to url %s with host %s\n", url, triggerHost)

		body, err := json.Marshal(params.Payload)
		if err != nil {
			fmt.Printf("Error marshaling trigger request body: %s\n", err)
			return triggers.NewFireTriggerInternalServerError()
		}
		fmt.Printf("Request Body: %s\n", body)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("Error creating http request for trigger: %s\n", err)
			return triggers.NewFireTriggerInternalServerError()
		}

		req.Host = triggerHost
		req.Header.Set("Content-Type", "application/json")

		var res *http.Response
		// Wait for the trigger to be ready
		readyTimeout := 5 * time.Minute
		err = wait.PollImmediate(1*time.Second, readyTimeout, func() (bool, error) {
			res, err = http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("Trigger not yet ready: %s\n", err)
				return false, err
			}
			return res.StatusCode != http.StatusNotFound && res.StatusCode != http.StatusServiceUnavailable, nil
		})
		if err != nil {
			fmt.Printf("Error calling trigger http endpoint: %s\n", err)
			return triggers.NewFireTriggerInternalServerError()
		}

		defer res.Body.Close()
		// resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("Response: %+v\n", res)

		if res.StatusCode == http.StatusAccepted {
			activationId, err := newActivationId()
			if err != nil {
				fmt.Printf("Error generating activationId: %s\n", err)
				return triggers.NewFireTriggerInternalServerError().WithPayload(errorMessageFromErr(err))
			}

			payload := &models.ActivationID{
				ActivationID: &activationId,
			}

			return triggers.NewFireTriggerAccepted().WithPayload(payload)
		} else {
			return triggers.NewFireTriggerInternalServerError()
		}
		return triggers.NewFireTriggerInternalServerError()
	}
}

// TODO: This entire function is a hack because I don't want to be
// forced to run kwsk inside Kubernetes for local development
func exposeChannelViaIstio(servingClient *serving.Clientset, channel *v1alpha1.Channel) error {
	virtualServiceName := channel.Status.VirtualService.Name
	namespace := channel.Namespace
	virtualService, err := servingClient.NetworkingV1alpha3().VirtualServices(namespace).Get(virtualServiceName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if len(virtualService.Spec.Gateways) == 0 {
		// TODO: don't hardcode this, obviously
		virtualService.Spec.Gateways = []string{"knative-shared-gateway.knative-serving.svc.cluster.local"}
		_, err := servingClient.NetworkingV1alpha3().VirtualServices(namespace).Update(virtualService)
		if err != nil {
			return err
		}
	}
	return nil
}

func getAllTriggersFunc(eventingClient *eventing.Clientset) triggers.GetAllTriggersHandlerFunc {
	return func(params triggers.GetAllTriggersParams, principal *models.Principal) middleware.Responder {
		namespace := namespaceOrDefault(params.Namespace)
		channels, err := eventingClient.ChannelsV1alpha1().Channels(namespace).List(metav1.ListOptions{})
		if err != nil {
			return triggers.NewGetAllTriggersInternalServerError().WithPayload(errorMessageFromErr(err))
		}
		var payload = make([]*models.Trigger, len(channels.Items))
		for i, channel := range channels.Items {
			payload[i] = channelToTrigger(&channel)
		}
		return triggers.NewGetAllTriggersOK().WithPayload(payload)
	}
}

func getTriggerByNameFunc(eventingClient *eventing.Clientset) triggers.GetTriggerByNameHandlerFunc {
	return func(params triggers.GetTriggerByNameParams, principal *models.Principal) middleware.Responder {
		trigger, err := getTriggerByName(eventingClient, params.TriggerName, params.Namespace)
		if err != nil {
			errorMessage := errorMessageFromErr(err)
			if k8sErrors.IsNotFound(err) {
				return triggers.NewGetTriggerByNameNotFound().WithPayload(errorMessage)
			}
			return triggers.NewGetTriggerByNameInternalServerError().WithPayload(errorMessage)
		}
		return triggers.NewGetTriggerByNameOK().WithPayload(trigger)
	}
}

func updateTriggerFunc(eventingClient *eventing.Clientset) triggers.UpdateTriggerHandlerFunc {
	return func(params triggers.UpdateTriggerParams, principal *models.Principal) middleware.Responder {
		name := params.TriggerName
		channelName := sanitizeObjectName(name)
		namespace := namespaceOrDefault(params.Namespace)
		version := params.Trigger.Version
		if version == "" {
			version = "0.0.1"
		}

		annotations := make(map[string]string)
		annotations[KwskName] = name
		annotations[KwskVersion] = version

		channel := &v1alpha1.Channel{
			ObjectMeta: metav1.ObjectMeta{
				Name:        channelName,
				Namespace:   namespace,
				Annotations: annotations,
			},
			Spec: v1alpha1.ChannelSpec{
				ClusterBus: "stub",
			},
		}

		dbg := fmt.Sprintf("Creating channel %+v\n", channel)
		fmt.Printf("%.2000s\n", dbg)
		channel, err := eventingClient.ChannelsV1alpha1().Channels(namespace).Create(channel)
		if err != nil {
			fmt.Println("Error updating trigger: ", err)
			return triggers.NewUpdateTriggerInternalServerError().WithPayload(errorMessageFromErr(err))
		}

		trigger, err := getTriggerByName(eventingClient, channelName, namespace)
		if err != nil {
			fmt.Println("Error retrieving updated trigger: ", err)
			return triggers.NewUpdateTriggerInternalServerError().WithPayload(errorMessageFromErr(err))
		}
		return triggers.NewUpdateTriggerOK().WithPayload(trigger)
	}
}

func getTriggerByName(eventingClient *eventing.Clientset, name string, namespace string) (*models.Trigger, error) {
	channelName := sanitizeObjectName(name)
	namespace = namespaceOrDefault(namespace)
	channel, err := eventingClient.ChannelsV1alpha1().Channels(namespace).Get(channelName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return channelToTrigger(channel), nil
}

func channelToTrigger(channel *v1alpha1.Channel) *models.Trigger {
	objectMeta := channel.ObjectMeta
	name := objectMeta.Annotations[KwskName]
	if name == "" {
		name = objectMeta.Name
	}
	version := objectMeta.Annotations[KwskVersion]
	publish := false

	return &models.Trigger{
		Name:      &name,
		Namespace: &objectMeta.Namespace,
		Version:   &version,
		Publish:   &publish,
		Limits:    map[string]interface{}{},
	}
}

func withChannelReady(eventingClient *eventing.Clientset, channel *v1alpha1.Channel) (*v1alpha1.Channel, error) {
	readyTimeout := 5 * time.Minute
	if !channelReady(channel) {
		err := wait.Poll(1*time.Second, readyTimeout, func() (bool, error) {
			newChannel, err := eventingClient.ChannelsV1alpha1().Channels(channel.Namespace).Get(channel.Name, metav1.GetOptions{})
			if err != nil {
				if k8sErrors.IsNotFound(err) {
					// If not found then keep trying, assuming it will
					// be found later
					return false, nil
				}
				fmt.Println("Error waiting for channel readiness: ", err)
				return false, err
			}
			channel = newChannel
			return channelReady(channel), nil
		})
		if err != nil {
			fmt.Println("Error waiting for channel readiness: ", err)
			return channel, err
		}
	}
	return channel, nil
}

func channelReady(channel *v1alpha1.Channel) bool {
	if c := channel.Status.GetCondition(v1alpha1.ChannelReady); c != nil {
		return c.Status == corev1.ConditionTrue
	}
	return false
}
