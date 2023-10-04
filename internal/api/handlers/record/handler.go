package record

import (
	"github.com/gin-gonic/gin"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var logger = logf.Log.WithName("apiHandler")

func Mount(r *gin.Engine, v JWSValidator) error {
	var recordApi RecordImpl
	apiMw, err := CreateApiMiddleware(v)
	if err != nil {
		return err
	}
	r.Use(apiMw...)
	api.RegisterHandlers(r, recordApi)
	return nil
}

type RecordImpl struct{}

func (r RecordImpl) AddRecord(c *gin.Context) {
	// get the post body
	var record api.Record
	if err := c.ShouldBindJSON(&record); err != nil {
		c.IndentedJSON(400, err)
		logger.Error(err, "failed to bind json")
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
