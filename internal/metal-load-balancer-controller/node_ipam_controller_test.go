// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package metal_load_balancer_controller

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"
)

var _ = Describe("Node IPAM Controller", func() {
	Context("When the PodCIDR is not populated", func() {
		It("should populate PodCIDR and PodCIDRs", func(ctx SpecContext) {
			node := &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "test-",
				},
				Spec: corev1.NodeSpec{
					PodCIDR:  "",
					PodCIDRs: []string{},
				},
				Status: corev1.NodeStatus{
					Addresses: []corev1.NodeAddress{
						{
							Type:    corev1.NodeInternalIP,
							Address: "1a10:code:ninja::1",
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, node)).Should(Succeed())
			DeferCleanup(k8sClient.Delete, node)

			Eventually(Object(node)).Should(SatisfyAll(
				HaveField("Spec.PodCIDR", Equal("1a10:code:ninja::1/64")),
				HaveField("Spec.PodCIDRs", ContainElement("1a10:code:ninja::1/64")),
			))

		})
	})

	Context("When the PodCIDR is already populated", func() {
		It("should not modify PodCIDR and PodCIDRs", func(ctx SpecContext) {
			node := &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "test-",
				},
				Spec: corev1.NodeSpec{
					PodCIDR: "existing-pod-cidr",
					PodCIDRs: []string{
						"existing-pod-cidr",
					},
				},
				Status: corev1.NodeStatus{
					Addresses: []corev1.NodeAddress{
						{
							Type:    corev1.NodeInternalIP,
							Address: "1a10:code:ninja::1",
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, node)).Should(Succeed())
			DeferCleanup(k8sClient.Delete, node)

			Eventually(Object(node)).Should(SatisfyAll(
				HaveField("Spec.PodCIDR", Equal("existing-pod-cidr")),
				HaveField("Spec.PodCIDRs", ContainElement("existing-pod-cidr")),
			))
		})
	})
})
