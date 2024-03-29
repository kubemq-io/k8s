package v1alpha1

import (
	corev1alpha1 "github.com/kubemq-io/k8s/api/v1alpha1"
	"github.com/kubemq-io/k8s/controller/config"
	"go.uber.org/zap/zapcore"
	ext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
)

var (
	scheme     *runtime.Scheme = runtime.NewScheme()
	testLog                    = ctrl.Log.WithName("objects_tests")
	namespace                  = "default"
	mgr        ctrl.Manager
	testConfig *config.Configuration
)

func TestMain(m *testing.M) {
	var err error
	ctrl.SetLogger(zap.New(zap.UseDevMode(true), zap.StacktraceLevel(zapcore.PanicLevel)))
	mgr, err = ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      ":8090",
		Port:                    9443,
		LeaderElection:          false,
		LeaderElectionNamespace: namespace,
		LeaderElectionID:        "kubemq-operator-lock",
		Namespace:               namespace,
	})
	if err != nil {
		testLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	testConfig = config.NewConfiguration().
		SetNamespace(namespace).
		SetClient(mgr.GetClient()).
		SetReader(mgr.GetAPIReader()).
		SetLog(testLog)
	os.Exit(m.Run())
}
func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(corev1alpha1.AddToScheme(scheme))
	utilruntime.Must(ext.AddToScheme(scheme))
}
