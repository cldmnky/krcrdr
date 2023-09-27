package recorder

import (
	"encoding/json"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/ghodss/yaml"
	jsonpatch6902 "github.com/mattbaird/jsonpatch"
	"github.com/sergi/go-diff/diffmatchpatch"

	admissionv1 "k8s.io/api/admission/v1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

//go:generate mockery --name RecorderMock
type Recorder interface {
	String() string
	FromAdmissionRequest(oldObject, newObject *unstructured.Unstructured, req *admissionv1.AdmissionRequest) error
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

func (r *recorder) String() string {
	return r.record.String()
}

type record struct {
	ChangeTimestamp time.Time                          `json:"changeTimestamp"`
	Operation       admissionv1.Operation              `json:"operation"`
	Cluster         string                             `json:"cluster"`
	Namespace       *string                            `json:"namespace"`
	Name            string                             `json:"name"`
	UID             types.UID                          `json:"uid"`
	UserInfo        authenticationv1.UserInfo          `json:"userInfo"`
	Kind            metav1.GroupVersionKind            `json:"kind"`
	Resource        metav1.GroupVersionResource        `json:"resource"`
	Generation      *int64                             `json:"generation"`
	JSONPatch       string                             `json:"jsonpatch"`
	JSONPatch6902   []jsonpatch6902.JsonPatchOperation `json:"jsonpatch6902"`
	DiffString      string                             `json:"diffstring"`
	ObjectMeta      metav1.ObjectMeta                  `json:"objectmeta"`
}

func fromAdmissionRequest(oldObject, newObject *unstructured.Unstructured, req *admissionv1.AdmissionRequest) (*record, error) {
	var generation int64
	var objectMeta metav1.ObjectMeta
	var uid types.UID
	var patch []byte
	var patch6902 []jsonpatch6902.JsonPatchOperation = nil
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

		patch, _ = jsonpatch.CreateMergePatch(old, new)
		patch6902, err = jsonpatch6902.CreatePatch(old, new)
		if err != nil {
			patch6902 = nil
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
		uid = oldObject.GetUID()
	}
	oldY, _ := yaml.JSONToYAML(old)
	newY, _ := yaml.JSONToYAML(new)
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(oldY), string(newY), false)

	// Todo: get cluster name from config
	cluster := "local"

	return &record{
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
		DiffString:      dmp.DiffPrettyText(diffs),
		JSONPatch:       string(patch),
		JSONPatch6902:   patch6902,
	}, nil
}

func (r *record) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
