package cache

import (
	"context"
	"encoding/json"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	executor "github.com/argoproj/argo-workflows/v3/workflow/artifacts"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewArtifactCache(artifact v1alpha1.Artifact) MemoizationCache {
	return &artifactCache{artifact: artifact}
}

//
//type artEntry struct {
//	NodeID            string
//	Outputs           v1alpha1.Outputs
//	CreationTimestamp v1.Time
//	LastHitTimestamp  v1.Time
//}

type artifactCache struct {
	namespace  string
	kubeClient kubernetes.Interface
	artifact   v1alpha1.Artifact
}

func (a *artifactCache) GetConfigMapKey(ctx context.Context, name, key string) (string, error) {
	//cm, err := a.kubeClient.CoreV1().ConfigMaps(a.namespace).Get(ctx, name, metav1.GetOptions{})
	//if err != nil {
	//	return "", err
	//}
	//return cm.Data[key], nil
	// TODO
	return "", nil
}

func (a *artifactCache) GetSecret(ctx context.Context, name, key string) (string, error) {
	secret, err := a.kubeClient.CoreV1().Secrets(a.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(secret.Data[key]), nil
}

func (a artifactCache) Load(ctx context.Context, key string) (*Entry, error) {

}

func (a artifactCache) Save(ctx context.Context, key string, nodeId string, value *v1alpha1.Outputs) error {

	dv, err := executor.NewDriver(ctx, &a.artifact, &a)
	entry := Entry{
		NodeID:            nodeId,
		Outputs:           value,
		CreationTimestamp: v1.Now(),
		LastHitTimestamp:  v1.Now(),
	}

	file, _ := json.MarshalIndent(entry, "", "  ")
	err = ioutil.WriteFile(key, file, 0644)
	if err != nil {
		return err
	}

	a.artifact.Path = key
	a.artifact.Name = key
	err = dv.Save(key, &a.artifact)
	if err != nil {
		return err
	}

	return nil
}
