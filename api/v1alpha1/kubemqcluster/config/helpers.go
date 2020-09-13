package config

import (
	apiv1 "k8s.io/api/core/v1"
	"reflect"
)

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
