package kubectl

import (
	"context"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (k *Kubectl) StatefulSetPods(ctx context.Context, namespace, name string) ([]core.Pod, error) {
	var list *core.PodList

	err := retry.OnError(retry.DefaultBackoff, func(err error) bool {
		return isConnectionError(err)
	}, func() (err error) {
		list, err = k.client.
			CoreV1().
			Pods(namespace).
			List(ctx, meta.ListOptions{})

		return err
	})

	if err != nil {
		return nil, err
	}

	result := make([]core.Pod, 0)

	for _, item := range list.Items {
		for _, reference := range item.OwnerReferences {
			if reference.Kind == "StatefulSet" && reference.Name == name {
				result = append(result, item)
			}
		}
	}

	return result, err
}

func (k *Kubectl) Pod(ctx context.Context, namespace, name string) (*core.Pod, error) {
	var p *core.Pod

	err := retry.OnError(retry.DefaultBackoff, func(err error) bool {
		return isConnectionError(err)
	}, func() (err error) {
		p, err = k.client.
			CoreV1().
			Pods(namespace).
			Get(ctx, name, meta.GetOptions{})

		return err
	})

	if err != nil {
		return nil, err
	}

	return p, err
}

func (k *Kubectl) StatefulSet(ctx context.Context, namespace, name string) (*apps.StatefulSet, error) {
	var ss *apps.StatefulSet

	err := retry.OnError(retry.DefaultBackoff, func(err error) bool {
		return isConnectionError(err)
	}, func() (err error) {
		ss, err = k.client.
			AppsV1().
			StatefulSets(namespace).
			Get(ctx, name, meta.GetOptions{})

		return
	})

	if err != nil {
		if k8s_errors.IsNotFound(err) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return ss, err
}

func (k *Kubectl) UpdateStatefulSetScale(ctx context.Context, namespace, name string, scale int) error {
	ss, err := k.StatefulSet(ctx, namespace, name)

	if err != nil {
		return err
	}

	i := int32(scale)
	ss.Spec.Replicas = &i

	err = retry.OnError(retry.DefaultBackoff, func(err error) bool {
		return isConnectionError(err)
	}, func() (err error) {
		_, err = k.client.
			AppsV1().
			StatefulSets(namespace).
			Update(ctx, ss, meta.UpdateOptions{})

		return
	})

	if err != nil {
		if k8s_errors.IsNotFound(err) {
			return ErrNotFound
		}

		return err
	}

	return nil
}
