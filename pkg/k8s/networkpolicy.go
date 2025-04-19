package k8s

import (

	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNetworkPolicy(
	name string,
	labels map[string]string,
	annotations map[string]string,
	ingressNamespaceName string,
) *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: v1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []v1.PolicyType{v1.PolicyTypeIngress, v1.PolicyTypeEgress},
			Ingress: []v1.NetworkPolicyIngressRule{
				{
					From: []v1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{},
						},
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"name": ingressNamespaceName,
								},
							},
						},
					},
				},
			},
			Egress: []v1.NetworkPolicyEgressRule{
				{
					To: []v1.NetworkPolicyPeer{
						{
							IPBlock: &v1.IPBlock{
								CIDR: "0.0.0.0/0",
							},
						},
					},
				},
			},
		},
	}
}
