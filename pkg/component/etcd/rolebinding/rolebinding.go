// Copyright (c) 2023 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rolebinding

import (
	"context"

	"github.com/gardener/gardener/pkg/controllerutils"

	gardenercomponent "github.com/gardener/gardener/pkg/operation/botanist/component"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type component struct {
	client client.Client
	values *Values
}

func (c component) Deploy(ctx context.Context) error {
	roleBinding := c.emptyRoleBinding()
	_, err := controllerutils.GetAndCreateOrStrategicMergePatch(ctx, c.client, roleBinding, func() error {
		roleBinding.Name = c.values.Name
		roleBinding.Namespace = c.values.Namespace
		roleBinding.Labels = c.values.Labels
		roleBinding.OwnerReferences = []metav1.OwnerReference{c.values.OwnerReference}
		roleBinding.RoleRef = rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     c.values.RoleName,
		}
		roleBinding.Subjects = []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      c.values.ServiceAccountName,
				Namespace: c.values.Namespace,
			},
		}
		return nil
	})
	return err
}

func (c component) Destroy(ctx context.Context) error {
	return client.IgnoreNotFound(c.client.Delete(ctx, c.emptyRoleBinding()))
}

func (c *component) emptyRoleBinding() *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.values.Name,
			Namespace: c.values.Namespace,
		},
	}
}

// New creates a new role binding deployer instance.
func New(c client.Client, value *Values) gardenercomponent.Deployer {
	return &component{
		client: c,
		values: value,
	}
}
