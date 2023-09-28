package webhook

import (
	"context"
	"reflect"
	"testing"

	"github.com/cldmnky/krcrdr/internal/recorder"
	admissionv1 "k8s.io/api/admission/v1"
	apimachineryruntime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestRecorderWebhook_Handle(t *testing.T) {
	scheme := apimachineryruntime.NewScheme()
	decoder := admission.NewDecoder(scheme)
	recorder := recorder.NewRecorder()
	webhook := &RecorderWebhook{
		Decoder:  decoder,
		Recorder: recorder,
	}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := webhook.Handle(context.Background(), tt.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecorderWebhook.Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}
