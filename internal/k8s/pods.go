package k8s

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetPodsInNamespace fetches all pods in the specified namespace.
func GetPodsInNamespace(clientset *kubernetes.Clientset, namespace string) ([]v1.Pod, error) {
	// List all pods in the given namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list pods: %v", err)
	}
	return pods.Items, nil
}

// DeleteRandomPods deletes the specified number of random pods from the given namespace.
func DeleteRandomPods(clientset *kubernetes.Clientset, namespace string, countPods int) error {
	// Get the list of all pods in the namespace
	pods, err := GetPodsInNamespace(clientset, namespace)
	if err != nil {
		return err
	}
	// LOOP OVER PODS AND PRINT NAMES
	for _, pod := range pods {
		fmt.Println(pod.Name)
	}

	// If there are fewer pods than the requested count, adjust the count
	if len(pods) < countPods {
		countPods = len(pods)
	}

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Randomly select pods and delete them
	for i := 0; i < countPods; i++ {
		randomIndex := rand.Intn(len(pods))
		podToDelete := pods[randomIndex]

		// Delete the selected pod
		err := clientset.CoreV1().Pods(namespace).Delete(context.Background(), podToDelete.Name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("Error deleting pod %s: %v\n", podToDelete.Name, err)
		} else {
			fmt.Printf("Successfully deleted pod: %s\n", podToDelete.Name)
		}

		// Remove the deleted pod from the list to avoid re-selection
		pods = append(pods[:randomIndex], pods[randomIndex+1:]...)
	}

	return nil
}
