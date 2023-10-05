package webhook

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/cldmnky/krcrdr/internal/recorder"
	mockapi "github.com/cldmnky/krcrdr/test/mocks/internal_/api/handlers/record/api"
	"github.com/stretchr/testify/mock"
	admissionv1 "k8s.io/api/admission/v1"
	apimachineryruntime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	ctx = context.Background()
)

func TestRecorderWebhook_Handle(t *testing.T) {
	// mock http client
	apiClient := mockapi.NewClientInterface(t)
	scheme := apimachineryruntime.NewScheme()
	decoder := admission.NewDecoder(scheme)
	recorder := recorder.NewRecorder(apiClient)
	webhook := &RecorderWebhook{
		Decoder:  decoder,
		Recorder: recorder,
		//Client:   apiClient,
	}
	dryRun := true
	// Add a test http.reponse to the mock client
	resp := &http.Response{
		StatusCode: http.StatusOK,
	}
	apiClient.EXPECT().AddRecord(ctx, mock.Anything).Return(resp, nil).Times(3)

	tests := []struct {
		name string
		req  admission.Request
		want admission.Response
	}{
		{
			name: "Test create operation",
			req: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Name:      "test",
					Namespace: "test",
					Operation: admissionv1.Create,
				},
			},
			want: admission.Allowed("recorded"),
		},
		{
			name: "Test update operation",
			req: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Name:      "test",
					Namespace: "test",
					Operation: admissionv1.Update,
				},
			},
			want: admission.Allowed("recorded"),
		},
		{
			name: "Test delete operation",
			req: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Name:      "test",
					Namespace: "test",
					Operation: admissionv1.Delete,
				},
			},
			want: admission.Allowed("recorded"),
		},
		{
			name: "Test dry-run operation",
			req: admission.Request{
				AdmissionRequest: admissionv1.AdmissionRequest{
					Name:      "test",
					Namespace: "test",
					Operation: admissionv1.Delete,
					DryRun:    &dryRun,
				},
			},
			want: admission.Allowed("dry-run"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := webhook.Handle(context.Background(), tt.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecorderWebhook.Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}
