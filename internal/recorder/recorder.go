package recorder

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/ghodss/yaml"
	"github.com/sergi/go-diff/diffmatchpatch"
	jsonpatch6902 "gomodules.xyz/jsonpatch/v2"

	admissionv1 "k8s.io/api/admission/v1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

//go:generate mockery --name RecorderMock
type Recorder interface {
	FromAdmissionRequest(oldObject, newObject *unstructured.Unstructured, req *admissionv1.AdmissionRequest) error
	ToYaml() (string, error)
	OperationType() admissionv1.Operation
	SendToApiServer() error
}

type recorder struct {
	record *record
}

func NewRecorder() Recorder {
	return &recorder{}
}

func (r *recorder) FromAdmissionRequest(oldObject, newObject *unstructured.Unstructured, req *admissionv1.AdmissionRequest) error {

	record, err := fromAdmissionRequest(oldObject, newObject, req)
	if err != nil {
		return err
	}
	r.record = record
	return nil
}

// SendToApiServer sends the record to the API server
func (r *recorder) SendToApiServer() error {
	if r.record == nil {
		return fmt.Errorf("no record to send")
	}
	// Todo: send to API server
	// For now just write to a file
	fileName := fmt.Sprintf("%s-%s-%s.yaml", r.record.Operation, *r.record.Namespace, r.record.Name)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	yamlString, err := r.ToYaml()
	if err != nil {
		return err
	}
	_, err = file.WriteString(yamlString)
	if err != nil {
		return err
	}
	return nil
}
func (r *recorder) ToYaml() (string, error) {
	return r.record.ToYaml()
}

func (r *recorder) OperationType() admissionv1.Operation {
	return r.record.Operation
}

type record struct {
	ChangeTimestamp time.Time                   `json:"changeTimestamp"`
	Operation       admissionv1.Operation       `json:"operation"`
	Cluster         string                      `json:"cluster"`
	Namespace       *string                     `json:"namespace"`
	Name            string                      `json:"name"`
	UID             types.UID                   `json:"uid"`
	UserInfo        authenticationv1.UserInfo   `json:"userInfo"`
	Kind            metav1.GroupVersionKind     `json:"kind"`
	Resource        metav1.GroupVersionResource `json:"resource"`
	Generation      *int64                      `json:"generation"`
	JSONPatch       string                      `json:"jsonPatch"`
	JSONPatch6902   string                      `json:"jsonPatch6902"`
	DiffString      string                      `json:"diffstring"`
	ObjectMeta      metav1.ObjectMeta           `json:"objectMeta"`
}

func fromAdmissionRequest(oldObject, newObject *unstructured.Unstructured, req *admissionv1.AdmissionRequest) (*record, error) {
	var generation int64
	var objectMeta metav1.ObjectMeta
	var uid types.UID
	var patch []byte
	var err error

	old, _ := json.Marshal(req.OldObject)
	new, _ := json.Marshal(req.Object)

	if req.Operation != admissionv1.Operation(admissionregistrationv1.Delete) {
		generation = newObject.GetGeneration()
		objectMeta = metav1.ObjectMeta{
			Name:              newObject.GetName(),
			Namespace:         newObject.GetNamespace(),
			Labels:            newObject.GetLabels(),
			Annotations:       newObject.GetAnnotations(),
			CreationTimestamp: newObject.GetCreationTimestamp(),
		}
		uid = newObject.GetUID()
	} else {

		generation = oldObject.GetGeneration()
		objectMeta = metav1.ObjectMeta{
			Name:              oldObject.GetName(),
			Namespace:         oldObject.GetNamespace(),
			Labels:            oldObject.GetLabels(),
			Annotations:       oldObject.GetAnnotations(),
			CreationTimestamp: oldObject.GetCreationTimestamp(),
		}
		uid = oldObject.GetUID()
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

	ret := &record{
		ChangeTimestamp: time.Now(),
		Operation:       req.Operation,
		Cluster:         cluster,
		Namespace:       &req.Namespace,
		Name:            req.Name,
		UID:             uid,
		UserInfo:        req.UserInfo,
		Kind:            req.Kind,
		Resource:        req.Resource,
		Generation:      &generation,
		ObjectMeta:      objectMeta,
		DiffString:      dmp.DiffPrettyHtml(diffs),
		JSONPatch:       string(patch),
		JSONPatch6902:   jp6902str,
	}
	return ret, nil
}

func (r *record) ToYaml() (string, error) {
	b, err := yaml.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
