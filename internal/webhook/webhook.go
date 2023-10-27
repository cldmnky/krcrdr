package webhook

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/cldmnky/krcrdr/internal/recorder"
)

var webhooklog = logf.Log.WithName("webhook")

// +kubebuilder:webhook:path=/recorder,mutating=false,failurePolicy=ignore,groups="*",resources="*",verbs=create;update;delete,versions="*",name=recorder.blahonga.me,sideEffects=None,admissionReviewVersions=v1

// RecorderWebhook implements admission.Handler.
type RecorderWebhook struct {
	Client   client.Client
	Decoder  *admission.Decoder
	Recorder recorder.Recorder
	Tracer   trace.Tracer
}

// Handle handles the admission request and records it.
// It decodes the old and new objects from the admission request,
// records the request using the Recorder, and returns an admission response
// indicating that the request was allowed and recorded.
func (v *RecorderWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	ctx, span := v.Tracer.Start(ctx, "Handle")
	defer span.End()
	// Skip dry-run requests
	if req.DryRun != nil && *req.DryRun {
		span.AddEvent("dry-run")
		return admission.Allowed("dry-run")
	}
	target := &unstructured.Unstructured{}
	object := &unstructured.Unstructured{}

	_ = v.Decoder.DecodeRaw(req.OldObject, target)
	_ = v.Decoder.DecodeRaw(req.Object, object)
	err := v.Recorder.FromAdmissionRequest(ctx, target, object, &req.AdmissionRequest)
	if err != nil {
		webhooklog.Error(err, "failed to record request")
		span.RecordError(err)
		return admission.Allowed(fmt.Sprintf("failed to record request: %v", err))
	}
	err = v.Recorder.SendToApiServer(ctx)
	if err != nil {
		webhooklog.Error(err, "failed to send request to API server")
		span.RecordError(err)
		return admission.Allowed(fmt.Sprintf("failed to send request to API server: %v", err))
	}
	span.AddEvent("recorded")
	return admission.Allowed("recorded")
}
