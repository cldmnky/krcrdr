package record

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/mitchellh/mapstructure"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
)

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
		return
	}

	c.IndentedJSON(200, gin.H{"addRecord": "tenant.ID"})
}

func (r RecordImpl) ListRecords(c *gin.Context) {
	// Get the tenant from the context
	t, err := getTenantClaimsFromContext(c)
	if err != nil {
		c.IndentedJSON(400, err)
		return
	}
	fmt.Println(t)
	c.IndentedJSON(200, "ListRecords")
}

func getTenantClaimsFromContext(c *gin.Context) (*Tenant, error) {
	ck, ok := c.Get(JWTClaimsContextKey)
	if !ok {
		return nil, ErrClaimsInvalid
	}
	// Get privateClaims from ck
	pc := ck.(jwt.Token).PrivateClaims()
	// Get the tenant claim from the private claims
	t, ok := pc[TenantClaim]
	if !ok {
		return nil, fmt.Errorf("tenant claim not found")
	}
	var tenant Tenant
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &tenant,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	if err := decoder.Decode(t); err != nil {
		return nil, fmt.Errorf("failed to decode tenant claim: %w", err)
	}
	return &tenant, nil
}
