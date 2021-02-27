package webhooks

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type ValidateRouteWebhook struct {
	Client  client.Client
	Log     logr.Logger
	decoder *admission.Decoder
}

// +kubebuilder:webhook:path=/validate-route,mutating=false,failurePolicy=fail,groups=route.openshift.io,resources=routes,verbs=create;update,versions=v1,name=validateroute.kb.io

func (w *ValidateRouteWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	log := w.Log.WithValues("route", req.Name)

	// Get Namespace
	namespace, err := w.GetNamespace(ctx, log, req)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	allowedIPWhitelist := w.GetAllowedIPWhitelist(log, namespace)
	forbiddenIPWhitelist := w.GetForbiddenIPWhitelist(log, namespace)
	requiredIPWhitelist := w.GetRequiredIPWhitelist(log, namespace)

	// Get Route
	routeIPWhitelist, err := w.GetRouteIPWhitelist(log, req)
	if err != nil {
		log.Error(err, "Error getting Route IP whitelist")
		return admission.Errored(http.StatusBadRequest, err)
	}

	// Check allowed IP whitelist only if declared
	if allowedIPWhitelist != nil {
		for _, routeIP := range routeIPWhitelist {
			if !isStringInSlice(routeIP, allowedIPWhitelist) {
				return admission.Denied(fmt.Sprintf("You are not authorized to use the following IP/CIDR in the Route IP whitelist: %v", routeIP))
			}
		}
	}

	// Check forbidden IP whitelist only if declared
	if forbiddenIPWhitelist != nil {
		for _, routeIP := range routeIPWhitelist {
			if isStringInSlice(routeIP, forbiddenIPWhitelist) {
				return admission.Denied(fmt.Sprintf("You are not authorized to use the following IP/CIDR in the Route IP whitelist: %v", routeIP))
			}
		}
	}

	// Check required IP whitelist only if declared
	if requiredIPWhitelist != nil {
		for _, requiredIP := range requiredIPWhitelist {
			if !isStringInSlice(requiredIP, routeIPWhitelist) {
				return admission.Denied(fmt.Sprintf("You have to include the following IP/CIDR in the Route IP whitelist: %v", requiredIP))
			}
		}
	}

	return admission.Allowed("Valid Route")
}

func (w *ValidateRouteWebhook) InjectDecoder(d *admission.Decoder) error {
	w.decoder = d
	return nil
}

func (w *ValidateRouteWebhook) GetNamespace(ctx context.Context, log logr.Logger, req admission.Request) (*corev1.Namespace, error) {
	namespacedName := types.NamespacedName{
		Name: req.Namespace,
	}
	namespace := &corev1.Namespace{}

	if err := w.Client.Get(ctx, namespacedName, namespace); err != nil {
		log.Error(err, "Unable to get Namespace")
		return nil, err
	}

	return namespace, nil
}

func (w *ValidateRouteWebhook) GetRequiredIPWhitelist(log logr.Logger, namespace *corev1.Namespace) []string {
	if ipWhitelistRaw, found := namespace.ObjectMeta.Annotations[RequiredIPWhitelistAnnotation]; found {
		return strings.Split(ipWhitelistRaw, ",")
	}

	return nil
}

func (w *ValidateRouteWebhook) GetAllowedIPWhitelist(log logr.Logger, namespace *corev1.Namespace) []string {
	if ipWhitelistRaw, found := namespace.ObjectMeta.Annotations[AllowedIPWhitelistAnnotation]; found {
		return strings.Split(ipWhitelistRaw, ",")
	}

	return nil
}

func (w *ValidateRouteWebhook) GetForbiddenIPWhitelist(log logr.Logger, namespace *corev1.Namespace) []string {
	if ipWhitelistRaw, found := namespace.ObjectMeta.Annotations[ForbiddenIPWhitelistAnnotation]; found {
		return strings.Split(ipWhitelistRaw, ",")
	}

	return nil
}

func (w *ValidateRouteWebhook) GetRouteIPWhitelist(log logr.Logger, req admission.Request) ([]string, error) {
	route := &routev1.Route{}
	err := w.decoder.Decode(req, route)
	if err != nil {
		log.Error(err, "Unable to decode request")
		return nil, err
	}

	if whitelistRaw, found := route.ObjectMeta.Annotations[RouteIPWhitelistAnnotation]; found {
		return strings.Split(whitelistRaw, ","), nil
	}

	return []string{}, nil
}
