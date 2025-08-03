package k8s

import (
	"context"
	"errors"
	"fmt"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
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

func (c *ClusterConnection) PatchCertificate(crtName, crtNamespace string, patch []byte) error {
	gvr := schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "certificates"}
	crt, err := c.Client.Resource(gvr).Namespace(crtNamespace).Get(context.Background(), crtName, v1.GetOptions{})

	if crt.GetName() != crtName {
		return err
	}
	_, err = c.Client.Resource(gvr).Namespace(crtNamespace).Patch(context.Background(), crtName, types.StrategicMergePatchType, patch, v1.PatchOptions{})
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

func MarkManuallyTriggered(obj *unstructured.Unstructured) error {
	now := time.Now().UTC().Format(time.RFC3339)

	newCond := map[string]interface{}{
		"type":               "Issuing",
		"status":             "True",
		"reason":             "ManuallyTriggered",
		"message":            "Re-issuance manually requested",
		"lastTransitionTime": now,
	}

	conds, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil {
		return err
	}
	if !found {
		conds = []interface{}{}
	}

	replaced := false
	for i, c := range conds {
		if m, ok := c.(map[string]interface{}); ok {
			if t, _ := m["type"].(string); t == "Issuing" {
				conds[i] = newCond
				replaced = true
				break
			}
		}
	}
	if !replaced {
		conds = append(conds, newCond)
	}

	return unstructured.SetNestedSlice(obj.Object, conds, "status", "conditions")
}
