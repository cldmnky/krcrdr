package recorder

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cldmnky/krcrdr/internal/api/handlers/record/api"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/ghodss/yaml"
	"github.com/mitchellh/mapstructure"
	"github.com/sergi/go-diff/diffmatchpatch"
	"go.opentelemetry.io/otel/trace"
	jsonpatch6902 "gomodules.xyz/jsonpatch/v2"

	admissionv1 "k8s.io/api/admission/v1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//go:generate mockery --name RecorderMock
type Recorder interface {
	FromAdmissionRequest(oldObject, newObject *unstructured.Unstructured, req *admissionv1.AdmissionRequest) error
	SendToApiServer(ctx context.Context) error
}

type recorder struct {
	record *api.Record
	client api.ClientInterface
	tracer trace.Tracer
}

func NewRecorder(client api.ClientInterface, tracer trace.Tracer) Recorder {
	return &recorder{
		client: client,
		tracer: tracer,
	}
}

func (r *recorder) FromAdmissionRequest(oldObject, newObject *unstructured.Unstructured, req *admissionv1.AdmissionRequest) error {
	_, span := r.tracer.Start(context.Background(), "record")
	defer span.End()
	record, err := fromAdmissionRequest(oldObject, newObject, req)
	if err != nil {
		span.RecordError(err)
		return err
	}
	r.record = record
	return nil
}

// SendToApiServer sends the record to the API server
func (r *recorder) SendToApiServer(ctx context.Context) error {
	ctx, span := r.tracer.Start(ctx, "send")
	defer span.End()
	if r.record == nil {
		span.RecordError(fmt.Errorf("no record to send"))
		return fmt.Errorf("no record to send")
	}
	resp, err := r.client.AddRecord(ctx, *r.record)
	if err != nil {
		span.RecordError(err)
		return err
	}
	if resp == nil {
		span.RecordError(fmt.Errorf("no response from API server"))
		return fmt.Errorf("no response from API server")
	}
	if resp.StatusCode > 399 {
		span.RecordError(fmt.Errorf("error sending record: %s", resp.Status))
		return fmt.Errorf("error sending record: %s", resp.Status)
	}
	span.AddEvent("sent")
	return nil
}

func (r *recorder) ToYaml() (string, error) {
	y, err := yaml.Marshal(r.record)
	if err != nil {
		return "", err
	}
	return string(y), nil
}

func fromAdmissionRequest(oldObject, newObject *unstructured.Unstructured, req *admissionv1.AdmissionRequest) (*api.Record, error) {
	var (
		patch      []byte
		err        error
		generation int64
		objectMeta metav1.ObjectMeta
	)

	old, _ := json.Marshal(req.OldObject)
	new, _ := json.Marshal(req.Object)

	switch req.Operation {
	case admissionv1.Create:
		// new object
	case admissionv1.Update:
		// update
	case admissionv1.Delete:
		// old object
	case admissionv1.Connect:
		//pass
	default:
		return nil, fmt.Errorf("unknown operation type: %s", req.Operation)
	}

	if req.Operation != admissionv1.Operation(admissionregistrationv1.Delete) {
		generation = newObject.GetGeneration()
		objectMeta = metav1.ObjectMeta{
			Name:              newObject.GetName(),
			Namespace:         newObject.GetNamespace(),
			Labels:            newObject.GetLabels(),
			Annotations:       newObject.GetAnnotations(),
			CreationTimestamp: newObject.GetCreationTimestamp(),
		}
	} else {

		generation = oldObject.GetGeneration()
		objectMeta = metav1.ObjectMeta{
			Name:              oldObject.GetName(),
			Namespace:         oldObject.GetNamespace(),
			Labels:            oldObject.GetLabels(),
			Annotations:       oldObject.GetAnnotations(),
			CreationTimestamp: oldObject.GetCreationTimestamp(),
		}
	}
	patch, err = jsonpatch.CreateMergePatch(old, new)
	if err != nil {
		patch = nil
	}
	p, err := jsonpatch6902.CreatePatch(old, new)
	if err != nil {
		p = nil
	}
	jp6902str := ""
	for _, op := range p {
		jp6902str += fmt.Sprintf("%s\n", op.Json())
	}

	oldY, _ := yaml.JSONToYAML(old)
	newY, _ := yaml.JSONToYAML(new)
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(oldY), string(newY), false)

	// Todo: get cluster name from config
	cluster := "local"

	record := &api.Record{}
	mapstructure.Decode(req, record)
	record.Name = req.Name
	record.Generation = generation
	record.Cluster = cluster
	record.JsonPatch = string(patch)
	record.JsonPatch6902 = jp6902str
	record.DiffString = dmp.DiffPrettyHtml(diffs)
	record.ObjectMeta = api.IoK8sApimachineryPkgApisMetaV1ObjectMeta{
		Name:      &objectMeta.Name,
		Namespace: &objectMeta.Namespace,
		Labels:    &objectMeta.Labels,
	}
	return record, nil
}
