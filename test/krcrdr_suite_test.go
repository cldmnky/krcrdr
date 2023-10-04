package webhook

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/cldmnky/krcrdr/internal/recorder"

	"github.com/cldmnky/krcrdr/internal/api"
	"github.com/cldmnky/krcrdr/internal/api/handlers/record"
	apiclient "github.com/cldmnky/krcrdr/internal/api/handlers/record/client"
	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	admissionv1 "k8s.io/api/admission/v1"
	//+kubebuilder:scaffold:imports
	appsv1 "k8s.io/api/apps/v1"
	apimachineryruntime "k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	krcrdrwebhook "github.com/cldmnky/krcrdr/internal/webhook"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var ctx context.Context
var cancel context.CancelFunc

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "krcrdr Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: false,

		// The BinaryAssetsDirectory is only required if you want to run the tests directly
		// without call the makefile target test. If not informed it will look for the
		// default path defined in controller-runtime which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform
		// the tests directly. When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join("..", "bin", "k8s",
			fmt.Sprintf("1.28.0-%s-%s", runtime.GOOS, runtime.GOARCH)),

		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: []string{filepath.Join("..", "config", "webhook")},
		},
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	scheme := apimachineryruntime.NewScheme()
	err = clientgoscheme.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	err = admissionv1.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	// setup auth
	fa, err := record.NewFakeAuthenticator()
	Expect(err).NotTo(HaveOccurred())

	// start webhook server using Manager
	webhookInstallOptions := &testEnv.WebhookInstallOptions
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme,
		WebhookServer: webhook.NewServer(webhook.Options{
			Host:    webhookInstallOptions.LocalServingHost,
			Port:    webhookInstallOptions.LocalServingPort,
			CertDir: webhookInstallOptions.LocalServingCertDir,
		}),
		LeaderElection: false,
		Metrics:        metricsserver.Options{BindAddress: "0"},
	})
	Expect(err).NotTo(HaveOccurred())
	dec := admission.NewDecoder(scheme)
	Expect(dec).NotTo(BeNil())
	wh := mgr.GetWebhookServer()
	Expect(wh).NotTo(BeNil())
	tenant := record.Tenant{
		ID:   "test",
		Role: "admin",
	}
	jwt, err := fa.CreateJWSWithClaims([]string{"records:w", "records:r"}, tenant)
	Expect(err).NotTo(HaveOccurred())
	apiClient, err := apiclient.NewApiClient(
		"http://127.0.0.1:8082",
		string(jwt),
	)
	Expect(err).NotTo(HaveOccurred())
	recorder := recorder.NewRecorder(apiClient)
	wh.Register("/recorder", &webhook.Admission{Handler: &krcrdrwebhook.RecorderWebhook{Client: mgr.GetClient(), Decoder: dec, Recorder: recorder}})

	//+kubebuilder:scaffold:webhook

	go func() {
		defer GinkgoRecover()
		err = mgr.Start(ctx)
		Expect(err).NotTo(HaveOccurred())
	}()

	// wait for the webhook server to get ready
	dialer := &net.Dialer{Timeout: time.Second}
	addrPort := fmt.Sprintf("%s:%d", webhookInstallOptions.LocalServingHost, webhookInstallOptions.LocalServingPort)
	Eventually(func() error {
		conn, err := tls.DialWithDialer(dialer, "tcp", addrPort, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return err
		}
		return conn.Close()
	}).Should(Succeed())

	// Start the API server
	opts := &api.Options{
		Addr:          "127.0.0.1:8082",
		Authenticator: fa,
		ApiLogger:     logf.Log.WithName("api"),
		Env:           "dev",
	}

	go func() {
		defer GinkgoRecover()
		err = api.NewServer(*opts).Run(ctx)
		Expect(err).NotTo(HaveOccurred())
	}()

	By("Applying intial resources")
	err = k8sClient.Create(ctx, testDeployment())
	Expect(err).ShouldNot(HaveOccurred())

})

var _ = AfterSuite(func() {
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("recoder webhook", func() {
	It("should handle update", func() {
		updatedDeployment := testDeployment()
		updatedDeployment.Spec.Template.Spec.Containers[0].Image = "updated:1234"
		updatedDeployment.Spec.Template.Spec.Containers[0].Ports = nil
		updatedDeployment.Spec.Paused = true
		updatedDeployment.ObjectMeta.Labels["addedLabel"] = "test"
		updatedDeployment.ObjectMeta.Labels["changed"] = "yes"

		err := k8sClient.Update(ctx, updatedDeployment)
		Expect(err).ShouldNot(HaveOccurred())

	})
	It("should handle delete", func() {
		deployment := testDeployment()
		err := k8sClient.Delete(ctx, deployment)
		Expect(err).ShouldNot(HaveOccurred())
	})
})

func testDeployment() *appsv1.Deployment {
	d, err := os.ReadFile(filepath.Join("test-data", "deployment.yaml"))
	Expect(err).ToNot(HaveOccurred())

	var deployment appsv1.Deployment

	err = yaml.Unmarshal(d, &deployment)
	Expect(err).ToNot(HaveOccurred())
	return &deployment
}
