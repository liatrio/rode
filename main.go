/*

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
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"strings"

	"github.com/go-logr/logr"

	"github.com/liatrio/rode/pkg/enforcer"

	rodev1alpha1 "github.com/liatrio/rode/api/v1alpha1"
	"github.com/liatrio/rode/controllers"
	"github.com/liatrio/rode/pkg/attester"
	"github.com/liatrio/rode/pkg/aws"
	"github.com/liatrio/rode/pkg/occurrence"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = rodev1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var healthAddr string
	var certDir string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&healthAddr, "health-addr", ":4000", "The address the health endpoint binds to.")
	flag.StringVar(&certDir, "cert-dir", "/certificates", "The path to tls certificates.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		HealthProbeBindAddress: healthAddr,
		CertDir:                certDir,
		LeaderElection:         enableLeaderElection,
		Port:                   9443,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	attesters := &controllers.AttesterReconciler{
		Client:    mgr.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("Attester"),
		Scheme:    mgr.GetScheme(),
		Attesters: make(map[string]attester.Attester),
	}
	if err = attesters.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Attester")
		os.Exit(1)
	}

	awsConfig := aws.NewAWSConfig(ctrl.Log.WithName("aws").WithName("AWSConfig"))

	grafeasTlsConfig, err := grafeasTlsConfig(setupLog)
	if err != nil {
		setupLog.Error(err, "error creating grafeas TLS config")
		os.Exit(1)
	}
	grafeasClient, err := occurrence.NewGrafeasClient(ctrl.Log.WithName("occurrence").WithName("GrafeasClient"), grafeasTlsConfig, os.Getenv("GRAFEAS_ENDPOINT"))
	if err != nil {
		setupLog.Error(err, "error initializing grafeas client")
		os.Exit(1)
	}

	occurrenceCreator := attester.NewAttestWrapper(ctrl.Log.WithName("attester").WithName("AttestWrapper"), grafeasClient, grafeasClient, attesters)

	if err = (&controllers.CollectorReconciler{
		Client:            mgr.GetClient(),
		Log:               ctrl.Log.WithName("controllers").WithName("Collector"),
		Scheme:            mgr.GetScheme(),
		AWSConfig:         awsConfig,
		OccurrenceCreator: occurrenceCreator,
		Workers:           make(map[string]*controllers.CollectorWorker),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Collector")
		os.Exit(1)
	}

	if err = (&controllers.EnforcerReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Enforcer"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Enforcer")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	excludeNS := strings.Split(os.Getenv("EXCLUDED_NAMESPACES"), ",")
	enforcer := enforcer.NewEnforcer(ctrl.Log.WithName("enforcer"), excludeNS, attesters, grafeasClient, mgr.GetClient())

	// TODO: add webhook route

	checker := func(req *http.Request) error {
		return nil
	}

	mgr.AddHealthzCheck("test", checker)
	mgr.AddReadyzCheck("test", checker)
	mgr.GetWebhookServer().Register("/validate-v1-pod", &webhook.Admission{Handler: enforcer})

	// TODO: add occurrences route

	// TODO: setup TLS for endpoints

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func grafeasTlsConfig(log logr.Logger) (*tls.Config, error) {
	clientCert, err := tls.LoadX509KeyPair(os.Getenv("TLS_CLIENT_CERT"), os.Getenv("TLS_CLIENT_KEY"))
	if err != nil {
		log.Error(err, "Unable to load client cert")
		return nil, err
	}

	cf, err := ioutil.ReadFile(os.Getenv("TLS_CA_CERT"))
	if err != nil {
		log.Error(err, "Unable to load CA cert")
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cf)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
		//InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}
