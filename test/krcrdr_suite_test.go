package webhook

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	//+kubebuilder:scaffold:imports

	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	"github.com/madflojo/testcerts"
	"github.com/nats-io/nats-server/v2/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	apimachineryruntime "k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/cldmnky/krcrdr/internal/api"
	"github.com/cldmnky/krcrdr/internal/api/auth"
	apiclient "github.com/cldmnky/krcrdr/internal/api/handlers/record/client"
	"github.com/cldmnky/krcrdr/internal/api/store"
	"github.com/cldmnky/krcrdr/internal/api/store/providers/frostdb"
	"github.com/cldmnky/krcrdr/internal/api/store/providers/nats"
	"github.com/cldmnky/krcrdr/internal/recorder"
	"github.com/cldmnky/krcrdr/internal/tracer"
	krcrdrwebhook "github.com/cldmnky/krcrdr/internal/webhook"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var ctx context.Context
var cancel context.CancelFunc
var ns *server.Server
var s store.Store
var traceBuffer bytes.Buffer
var traceExporter tracer.Exporter

func TestKrcrdr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "krcrdr Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.Background())

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
	fa, err := auth.NewFakeAuthenticator()
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
	tenant := auth.Tenant{
		ID:   uuid.NewString(),
		Role: "admin",
	}
	jwt, err := fa.CreateJWSWithClaims([]string{"records:w", "records:r"}, tenant)
	Expect(err).NotTo(HaveOccurred())
	apiClient, err := apiclient.NewApiClient(
		"https://127.0.0.1:8443",
		string(jwt),
		true,
	)
	Expect(err).NotTo(HaveOccurred())
	// Setup tracing
	traceExporter, err = tracer.NewExporter("noop", "127.0.0.1:4317", &traceBuffer)
	Expect(err).NotTo(HaveOccurred())
	traceProvider, err := tracer.NewProvider(ctx, "version", traceExporter)
	Expect(err).NotTo(HaveOccurred())
	go func() {
		defer GinkgoRecover()
		err = tracer.StartTracer(ctx, traceExporter)
		Expect(err).NotTo(HaveOccurred())
	}()

	// Setup recorder
	recorder := recorder.NewRecorder(apiClient, traceProvider.Tracer("recorder"))
	wh.Register("/recorder", &webhook.Admission{
		Handler: &krcrdrwebhook.RecorderWebhook{
			Client:   mgr.GetClient(),
			Decoder:  dec,
			Recorder: recorder,
			Tracer:   traceProvider.Tracer("recorder"),
		},
	})

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

	// Start nats
	dir, err := os.MkdirTemp("", "store")
	Expect(err).NotTo(HaveOccurred())
	natsOpts := &server.Options{
		JetStream: true,
		Debug:     true,
		Host:      "127.0.0.1",
		StoreDir:  dir,
	}
	ns, err = server.NewServer(natsOpts)
	Expect(err).NotTo(HaveOccurred())
	ns.Start()
	Expect(ns.ReadyForConnections(20 * time.Second)).To(BeTrue())

	// Setup the store
	stream, err := nats.NewStream(fmt.Sprintf("nats://%s:%d", natsOpts.Host, natsOpts.Port), traceProvider.Tracer("stream"))
	Expect(err).NotTo(HaveOccurred())
	kv, err := nats.NewKV(fmt.Sprintf("nats://%s:%d", natsOpts.Host, natsOpts.Port), traceProvider.Tracer("kv"))
	Expect(err).NotTo(HaveOccurred())
	index, err := frostdb.NewIndex(kv, traceProvider.Tracer("index"))
	Expect(err).NotTo(HaveOccurred())
	s = store.NewStore(stream, kv, index)
	// Start the indexer
	err = s.StartIndexer(ctx)
	Expect(err).NotTo(HaveOccurred())

	// Create certs for the api server
	cert, key, err := testcerts.GenerateCertsToTempFile("/tmp")
	Expect(err).NotTo(HaveOccurred())
	// get filename from path
	cert = filepath.Base(cert)
	key = filepath.Base(key)

	// Start the API server

	opts := &api.Options{
		Host:          "127.0.0.1",
		Authenticator: fa,
		ApiLogger:     logf.Log.WithName("api"),
		Store:         s,
		CertDir:       "/tmp",
		CertName:      cert,
		KeyName:       key,
		Tracer:        traceProvider.Tracer("api"),
	}

	go func() {
		defer GinkgoRecover()
		err = api.NewServer(*opts).Start(ctx)
		Expect(err).NotTo(HaveOccurred())
	}()

	By("Creating the tenant")
	s.CreateTenant(ctx, &store.Tenant{
		Name: "test",
		Id:   tenant.ID,
	})

	By("Applying intial resources")
	err = k8sClient.Create(ctx, testDeployment())
	Expect(err).ShouldNot(HaveOccurred())

})

var _ = AfterSuite(func() {
	time.Sleep(10 * time.Second)
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
	ns.Shutdown()
	err = traceExporter.Shutdown(ctx)
	Expect(err).ToNot(HaveOccurred())

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
		// loop 100 changes
		for i := 0; i < 10; i++ {
			updatedDeployment.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("updated%d", i)
			updatedDeployment.Spec.Template.Spec.Containers[0].Ports = nil
			updatedDeployment.Spec.Paused = true
			updatedDeployment.ObjectMeta.Labels["addedLabel"] = "test"
			updatedDeployment.ObjectMeta.Labels["changed"] = fmt.Sprintf("yes%d", i)

			err := k8sClient.Update(ctx, updatedDeployment)
			Expect(err).ShouldNot(HaveOccurred())
		}

	})
	It("should handle delete", func() {
		deployment := testDeployment()
		err := k8sClient.Delete(ctx, deployment)
		Expect(err).ShouldNot(HaveOccurred())
		s, err := s.ListTenants(ctx)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(s).To(HaveLen(1))
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
