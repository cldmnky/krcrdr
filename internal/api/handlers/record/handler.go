package record

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	"github.com/cldmnky/krcrdr/internal/api/store"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var logger = logf.Log.WithName("apiHandler")

func Mount(r *gin.Engine, v JWSValidator, store store.Store, tracer trace.Tracer) error {
	r.Use(otelgin.Middleware("api"))
	recordApi := NewRecordHandler(store, tracer)
	apiMw, err := CreateApiMiddleware(v)
	if err != nil {
		return err
	}
	r.Use(apiMw...)
	api.RegisterHandlers(r, recordApi)
	return nil
}

func NewRecordHandler(store store.Store, tracer trace.Tracer) *RecordImpl {
	return &RecordImpl{
		store:  store,
		tracer: tracer,
	}
}

type RecordImpl struct {
	store  store.Store
	tracer trace.Tracer
}

func (r RecordImpl) AddRecord(c *gin.Context) {
	_, span := r.tracer.Start(c, "AddRecord")
	defer span.End()
	// get the post body
	var record api.Record
	if err := c.ShouldBindJSON(&record); err != nil {
		span.RecordError(err)
		c.IndentedJSON(400, err)
		logger.Error(err, "failed to bind json")
		return
	}
	// get the tenant from the context
	tenant := getTenant(c)
	if tenant == nil {
		c.IndentedJSON(400, gin.H{"error": "no tenant in context"})
		return
	}
	storeTennant, err := r.store.GetTenant(c, tenant.ID)
	if err != nil {
		c.IndentedJSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := r.store.WriteStream(c, storeTennant.Id, &record); err != nil {
		c.IndentedJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(200, gin.H{"addRecord": "tenant.ID"})
}

func (r RecordImpl) ListRecords(c *gin.Context) {
	c.IndentedJSON(200, "ListRecords")
}

func getTenant(c *gin.Context) *Tenant {
	t, ok := c.Get("tenant")
	if !ok {
		return nil
	}
	return t.(*Tenant)
}
