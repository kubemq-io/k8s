package controller
import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	corev1alpha1 "github.com/kubemq-io/k8s/api/v1alpha1"
	ext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)
var (
	scheme    = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(corev1alpha1.AddToScheme(scheme))
	utilruntime.Must(ext.AddToScheme(scheme))
}
