package k8s

import (

	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"context"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *ClusterConnection) CreateNetworkPolicyObject(
	name string,
	namespace string,
	labels map[string]string,
	annotations map[string]string,
) *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
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
									"name": c.ingressNamespace,
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

func (c *ClusterConnection) ApplyNetworkPolicy(np *v1.NetworkPolicy) error {
	networkPolicy, _ := c.Clientset.NetworkingV1().NetworkPolicies(np.Namespace).Get(context.Background(), np.Name, metav1.GetOptions{})
	if networkPolicy.Name != np.Name {
		_, err := c.Clientset.NetworkingV1().NetworkPolicies(np.Namespace).Create(context.Background(), np, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := c.Clientset.NetworkingV1().NetworkPolicies(np.Namespace).Update(context.Background(), np, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterConnection) DeleteNetworkPolicy(npName string, npNamespace string) error {
	networkPolicy, _ := c.Clientset.NetworkingV1().NetworkPolicies(npNamespace).Get(context.Background(), npName, metav1.GetOptions{})
	if networkPolicy.Name != npName {
		return nil
	}
	err := c.Clientset.NetworkingV1().NetworkPolicies(npNamespace).Delete(context.Background(), npName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
