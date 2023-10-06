package store

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	"github.com/cldmnky/krcrdr/internal/api/store/providers/nats"
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
		stream, err := nats.NewStream(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port))
		Expect(err).NotTo(HaveOccurred())
		kv, err := nats.NewKV(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port))
		Expect(err).NotTo(HaveOccurred())
		s = NewStore(stream, kv)
	})
	It("should create tenants", func() {
		tenant := NewTenant("foo")
		ret, err := s.CreateTenant(context.Background(), tenant)
		Expect(err).NotTo(HaveOccurred())
		Expect(ret).To(Equal(tenant))
		storeTenant, err := s.GetTenant(context.Background(), tenant.Id)
		Expect(err).NotTo(HaveOccurred())
		Expect(tenant).To(Equal(storeTenant))
	})
	It("should get tenants", func() {
		tenant := NewTenant("bar")
		ret, err := s.CreateTenant(context.Background(), tenant)
		Expect(err).NotTo(HaveOccurred())
		Expect(ret).To(Equal(tenant))
		ret, err = s.GetTenant(context.Background(), tenant.Id)
		Expect(err).NotTo(HaveOccurred())
		Expect(ret).To(Equal(tenant))
	})
	// Do not overwrite existing tenants
	It("should not overwrite existing tenants", func() {
		// get a tenant
		existingTenants, err := s.ListTenants(context.Background())
		Expect(err).NotTo(HaveOccurred())
		// create a tenant with the same id
		tenant := &Tenant{
			Name: "foo",
			Id:   existingTenants[0],
		}
		ret, err := s.CreateTenant(context.Background(), tenant)
		Expect(err).To(HaveOccurred())
		Expect(ret).To(BeNil())
	})
	It("should list tenants", func() {
		tenant2 := NewTenant("baz")
		ret, err := s.CreateTenant(context.Background(), tenant2)
		Expect(err).NotTo(HaveOccurred())
		Expect(ret).To(Equal(tenant2))
		tenants, err := s.ListTenants(context.Background())
		Expect(err).NotTo(HaveOccurred())
		// len
		Expect(tenants).To(HaveLen(3))
	})

	var _ = Describe("Streams", func() {
		It("should write streams when tenant exists", func() {
			stream, err := nats.NewStream(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port))
			Expect(err).NotTo(HaveOccurred())
			kv, err := nats.NewKV(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port))
			Expect(err).NotTo(HaveOccurred())
			s := NewStore(stream, kv)
			tenants, err := s.ListTenants(context.Background())
			Expect(err).NotTo(HaveOccurred())
			err = s.WriteStream(context.Background(), tenants[0], &api.Record{
				Name:      "foo",
				Namespace: "bar",
				Cluster:   "baz",
				Kind: struct {
					Group   string `json:"group"`
					Kind    string `json:"kind"`
					Version string `json:"version"`
				}{
					Group:   "apps",
					Kind:    "Deployment",
					Version: "v1",
				},
			})
			Expect(err).NotTo(HaveOccurred())

		})
		It("should not write streams when tenant does not exist", func() {
			stream, err := nats.NewStream(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port))
			Expect(err).NotTo(HaveOccurred())
			kv, err := nats.NewKV(fmt.Sprintf("nats://%s:%d", opts.Host, opts.Port))
			Expect(err).NotTo(HaveOccurred())
			s := NewStore(stream, kv)
			err = s.WriteStream(context.Background(), "doesnotexist", &api.Record{
				Name:      "foo",
				Namespace: "bar",
				Cluster:   "baz",
				Kind: struct {
					Group   string `json:"group"`
					Kind    string `json:"kind"`
					Version string `json:"version"`
				}{
					Group:   "apps",
					Kind:    "Deployment",
					Version: "v1",
				},
			})
			Expect(err).To(HaveOccurred())
		})
	})
})
