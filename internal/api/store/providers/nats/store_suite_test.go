package nats

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/trace"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/cldmnky/krcrdr/internal/tracer"
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	// output logs

	RunSpecs(t, "Store Suite")
}

var (
	opts          *server.Options
	ns            *server.Server
	traceProvider trace.TracerProvider
)

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	traceExporter, err := tracer.NewExporter("noop", "127.0.0.1:4317", GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	traceProvider, err = tracer.NewProvider(context.Background(), "version", traceExporter)
	// create a temporary directory for the store
	dir, err := os.MkdirTemp("", "store")
	Expect(err).NotTo(HaveOccurred())
	opts = &server.Options{
		JetStream: true,
		Debug:     true,
		Host:      "127.0.0.1",
		// mktmpdir
		StoreDir: dir,
	}
	ns, err = server.NewServer(opts)
	Expect(err).NotTo(HaveOccurred())
	ns.Start()
	Expect(ns.ReadyForConnections(10 * time.Second)).To(BeTrue())
})

var _ = AfterSuite(func() {
	ns.Shutdown()
})

var _ = Describe("Tenants", func() {
	var stream *NatsStore
	var kv *NatsStore
	var err error
	BeforeEach(func() {
		stream, err = NewStream(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port), traceProvider.Tracer("stream"))
		Expect(err).NotTo(HaveOccurred())
		Expect(stream).NotTo(BeNil())
		kv, err = NewKV(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port), traceProvider.Tracer("kv"))
		Expect(err).NotTo(HaveOccurred())
		Expect(kv).NotTo(BeNil())
	})
	It("should create tenants", func() {
	})
	It("should get tenants", func() {

	})
	// Do not overwrite existing tenants
	It("should not overwrite existing tenants", func() {

	})
	It("should list tenants", func() {

	})

	var _ = Describe("Streams", func() {
		It("should write streams when tenant exists", func() {

		})
		It("should not write streams when tenant does not exist", func() {

		})
	})
})
