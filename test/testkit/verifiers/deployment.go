package verifiers

import (
	"context"
	"fmt"

	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kyma-project/telemetry-manager/test/testkit/periodic"
)

func DeploymentShouldBeReady(ctx context.Context, k8sClient client.Client, name types.NamespacedName) {
	Eventually(func(g Gomega) {
		ready, err := isDeploymentReady(ctx, k8sClient, name)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(ready).To(BeTrue())
	}, periodic.EventuallyTimeout, periodic.DefaultInterval).Should(Succeed())
}

func isDeploymentReady(ctx context.Context, k8sClient client.Client, name types.NamespacedName) (bool, error) {
	var deployment appsv1.Deployment
	err := k8sClient.Get(ctx, name, &deployment)
	if err != nil {
		return false, fmt.Errorf("failed to get deployment: %w", err)
	}
	listOptions := client.ListOptions{
		LabelSelector: labels.SelectorFromSet(deployment.Spec.Selector.MatchLabels),
		Namespace:     name.Namespace,
	}

	return IsPodReady(ctx, k8sClient, listOptions)
}

func DeploymentShouldHaveCorrectPodEnv(ctx context.Context, k8sClient client.Client, name types.NamespacedName, expectedSecretRefName string) {
	Eventually(func(g Gomega) {
		var deployment appsv1.Deployment
		g.Expect(k8sClient.Get(ctx, name, &deployment)).To(Succeed())

		container := deployment.Spec.Template.Spec.Containers[0]
		env := container.EnvFrom[0]

		g.Expect(env.SecretRef.LocalObjectReference.Name).To(Equal(expectedSecretRefName))
		g.Expect(*env.SecretRef.Optional).To(BeTrue())
	}, periodic.EventuallyTimeout, periodic.DefaultInterval).Should(Succeed())
}

func DeploymentShouldHaveCorrectPodMetadata(ctx context.Context, k8sClient client.Client, name types.NamespacedName) {
	Eventually(func(g Gomega) {
		var deployment appsv1.Deployment
		g.Expect(k8sClient.Get(ctx, name, &deployment)).To(Succeed())

		g.Expect(deployment.Spec.Template.ObjectMeta.Labels["sidecar.istio.io/inject"]).To(Equal("false"))
		g.Expect(deployment.Spec.Template.ObjectMeta.Annotations["checksum/config"]).ToNot(BeEmpty())
	}, periodic.EventuallyTimeout, periodic.DefaultInterval).Should(Succeed())
}

func DeploymentShouldHaveCorrectPodPriorityClass(ctx context.Context, k8sClient client.Client, name types.NamespacedName, expectedPriorityClassName string) {
	Eventually(func(g Gomega) {
		var deployment appsv1.Deployment
		g.Expect(k8sClient.Get(ctx, name, &deployment)).To(Succeed())

		g.Expect(deployment.Spec.Template.Spec.PriorityClassName).To(Equal(expectedPriorityClassName))
	}, periodic.EventuallyTimeout, periodic.DefaultInterval).Should(Succeed())
}
