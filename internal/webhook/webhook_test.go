package webhook

import (
	"context"
	"net/http"
	"os"
	"reflect"
	"testing"

	apiMocks "github.com/cldmnky/krcrdr/test/mocks/internal_/api/handlers/record/api"
	"github.com/stretchr/testify/mock"

	"github.com/cldmnky/krcrdr/internal/recorder"
	"github.com/cldmnky/krcrdr/internal/tracer"
	admissionv1 "k8s.io/api/admission/v1"
	apimachineryruntime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	ctx = context.Background()
)

func TestRecorderWebhook_Handle(t *testing.T) {
	// mock http client
	// Setup tracing
	consoleExporter, err := tracer.NewExporter(string(tracer.ExporterTypeConsole), "", os.Stdout)
	if err != nil {
		t.Errorf("failed to create console exporter: %v", err)
	}
	traceProvider, err := tracer.NewProvider(ctx, "version", consoleExporter)
	if err != nil {
		t.Errorf("failed to create trace provider: %v", err)
	}
	apiClient := apiMocks.NewClientInterface(t)
	scheme := apimachineryruntime.NewScheme()
	decoder := admission.NewDecoder(scheme)
	recorder := recorder.NewRecorder(apiClient, traceProvider.Tracer("recorder"))

	webhook := &RecorderWebhook{
		Decoder:  decoder,
		Recorder: recorder,
		Tracer:   traceProvider.Tracer("recorder"),
	}
	dryRun := true
	// Add a test http.reponse to the mock client
	resp := &http.Response{
		StatusCode: http.StatusOK,
	}
	apiClient.EXPECT().AddRecord(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("api.Record")).Return(resp, nil).Times(3)

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
