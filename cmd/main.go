/*
Copyright 2023 Magnus Bengtsson.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"
	"sync"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	recorderv1beta1 "github.com/cldmnky/krcrdr/api/v1beta1"
	//+kubebuilder:scaffold:imports
)

var (
	scheme               = runtime.NewScheme()
	setupLog             = ctrl.Log.WithName("setup")
	metricsAddr          string
	enableLeaderElection bool
	probeAddr            string
	enablecontroller     bool
	enableApi            bool
	frontend             bool
	apiAddr              string
	apiRemoteAddr        string
	frontendAddr         string
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(recorderv1beta1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	flag.BoolVar(&enablecontroller, "controller", false, "Enable webhook and controller manager")
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&enableApi, "api", false, "Enable backend krcrdr API server")
	flag.StringVar(&apiAddr, "api-bind-address", ":8082", "The address the API endpoint binds to.")
	flag.StringVar(&apiRemoteAddr, "api-remote-address", "http://localhost:8082", "The address of the API endpoint.")
	flag.BoolVar(&frontend, "frontend", false, "Enable frontend krcrdr API server")
	flag.StringVar(&frontendAddr, "frontend-bind-address", ":8083", "The address the frontend endpoint binds to.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if !enablecontroller && !enableApi && !frontend {
		setupLog.Error(nil, "No server enabled")
		os.Exit(1)
	}
	wg := sync.WaitGroup{}
	if enablecontroller {
		// Start controller manager
		wg.Add(1)
		go runManager(&wg)
	}
	if enableApi {
		// Start API server
		wg.Add(1)
		go runApiServer(&wg)
	}
	wg.Wait()
}
