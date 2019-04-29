package informers

import (
	astrid_types "github.com/SunSince90/ASTRID-kube/types"
)

func New(what astrid_types.InformerType, namespace string) Informer {

	switch what {
	case astrid_types.Deployments:
		return newDeploymentsInformer(namespace)
	case astrid_types.Services:
		return newServicesInformer(namespace)
	}

	return nil
}
