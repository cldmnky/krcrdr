package webhook

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/cldmnky/krcrdr/internal/recorder"
)

// +kubebuilder:webhook:path=/recorder,mutating=false,failurePolicy=ignore,groups="*",resources="*",verbs=create;update;delete,versions="*",name=recorder.blahonga.me,sideEffects=None,admissionReviewVersions=v1

type RecorderWebhook struct {
	Client   client.Client
	Decoder  *admission.Decoder
	Recorder recorder.Recorder
}

func (v *RecorderWebhook) Handle(ctx context.Context, req admission.Request) admission.Response {
	target := &unstructured.Unstructured{}
	object := &unstructured.Unstructured{}

	_ = v.Decoder.DecodeRaw(req.OldObject, target)
	_ = v.Decoder.DecodeRaw(req.Object, object)

	err := v.Recorder.FromAdmissionRequest(target, object, &req.AdmissionRequest)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(v.Recorder.String())

	return admission.Allowed("recorded")
}
