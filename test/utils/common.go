package utils

import (
	"fmt"

	"github.com/google/go-cmp/cmp/cmpopts"
	fleetv1beta1 "go.goms.io/fleet/apis/placement/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.goms.io/fleet/pkg/utils"
)

// IgnoreConditionLTTAndMessageFields is a cmpopts.IgnoreFields that ignores the LastTransitionTime and Message fields
var IgnoreConditionLTTAndMessageFields = cmpopts.IgnoreFields(metav1.Condition{}, "LastTransitionTime", "Message")

// LessFuncResourceIdentifier is a less function for sorting resource identifiers
var LessFuncResourceIdentifier = func(a, b fleetv1beta1.ResourceIdentifier) bool {
	aStr := fmt.Sprintf(utils.ResourceIdentifierStringFormat, a.Group, a.Version, a.Kind, a.Namespace, a.Name)
	bStr := fmt.Sprintf(utils.ResourceIdentifierStringFormat, b.Group, b.Version, b.Kind, b.Namespace, b.Name)
	return aStr < bStr
}
