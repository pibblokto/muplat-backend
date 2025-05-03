package k8s

import (
	"context"
	"errors"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *ClusterConnection) CreateCertificateObject(
	name string,
	namespace string,
	domain string,
	labels map[string]string,
	annotations map[string]string,
) *unstructured.Unstructured {

	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cert-manager.io/v1",
			"kind":       "Certificate",
			"metadata": map[string]interface{}{
				"name":        name,
				"namespace":   namespace,
				"labels":      labels,
				"annotations": annotations,
			},
			"spec": map[string]interface{}{
				"secretName": fmt.Sprintf("crt-%s", name),
				"issuerRef": map[string]interface{}{
					"kind": "ClusterIssuer",
					"name": c.ClusterIssuerName,
				},
				"commonName":  domain,
				"dnsNames":    []string{domain},
				"renewBefore": "360h",
			},
		},
	}
}

func (c *ClusterConnection) ApplyCertificate(crt *unstructured.Unstructured) error {
	crtName := crt.GetName()
	crtNamespace := crt.GetNamespace()
	gvr := schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "certificates"}
	certificate, _ := c.Client.Resource(gvr).Namespace(crtNamespace).Get(context.Background(), crtName, v1.GetOptions{})

	if certificate == nil {
		_, err := c.Client.Resource(gvr).Namespace(crtNamespace).Create(context.Background(), crt, v1.CreateOptions{})
		if err != nil {
			return err
		}
	} else if certificate.GetName() != crtName {
		_, err := c.Client.Resource(gvr).Namespace(crtNamespace).Update(context.Background(), crt, v1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterConnection) DeleteCertificate(crtName, crtNamespace string) error {
	gvr := schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "certificates"}
	certificate, _ := c.Client.Resource(gvr).Namespace(crtNamespace).Get(context.Background(), crtName, v1.GetOptions{})

	if certificate == nil {
		return nil
	}
	err := c.Client.Resource(gvr).Namespace(crtNamespace).Delete(context.Background(), crtName, v1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterConnection) GetLiveCertificate(crtName string, crtNamespace string) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "certificates"}
	certificate, _ := c.Client.Resource(gvr).Namespace(crtNamespace).Get(context.Background(), crtName, v1.GetOptions{})
	if certificate == nil {
		return nil, errors.New("no certificate was found")
	}
	return certificate, nil
}
