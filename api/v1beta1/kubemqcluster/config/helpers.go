package config

import (
	"reflect"

	"github.com/kubemq-io/k8s/api/v1beta1/kubemqcluster/deployment"
	apiv1 "k8s.io/api/core/v1"
)

// applyServiceExposure applies a connector's shared Service-exposure surface —
// expose, sessionAffinity and pinned node ports — to its Service.
//
// nodePorts maps a Service port NAME to the CRD's pinned value and is honoured only
// when expose is "NodePort": a nodePort on a ClusterIP/LoadBalancer Service is either
// rejected or silently inert, so pinning it there would be a lie. Nil entries leave
// the port kernel-assigned.
func applyServiceExposure(svc *deployment.ServiceConfig, expose, sessionAffinity *string, nodePorts map[string]*int32) {
	if expose != nil {
		svc.SetExpose(*expose)
	}
	if sessionAffinity != nil {
		svc.SetSessionAffinity(*sessionAffinity)
	}
	if expose == nil || *expose != "NodePort" {
		return
	}
	for name, np := range nodePorts {
		if np != nil {
			svc.SetPortNodePort(name, *np)
		}
	}
}

func EqualConfigMaps(a, b *apiv1.ConfigMap) bool {
	if a == nil || b == nil {
		if a == nil && b == nil {
			return true
		} else {
			return false
		}
	}
	return reflect.DeepEqual(a.Data, b.Data) && reflect.DeepEqual(a.BinaryData, b.BinaryData)
}

func EqualSecrets(a, b *apiv1.Secret) bool {
	if a == nil || b == nil {
		if a == nil && b == nil {
			return true
		} else {
			return false
		}
	}
	return reflect.DeepEqual(a.Data, b.Data) && reflect.DeepEqual(a.StringData, b.StringData)
}
