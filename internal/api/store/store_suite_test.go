package store

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	"github.com/nats-io/nats-server/v2/server"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	// output logs

	RunSpecs(t, "Store Suite")
}

var (
	opts *server.Options
	ns   *server.Server
	//err  error
)

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
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
	var s Store
	BeforeEach(func() {
		stream, err := NewNatsStream(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port))
		Expect(err).NotTo(HaveOccurred())
		kv, err := NewNatsKV(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port))
		Expect(err).NotTo(HaveOccurred())
		s = NewStore(stream, kv)
	})
	It("should create tenants", func() {
		tenant := &Tenant{
			Name: "foo",
		}
		ret, err := s.CreateTenant(context.Background(), "foo", tenant)
		Expect(err).NotTo(HaveOccurred())
		Expect(ret).To(Equal(tenant))
	})
	It("should get tenants", func() {
		tenant := &Tenant{
			Name: "foo",
		}
		ret, err := s.CreateTenant(context.Background(), "bar", tenant)
		Expect(err).NotTo(HaveOccurred())
		Expect(ret).To(Equal(tenant))
		ret, err = s.GetTenant(context.Background(), "bar")
		Expect(err).NotTo(HaveOccurred())
		Expect(ret).To(Equal(tenant))
	})
	// Do not overwrite existing tenants
	It("should not overwrite existing tenants", func() {
		tenant := &Tenant{
			Name: "foo",
		}
		ret, err := s.CreateTenant(context.Background(), "foo", tenant)
		Expect(err).To(HaveOccurred())
		Expect(ret).To(BeNil())
	})

	var _ = Describe("Streams", func() {
		It("should write streams when tenant exists", func() {
			stream, err := NewNatsStream(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port))
			Expect(err).NotTo(HaveOccurred())
			kv, err := NewNatsKV(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port))
			Expect(err).NotTo(HaveOccurred())
			s := NewStore(stream, kv)
			err = s.WriteStream(context.Background(), "foo", "FOO.bar.baz", &api.Record{})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
