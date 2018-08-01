package restapi

import (
	"fmt"
	"strings"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/projectodd/kwsk/models"
	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/rules"

	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/knative/eventing/pkg/apis/channels/v1alpha1"
	eventing "github.com/knative/eventing/pkg/client/clientset/versioned"
)

func configureRules(api *operations.KwskAPI, eventingClient *eventing.Clientset) {
	api.RulesDeleteRuleHandler = deleteRuleFunc(eventingClient)
	api.RulesGetAllRulesHandler = getAllRulesFunc(eventingClient)
	api.RulesGetRuleByNameHandler = getRuleByNameFunc(eventingClient)
	api.RulesSetStateHandler = rules.SetStateHandlerFunc(func(params rules.SetStateParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation rules.SetState has not yet been implemented")
	})
	api.RulesUpdateRuleHandler = updateRuleFunc(eventingClient)
}

func deleteRuleFunc(eventingClient *eventing.Clientset) rules.DeleteRuleHandlerFunc {
	return func(params rules.DeleteRuleParams, principal *models.Principal) middleware.Responder {
		subName := sanitizeObjectName(params.RuleName)
		namespace := namespaceOrDefault(params.Namespace)
		err := eventingClient.ChannelsV1alpha1().Subscriptions(namespace).Delete(subName, &metav1.DeleteOptions{})
		if err != nil {
			errorMessage := errorMessageFromErr(err)
			if k8sErrors.IsNotFound(err) {
				return rules.NewDeleteRuleNotFound().WithPayload(errorMessage)
			}
			return rules.NewDeleteRuleInternalServerError().WithPayload(errorMessage)
		}

		return rules.NewDeleteRuleOK()
	}
}

func getAllRulesFunc(eventingClient *eventing.Clientset) rules.GetAllRulesHandlerFunc {
	return func(params rules.GetAllRulesParams, principal *models.Principal) middleware.Responder {
		namespace := namespaceOrDefault(params.Namespace)
		subscriptions, err := eventingClient.ChannelsV1alpha1().Subscriptions(namespace).List(metav1.ListOptions{})
		if err != nil {
			return rules.NewGetAllRulesInternalServerError().WithPayload(errorMessageFromErr(err))
		}
		var payload = make([]*models.Rule, len(subscriptions.Items))
		for i, subscription := range subscriptions.Items {
			payload[i] = subscriptionToRule(&subscription)
		}
		return rules.NewGetAllRulesOK().WithPayload(payload)
	}
}

func getRuleByNameFunc(eventingClient *eventing.Clientset) rules.GetRuleByNameHandlerFunc {
	return func(params rules.GetRuleByNameParams, principal *models.Principal) middleware.Responder {
		rule, err := getRuleByName(eventingClient, params.RuleName, params.Namespace)
		if err != nil {
			errorMessage := errorMessageFromErr(err)
			if k8sErrors.IsNotFound(err) {
				return rules.NewGetRuleByNameNotFound().WithPayload(errorMessage)
			}
			return rules.NewGetRuleByNameInternalServerError().WithPayload(errorMessage)
		}
		return rules.NewGetRuleByNameOK().WithPayload(rule)
	}
}

func updateRuleFunc(eventingClient *eventing.Clientset) rules.UpdateRuleHandlerFunc {
	return func(params rules.UpdateRuleParams, principal *models.Principal) middleware.Responder {
		name := params.RuleName
		subName := sanitizeObjectName(name)
		namespace := namespaceOrDefault(params.Namespace)
		version := params.Rule.Version
		if version == "" {
			version = "0.0.1"
		}

		triggerPath := params.Rule.Trigger
		triggerName, _ := splitPathIntoNameAndNamespace(triggerPath)
		channelName := sanitizeObjectName(triggerName)

		actionPath := params.Rule.Action
		actionName, actionNamespace := splitPathIntoNameAndNamespace(actionPath)
		actionNamespace = namespaceOrDefault(actionNamespace)
		serviceName := sanitizeObjectName(actionName)
		serviceHost := fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, actionNamespace)

		annotations := make(map[string]string)
		annotations[KwskName] = name
		annotations[KwskVersion] = version
		annotations["kwsk_trigger_path"] = triggerPath
		annotations["kwsk_trigger_name"] = triggerName
		annotations["kwsk_action_path"] = actionPath
		annotations["kwsk_action_name"] = actionName

		subscription := &v1alpha1.Subscription{
			ObjectMeta: metav1.ObjectMeta{
				Name:        subName,
				Namespace:   namespace,
				Annotations: annotations,
			},
			Spec: v1alpha1.SubscriptionSpec{
				Channel:    channelName,
				Subscriber: serviceHost,
			},
		}

		dbg := fmt.Sprintf("Creating subscription %+v\n", subscription)
		fmt.Printf("%.2000s\n", dbg)
		subscription, err := eventingClient.ChannelsV1alpha1().Subscriptions(namespace).Create(subscription)
		if err != nil {
			fmt.Println("Error updating rule: ", err)
			return rules.NewUpdateRuleInternalServerError().WithPayload(errorMessageFromErr(err))
		}

		rule, err := getRuleByName(eventingClient, subName, namespace)
		if err != nil {
			fmt.Println("Error retrieving updated rule: ", err)
			return rules.NewUpdateRuleInternalServerError().WithPayload(errorMessageFromErr(err))
		}
		return rules.NewUpdateRuleOK().WithPayload(rule)
	}
}

func getRuleByName(eventingClient *eventing.Clientset, name string, namespace string) (*models.Rule, error) {
	subName := sanitizeObjectName(name)
	namespace = namespaceOrDefault(namespace)
	subscription, err := eventingClient.ChannelsV1alpha1().Subscriptions(namespace).Get(subName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return subscriptionToRule(subscription), nil
}

func subscriptionToRule(subscription *v1alpha1.Subscription) *models.Rule {
	objectMeta := subscription.ObjectMeta
	name := objectMeta.Annotations[KwskName]
	version := objectMeta.Annotations[KwskVersion]
	publish := false

	triggerName := objectMeta.Annotations["kwsk_trigger_name"]
	triggerPath := objectMeta.Annotations["kwsk_trigger_path"]
	actionName := objectMeta.Annotations["kwsk_action_name"]
	actionPath := objectMeta.Annotations["kwsk_action_path"]

	return &models.Rule{
		Name:      &name,
		Namespace: &objectMeta.Namespace,
		Version:   &version,
		Publish:   &publish,
		Action: &models.PathName{
			Name: &actionName,
			Path: &actionPath,
		},
		Trigger: &models.PathName{
			Name: &triggerName,
			Path: &triggerPath,
		},
	}
}

func splitPathIntoNameAndNamespace(path string) (string, string) {
	fmt.Printf("Splitting: %s\n", path)
	parts := strings.Split(path, "/")
	fmt.Printf("Parts: %s, length: %d\n", parts, len(parts))
	var name, namespace string
	if len(parts) == 3 {
		namespace = parts[1]
		name = parts[2]
	} else if len(parts) == 2 {
		namespace = parts[0]
		name = parts[1]
	} else if len(parts) == 1 {
		namespace = "_"
		name = parts[0]
	}
	return name, namespace
}
