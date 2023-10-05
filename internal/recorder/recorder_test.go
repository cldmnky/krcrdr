package recorder

import (
	"testing"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

func TestFromAdmissionRequest(t *testing.T) {
	oldObject := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "test-old-object-name",
				"uid":  "test-old-object-uid",
			},
		},
	}
	newObject := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "test-object-name",
				"uid":  "test-object-uid",
			},
		},
	}
	req := &admissionv1.AdmissionRequest{
		UID:       types.UID("test-uid"),
		Operation: admissionv1.Create,
		Name:      "test-name",
		Namespace: "test-namespace",
		UserInfo: authenticationv1.UserInfo{
			Username: "test-user",
			UID:      "test-user-uid",
			Groups:   []string{"test-group"},
		},
		Kind: metav1.GroupVersionKind{
			Group:   "test-group",
			Version: "v1",
			Kind:    "test-kind",
		},
		Resource: metav1.GroupVersionResource{
			Group:    "test-group",
			Version:  "v1",
			Resource: "test-resource",
		},
		Object: runtime.RawExtension{
			// name and uid
			Raw: []byte(`{"metadata":{"name":"test-object-name","uid":"test-object-uid"}}`),
			Object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 1,
					Name:       "test-object-name",
					UID:        "test-object-uid",
				},
			},
		},
		OldObject: runtime.RawExtension{
			// name and uid
			Raw: []byte(`{"metadata":{"name":"test-old-object-name","uid":"test-old-object-uid"}}`),
			Object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 0,
					Name:       "test-old-object-name",
					UID:        "test-old-object-uid",
				},
			},
		},
	}

	record, err := fromAdmissionRequest(oldObject, newObject, req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if record.Operation != "CREATE" {
		t.Errorf("unexpected operation: %v", record.Operation)
	}

	if record.Name != "test-name" {
		t.Errorf("unexpected name: %v", record.Name)
	}

	if record.Namespace != "test-namespace" {
		t.Errorf("unexpected namespace: %v", record.Namespace)
	}
	userName := "test-user"
	recordUserName := *record.UserInfo.Username
	if recordUserName != userName {
		t.Errorf("unexpected username: %v", *record.UserInfo.Username)
	}

	if record.Kind.Group != "test-group" {
		t.Errorf("unexpected kind group: %v", record.Kind.Group)
	}

	resourceName := "test-resource"
	if record.Resource.Resource != resourceName {
		t.Errorf("unexpected resource: %v", record.Resource.Resource)
	}

	objMetaName := "test-object-name"
	recordObjMeta := *record.ObjectMeta.Name
	if recordObjMeta != objMetaName {
		t.Errorf("unexpected object name: %v", record.ObjectMeta.Name)
	}

	if record.DiffString == "" {
		t.Errorf("unexpected empty diff string")
	}
}
