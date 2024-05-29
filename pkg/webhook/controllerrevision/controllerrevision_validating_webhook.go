package controllerrevision

import (
	"context"
	"fmt"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"go.goms.io/fleet/pkg/utils"
)

var (
	// ValidationPath is the webhook service path which admission requests are routed to for validating ControllerRevision resources.
	ValidationPath = fmt.Sprintf(utils.ValidationPathFmt, appsv1.SchemeGroupVersion.Group, appsv1.SchemeGroupVersion.Version, "controllerrevision")
)

const (
	deniedControllerRevisionResource  = "controller revision creation is disallowed in the fleet hub cluster"
	allowedControllerRevisionResource = "controller revision creation is allowed in the fleet hub cluster"
	controllerRevisionDeniedFormat    = "controller revision %s/%s creation is disallowed in the fleet hub cluster"
)

type controllerRevisionValidator struct {
	decoder *admission.Decoder
}

// Add registers the webhook for K8s bulit-in object types.
func Add(mgr manager.Manager) error {
	hookServer := mgr.GetWebhookServer()
	hookServer.Register(ValidationPath, &webhook.Admission{Handler: &controllerRevisionValidator{admission.NewDecoder(mgr.GetScheme())}})
	return nil
}

// Handle controllerRevisionValidator denies a pod if it is not created in the system namespaces.
func (v *controllerRevisionValidator) Handle(_ context.Context, req admission.Request) admission.Response {
	namespacedName := types.NamespacedName{Name: req.Name, Namespace: req.Namespace}
	if req.Operation == admissionv1.Create {
		klog.V(2).InfoS("handling controller revision resource", "operation", req.Operation, "subResource", req.SubResource, "namespacedName", namespacedName)
		cr := &appsv1.ControllerRevision{}
		err := v.decoder.Decode(req, cr)
		if err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
		if !utils.IsReservedNamespace(cr.Namespace) {
			klog.V(2).InfoS(deniedControllerRevisionResource, "user", req.UserInfo.Username, "groups", req.UserInfo.Groups, "operation", req.Operation, "GVK", req.RequestKind, "subResource", req.SubResource, "namespacedName", namespacedName)
			return admission.Denied(fmt.Sprintf(controllerRevisionDeniedFormat, cr.Namespace, cr.Name))
		}
	}
	klog.V(3).InfoS(allowedControllerRevisionResource, "user", req.UserInfo.Username, "groups", req.UserInfo.Groups, "operation", req.Operation, "GVK", req.RequestKind, "subResource", req.SubResource, "namespacedName", namespacedName)
	return admission.Allowed("")
}
