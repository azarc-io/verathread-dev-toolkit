package e2e

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestKubernetesSystem(t *testing.T) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not resolve home directory: %v", err)
	}

	cfgPath := path.Join(dirname, ".kube/config")

	config, err := clientcmd.BuildConfigFromFlags(
		"", cfgPath,
	)
	if err != nil {
		t.Fatalf("unable to load kubeconfig from %s: %v", cfgPath, err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("unable to create a client: %v", err)
	}

	t.Run("All system pods should be healthy", func(t *testing.T) {
		pods, err := client.CoreV1().Pods("kube-system").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			t.Fatalf("error getting pods: %v", err)
		}

		for _, item := range pods.Items {
			t.Logf("verifying pod %s: %v", item.Name, item.Status.Phase)
			if len(item.ObjectMeta.OwnerReferences) > 0 &&
				item.ObjectMeta.OwnerReferences[0].Kind == "Job" {
				assert.Equal(t, v1.PodSucceeded, item.Status.Phase)
			} else {
				assert.Equal(t, v1.PodRunning, item.Status.Phase)
			}
		}
	})

	t.Run("All dapr pods should be healthy", func(t *testing.T) {
		pods, err := client.CoreV1().Pods("dapr-system").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			t.Fatalf("error getting pods: %v", err)
		}

		for _, item := range pods.Items {
			t.Logf("verifying pod %s: %v", item.Name, item.Status.Phase)
			if len(item.ObjectMeta.OwnerReferences) > 0 &&
				item.ObjectMeta.OwnerReferences[0].Kind == "Job" {
				assert.Equal(t, v1.PodSucceeded, item.Status.Phase)
			} else {
				assert.Equal(t, v1.PodRunning, item.Status.Phase)
			}
		}
	})

	t.Run("All service pods should be healthy", func(t *testing.T) {
		pods, err := client.CoreV1().Pods("test-dev").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			t.Fatalf("error getting pods: %v", err)
		}

		for _, item := range pods.Items {
			t.Logf("verifying pod %s: %v", item.Name, item.Status.Phase)
			if len(item.ObjectMeta.OwnerReferences) > 0 &&
				item.ObjectMeta.OwnerReferences[0].Kind == "Job" {
				assert.Equal(t, v1.PodSucceeded, item.Status.Phase)
			} else {
				assert.Equal(t, v1.PodRunning, item.Status.Phase)
			}
		}
	})
}
