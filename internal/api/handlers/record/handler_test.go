// BEGIN: yz9d8f4g5h6j
package record

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
)

func TestApi(t *testing.T) {
	// Create a new gin engine
	r := gin.New()

	// setup the validator
	fa, err := NewFakeAuthenticator()
	require.NoError(t, err)

	// Mount the API on the gin engine
	if err := Mount(r, fa); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create a JWT token with the fake authenticator that has the write and read premissions
	tenant := Tenant{
		ID:   "test",
		Role: "admin",
	}
	wJWT, err := fa.CreateJWSWithClaims([]string{"records:w", "records:r"}, tenant)
	require.NoError(t, err)

	// Create a new HTTP request to the /records endpoint, add bearer authenticaion
	req, err := http.NewRequest("GET", "/record", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bearer := "Bearer " + string(wJWT)
	req.Header.Add("Authorization", bearer)

	// Create a new HTTP response recorder
	w := httptest.NewRecorder()

	// Dispatch the HTTP request
	r.ServeHTTP(w, req)

	// Check the HTTP response status code
	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code: %v", w.Code)
	}
}

func TestRecordImpl_AddRecord(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	// Create a new gin context and RecordImpl instance
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	recordApi := RecordImpl{}

	// add json to the request body
	c.Request, _ = http.NewRequest("POST", "/record", bytes.NewBuffer([]byte(`
	{
		"changeTimestamp": "2023-09-28T18:32:44.73441+02:00",
		"cluster": "local",
		"diffstring": "<span>apiVersion: apps/v1&para;<br>kind: Deployment&para;<br>metadata:&para;<br>  creationTimestamp: &#34;2023-09-28T16:32:44Z&#34;&para;<br>  generation: </span><del style=\"background:#ffe6e6;\">1</del><ins style=\"background:#e6ffe6;\">2</ins><span>&para;<br>  labels:&para;<br>    a</span><ins style=\"background:#e6ffe6;\">ddedLabel: test&para;<br>    a</ins><span>pp: test&para;<br>    changed: &#34;</span><del style=\"background:#ffe6e6;\">no</del><ins style=\"background:#e6ffe6;\">yes</ins><span>&#34;&para;<br>  managedFields:&para;<br>  - apiVersion: apps/v1&para;<br>    fieldsType: FieldsV1&para;<br>    fieldsV1:&para;<br>      f:metadata:&para;<br>        f:labels:&para;<br>          .: {}&para;<br>          f:a</span><ins style=\"background:#e6ffe6;\">ddedLabel: {}&para;<br>          f:a</ins><span>pp: {}&para;<br>          f:changed: {}&para;<br>      f:spec:</span><ins style=\"background:#e6ffe6;\">&para;<br>        f:paused: {}</ins><span>&para;<br>        f:progressDeadlineSeconds: {}&para;<br>        f:replicas: {}&para;<br>        f:revisionHistoryLimit: {}&para;<br>        f:selector: {}&para;<br>        f:strategy:&para;<br>          f:rollingUpdate:&para;<br>            .: {}&para;<br>            f:maxSurge: {}&para;<br>            f:maxUnavailable: {}&para;<br>          f:type: {}&para;<br>        f:template:&para;<br>          f:metadata:&para;<br>            f:labels:&para;<br>              .: {}&para;<br>              f:app: {}&para;<br>          f:spec:&para;<br>            f:containers:&para;<br>              k:{&#34;name&#34;:&#34;test&#34;}:&para;<br>                .: {}&para;<br>                f:image: {}&para;<br>                f:imagePullPolicy: {}&para;<br>                f:name: {}&para;<br>                f:</span><del style=\"background:#ffe6e6;\">ports:&para;<br>                  .: {}&para;<br>                  k:{&#34;containerPort&#34;:80,&#34;protocol&#34;:&#34;TCP&#34;}:&para;<br>                    .: {}&para;<br>                    f:containerPort: {}&para;<br>                    f:protocol: {}&para;<br>                f:</del><span>resources: {}&para;<br>                f:terminationMessagePath: {}&para;<br>                f:terminationMessagePolicy: {}&para;<br>            f:dnsPolicy: {}&para;<br>            f:restartPolicy: {}&para;<br>            f:schedulerName: {}&para;<br>            f:securityContext: {}&para;<br>            f:terminationGracePeriodSeconds: {}&para;<br>    manager: webhook.test&para;<br>    operation: Update&para;<br>    time: &#34;2023-09-28T16:32:44Z&#34;&para;<br>  name: deployment&para;<br>  namespace: default&para;<br>  resourceVersion: &#34;198&#34;&para;<br>  uid: 64cb5f7f-d739-4b63-9f8c-32535c8fd6a6&para;<br>spec:&para;<br>  p</span><ins style=\"background:#e6ffe6;\">aused: true&para;<br>  p</ins><span>rogressDeadlineSeconds: 600&para;<br>  replicas: 3&para;<br>  revisionHistoryLimit: 10&para;<br>  selector:&para;<br>    matchLabels:&para;<br>      app: test&para;<br>  strategy:&para;<br>    rollingUpdate:&para;<br>      maxSurge: 25%&para;<br>      maxUnavailable: 25%&para;<br>    type: RollingUpdate&para;<br>  template:&para;<br>    metadata:&para;<br>      creationTimestamp: null&para;<br>      labels:&para;<br>        app: test&para;<br>    spec:&para;<br>      containers:&para;<br>      - image: </span><ins style=\"background:#e6ffe6;\">upda</ins><span>te</span><del style=\"background:#ffe6e6;\">st</del><ins style=\"background:#e6ffe6;\">d</ins><span>:1</span><del style=\"background:#ffe6e6;\">.</del><span>2</span><del style=\"background:#ffe6e6;\">.</del><span>3</span><ins style=\"background:#e6ffe6;\">4</ins><span>&para;<br>        imagePullPolicy: IfNotPresent&para;<br>        name: test</span><del style=\"background:#ffe6e6;\">&para;<br>        ports:&para;<br>        - containerPort: 80&para;<br>          protocol: TCP</del><span>&para;<br>        resources: {}&para;<br>        terminationMessagePath: /dev/termination-log&para;<br>        terminationMessagePolicy: File&para;<br>      dnsPolicy: ClusterFirst&para;<br>      restartPolicy: Always&para;<br>      schedulerName: default-scheduler&para;<br>      securityContext: {}&para;<br>      terminationGracePeriodSeconds: 30&para;<br>status: {}&para;<br></span>",
		"generation": 2,
		"jsonPatch": "{\"metadata\":{\"generation\":2,\"labels\":{\"addedLabel\":\"test\",\"changed\":\"yes\"},\"managedFields\":[{\"apiVersion\":\"apps/v1\",\"fieldsType\":\"FieldsV1\",\"fieldsV1\":{\"f:metadata\":{\"f:labels\":{\".\":{},\"f:addedLabel\":{},\"f:app\":{},\"f:changed\":{}}},\"f:spec\":{\"f:paused\":{},\"f:progressDeadlineSeconds\":{},\"f:replicas\":{},\"f:revisionHistoryLimit\":{},\"f:selector\":{},\"f:strategy\":{\"f:rollingUpdate\":{\".\":{},\"f:maxSurge\":{},\"f:maxUnavailable\":{}},\"f:type\":{}},\"f:template\":{\"f:metadata\":{\"f:labels\":{\".\":{},\"f:app\":{}}},\"f:spec\":{\"f:containers\":{\"k:{\\\"name\\\":\\\"test\\\"}\":{\".\":{},\"f:image\":{},\"f:imagePullPolicy\":{},\"f:name\":{},\"f:resources\":{},\"f:terminationMessagePath\":{},\"f:terminationMessagePolicy\":{}}},\"f:dnsPolicy\":{},\"f:restartPolicy\":{},\"f:schedulerName\":{},\"f:securityContext\":{},\"f:terminationGracePeriodSeconds\":{}}}}},\"manager\":\"webhook.test\",\"operation\":\"Update\",\"time\":\"2023-09-28T16:32:44Z\"}]},\"spec\":{\"paused\":true,\"template\":{\"spec\":{\"containers\":[{\"image\":\"updated:1234\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"test\",\"resources\":{},\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\"}]}}}}",
		"jsonPatch6902": "{\"op\":\"replace\",\"path\":\"/metadata/generation\",\"value\":2}\n{\"op\":\"add\",\"path\":\"/metadata/labels/addedLabel\",\"value\":\"test\"}\n{\"op\":\"replace\",\"path\":\"/metadata/labels/changed\",\"value\":\"yes\"}\n{\"op\":\"add\",\"path\":\"/metadata/managedFields/0/fieldsV1/f:metadata/f:labels/f:addedLabel\",\"value\":{}}\n{\"op\":\"remove\",\"path\":\"/metadata/managedFields/0/fieldsV1/f:spec/f:template/f:spec/f:containers/k:{\\\"name\\\":\\\"test\\\"}/f:ports\"}\n{\"op\":\"add\",\"path\":\"/metadata/managedFields/0/fieldsV1/f:spec/f:paused\",\"value\":{}}\n{\"op\":\"replace\",\"path\":\"/spec/template/spec/containers/0/image\",\"value\":\"updated:1234\"}\n{\"op\":\"remove\",\"path\":\"/spec/template/spec/containers/0/ports\"}\n{\"op\":\"add\",\"path\":\"/spec/paused\",\"value\":true}\n",
		"kind": {
		  "group": "apps",
		  "kind": "Deployment",
		  "version": "v1"
		},
		"name": "deployment",
		"namespace": "default",
		"objectMeta": {
		  "creationTimestamp": "2023-09-28T16:32:44Z",
		  "labels": {
			"addedLabel": "test",
			"app": "test",
			"changed": "yes"
		  },
		  "name": "deployment",
		  "namespace": "default"
		},
		"operation": "UPDATE",
		"resource": {
		  "group": "apps",
		  "resource": "deployments",
		  "version": "v1"
		},
		"uid": "64cb5f7f-d739-4b63-9f8c-32535c8fd6a6",
		"userInfo": {
		  "groups": [
			"system:masters",
			"system:authenticated"
		  ],
		  "username": "admin"
		}
	  }
	`)))

	// Call the AddRecord function
	recordApi.AddRecord(c)

	// Check that the response status code is 200 and the response body is "AddRecord"
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}

}

func TestRecordImpl_ListRecords(t *testing.T) {
	// Create a new gin context and RecordImpl instance
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Set the tenant in the context
	c.Set("tenant", &Tenant{ID: "test"})

	recordApi := RecordImpl{}

	// Call the ListRecords function
	recordApi.ListRecords(c)

	// Check that the response status code is 200"
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
	}
}
