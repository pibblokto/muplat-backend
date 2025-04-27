package k8s

import (
	"context"
	"fmt"

	// CRD Go types :contentReference[oaicite:0]{index=0}
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *ClusterConnection) CreatePostgresClusterObject(
	name string,
	namespace string,
	database string,
	labels map[string]string,
	annotations map[string]string,
	storageSize uint,
) *unstructured.Unstructured {

	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "postgresql.cnpg.io/v1",
			"kind":       "Cluster",
			"metadata": map[string]interface{}{
				"name":        name,
				"namespace":   namespace,
				"labels":      labels,
				"annotations": annotations,
			},
			"spec": map[string]interface{}{
				"instances": 1,
				"storage": map[string]interface{}{
					"size": fmt.Sprintf("%dGi", storageSize),
				},
				"bootstrap": map[string]interface{}{
					"initdb": map[string]interface{}{
						"database": database,
						"owner":    database,
					},
				},
			},
		},
	}
}

func (c *ClusterConnection) ApplyPostgresCluster(pc *unstructured.Unstructured) error {
	pcName := pc.GetName()
	pcNamespace := pc.GetNamespace()
	gvr := schema.GroupVersionResource{Group: "postgresql.cnpg.io", Version: "v1", Resource: "clusters"}
	postgresCluster, _ := c.Client.Resource(gvr).Namespace(pcNamespace).Get(context.Background(), pcName, v1.GetOptions{})

	if postgresCluster == nil {
		_, err := c.Client.Resource(gvr).Namespace(pcNamespace).Create(context.Background(), pc, v1.CreateOptions{})
		if err != nil {
			return err
		}
	} else if postgresCluster.GetName() != pcName {
		_, err := c.Client.Resource(gvr).Namespace(pcNamespace).Update(context.Background(), pc, v1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterConnection) DeletePostgresCluster(pcName, pcNamespace string) error {
	gvr := schema.GroupVersionResource{Group: "postgresql.cnpg.io", Version: "v1", Resource: "clusters"}
	postgresCluster, _ := c.Client.Resource(gvr).Namespace(pcNamespace).Get(context.Background(), pcName, v1.GetOptions{})

	if postgresCluster == nil {
		return nil
	}
	err := c.Client.Resource(gvr).Namespace(pcNamespace).Delete(context.Background(), pcName, v1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
