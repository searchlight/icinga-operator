package v1beta1

// Auto-generated. DO NOT EDIT.
import (
	"github.com/appscode/api/dtypes"
	"github.com/golang/glog"
	"github.com/xeipuuv/gojsonschema"
)

var subscribeRequestSchema *gojsonschema.Schema

func init() {
	var err error
	subscribeRequestSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(`{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "properties": {
    "auto_extend": {
      "type": "boolean"
    },
    "product_id": {
      "type": "string"
    },
    "start_time": {
      "type": "integer"
    }
  },
  "type": "object"
}`))
	if err != nil {
		glog.Fatal(err)
	}
}

func (m *SubscribeRequest) IsValid() (*gojsonschema.Result, error) {
	return subscribeRequestSchema.Validate(gojsonschema.NewGoLoader(m))
}
func (m *SubscribeRequest) IsRequest() {}

func (m *GetQuotaResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *CreateClientTokenResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
func (m *GetSubscriptionResponse) SetStatus(s *dtypes.Status) {
	m.Status = s
}
