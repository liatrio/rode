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

package controllers

import (
	"context"
  "bytes"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	rodev1 "github.com/liatrio/rode/api/v1"
	"github.com/liatrio/rode/pkg/attester"
)

// AttesterReconciler reconciles a Attester object
type AttesterReconciler struct {
	client.Client
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Attesters []attester.Attester
}

func (r *AttesterReconciler) ListAttesters() []attester.Attester {
	return r.Attesters
}

var (
    attesterFinalizerName = "attester.finalizers.rode.liatr.io"
)

// +kubebuilder:rbac:groups=rode.liatr.io,resources=attesters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rode.liatr.io,resources=attesters/status,verbs=get;update;patch

func (r *AttesterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("attester", req.NamespacedName)
	opaTrace := false

	log.Info("Reconciling attester")

	att := &rodev1.Attester{}
	err := r.Get(ctx, req.NamespacedName, att)
	if err != nil {
		log.Error(err, "Unable to load attester")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

    err = r.registerFinalizer(log, att)
    if err != nil {
        log.Error(err, "Error registering finalizer")
    }

    if !att.ObjectMeta.DeletionTimestamp.IsZero() && containsFinalizer(att.ObjectMeta.Finalizers, attesterFinalizerName) {
        log.Info("Removing finalizer")

        att.ObjectMeta.Finalizers = removeFinalizer(att.ObjectMeta.Finalizers, attesterFinalizerName)
        err := r.Update(ctx, att)

        return ctrl.Result{}, err
    }

	policy, err := attester.NewPolicy(req.Name, att.Spec.Policy, opaTrace)
	if err != nil {
		log.Error(err, "Unable to create policy")
        // TODO: set status condition to something bad
		return ctrl.Result{}, err
	}

	// TODO: update status based on results of compiling the policy
    condition := rodev1.AttesterCondition{
        Type: rodev1.AttesterConditionReady,
        Status: rodev1.ConditionTrue,
    }

    att.Status.Conditions = append(att.Status.Conditions, condition)

	if err := r.Status().Update(ctx, att); err != nil {
		log.Error(err, "Unable to update Attester status")
		return ctrl.Result{}, err
	}

	// TODO: check if secret already exists before doing this...maybe just load secret

	signerSecret := &corev1.Secret{}
	var signer attester.Signer

    err = r.Get(ctx, client.ObjectKey{
            Namespace: req.Namespace,
            Name: req.Name,
            }, signerSecret)
    if err != nil {

        // If the secret wasn't found then create the secret
        if err.Error() != "Secret \"" + req.Name + "\" not found" {
            log.Error(err, "Unable to get the secret")
            return ctrl.Result{}, err
        }

        log.Info("Couldn't find secret, creating a new one")

        // Create a new signer
        signer, err = attester.NewSigner(req.Name)
		if err != nil {
			log.Error(err, "Unable to create signer")
			return ctrl.Result{}, err
		}
        log.Info("Created the signer successfully")

        var buffer []byte
        buf := bytes.NewBuffer(buffer)

        // signer.Serialize writes the public and private keys to a buffer object buf
        err = signer.Serialize(buf)
        if err != nil {
            log.Error(err, "Unable to read the private key data from the signer")
            return ctrl.Result{}, err
        }

        // buf writes the private and public key to the signerData string
        signerData := buf.String()

        signerSecret = &corev1.Secret{
            ObjectMeta: metav1.ObjectMeta{
                Namespace: req.Namespace,
                Name:      req.Name,
            },
            Data: map[string][]byte{ "keys": []byte(signerData), },
        }

        err = r.Create(ctx, signerSecret)
        if err != nil {
            log.Error(err, "Unable to create signer Secret")
            return ctrl.Result{}, err
        }

        log.Info("Created the signer secret")
    }

    // Secret did exist, or we just created it 

    // TODO: Update the secret with the signer information?




	r.Attesters = append(r.Attesters, attester.NewAttester(req.Name, policy, signer))

	return ctrl.Result{}, nil
}

func (r *AttesterReconciler) registerFinalizer(logger logr.Logger, attester *rodev1.Attester) error {
	if attester.ObjectMeta.DeletionTimestamp.IsZero() && !containsFinalizer(attester.ObjectMeta.Finalizers, attesterFinalizerName) {
		logger.Info("Creating attester finalizer...")
		attester.ObjectMeta.Finalizers = append(attester.ObjectMeta.Finalizers, attesterFinalizerName)

		if err := r.Update(context.Background(), attester); err != nil {
			return err
		}
	}

	return nil
}

func (r *AttesterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rodev1.Attester{}).
		Complete(r)
}
