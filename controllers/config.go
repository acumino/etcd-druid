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

package controllers

import (
	"flag"

	"github.com/gardener/etcd-druid/controllers/compaction"
	"github.com/gardener/etcd-druid/controllers/custodian"
	"github.com/gardener/etcd-druid/controllers/etcd"
	"github.com/gardener/etcd-druid/controllers/etcdcopybackupstask"
	"github.com/gardener/etcd-druid/controllers/secret"
	"github.com/gardener/etcd-druid/controllers/utils"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

const (
	metricsAddrFlagName                = "metrics-addr"
	enableLeaderElectionFlagName       = "enable-leader-election"
	leaderElectionIDFlagName           = "leader-election-id"
	leaderElectionResourceLockFlagName = "leader-election-resource-lock"
	ignoreOperationAnnotationFlagName  = "ignore-operation-annotation"
	disableLeaseCacheFlagName          = "disable-lease-cache"

	defaultMetricsAddr                = ":8080"
	defaultEnableLeaderElection       = false
	defaultLeaderElectionID           = "druid-leader-election"
	defaultLeaderElectionResourceLock = resourcelock.LeasesResourceLock
	defaultIgnoreOperationAnnotation  = false
	defaultDisableLeaseCache          = false
)

// LeaderElectionConfig defines the configuration for the leader election for the controller manager.
type LeaderElectionConfig struct {
	// EnableLeaderElection specifies whether to enable leader election for controller manager.
	EnableLeaderElection bool
	// LeaderElectionID is the name of the resource that leader election will use for holding the leader lock.
	LeaderElectionID string
	// LeaderElectionResourceLock specifies which resource type to use for leader election.
	LeaderElectionResourceLock string
}

// ManagerConfig defines the configuration for the controller manager.
type ManagerConfig struct {
	// MetricsAddr is the address the metric endpoint binds to.
	MetricsAddr string
	LeaderElectionConfig
	// DisableLeaseCache specifies whether to disable cache for lease.coordination.k8s.io resources.
	DisableLeaseCache bool
	// IgnoreOperationAnnotation specifies whether to ignore or honour the operation annotation on resources to be reconciled.
	IgnoreOperationAnnotation bool
	// EtcdControllerConfig is the configuration required for etcd controller.
	EtcdControllerConfig *etcd.Config
	// CustodianControllerConfig is the configuration required for custodian controller.
	CustodianControllerConfig *custodian.Config
	// CompactionControllerConfig is the configuration required for compaction controller.
	CompactionControllerConfig *compaction.Config
	// EtcdCopyBackupsTaskControllerConfig is the configuration required for etcd-copy-backup-tasks controller.
	EtcdCopyBackupsTaskControllerConfig *etcdcopybackupstask.Config
	// SecretControllerConfig is the configuration required for secret controller.
	SecretControllerConfig *secret.Config
}

// InitFromFlags initializes the controller manager config from the provided CLI flag set.
func InitFromFlags(fs *flag.FlagSet, cfg *ManagerConfig) {
	flag.StringVar(&cfg.MetricsAddr, metricsAddrFlagName, defaultMetricsAddr, ""+
		"The address the metric endpoint binds to.")
	flag.BoolVar(&cfg.EnableLeaderElection, enableLeaderElectionFlagName, defaultEnableLeaderElection,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&cfg.LeaderElectionID, leaderElectionIDFlagName, defaultLeaderElectionID,
		"Name of the resource that leader election will use for holding the leader lock")
	flag.StringVar(&cfg.LeaderElectionResourceLock, leaderElectionResourceLockFlagName, defaultLeaderElectionResourceLock,
		"Specifies which resource type to use for leader election. Supported options are 'endpoints', 'configmaps', 'leases', 'endpointsleases' and 'configmapsleases'.")
	flag.BoolVar(&cfg.DisableLeaseCache, disableLeaseCacheFlagName, defaultDisableLeaseCache,
		"Disable cache for lease.coordination.k8s.io resources.")
	flag.BoolVar(&cfg.IgnoreOperationAnnotation, ignoreOperationAnnotationFlagName, defaultIgnoreOperationAnnotation,
		"Specifies whether to ignore or honour the operation annotation on resources to be reconciled.")

	cfg.EtcdControllerConfig = &etcd.Config{}
	etcd.InitFromFlags(fs, cfg.EtcdControllerConfig)

	cfg.CustodianControllerConfig = &custodian.Config{}
	custodian.InitFromFlags(fs, cfg.CustodianControllerConfig)

	cfg.CompactionControllerConfig = &compaction.Config{}
	compaction.InitFromFlags(fs, cfg.CompactionControllerConfig)

	cfg.EtcdCopyBackupsTaskControllerConfig = &etcdcopybackupstask.Config{}
	etcdcopybackupstask.InitFromFlags(fs, cfg.EtcdCopyBackupsTaskControllerConfig)

	cfg.SecretControllerConfig = &secret.Config{}
	secret.InitFromFlags(fs, cfg.SecretControllerConfig)
}

// Validate validates the controller manager config.
func (cfg *ManagerConfig) Validate() error {
	if err := utils.ShouldBeOneOfAllowedValues("LeaderElectionResourceLock", getAllowedLeaderElectionResourceLocks(), cfg.LeaderElectionResourceLock); err != nil {
		return err
	}
	if err := cfg.EtcdControllerConfig.Validate(); err != nil {
		return err
	}

	if err := cfg.CustodianControllerConfig.Validate(); err != nil {
		return err
	}

	if err := cfg.CompactionControllerConfig.Validate(); err != nil {
		return err
	}

	if err := cfg.EtcdCopyBackupsTaskControllerConfig.Validate(); err != nil {
		return err
	}

	if err := cfg.SecretControllerConfig.Validate(); err != nil {
		return err
	}

	return nil
}

// getAllowedLeaderElectionResourceLocks gives the resource names that can be used for leader election.
// TODO: This should be changed as lease is the default choice now for leader election. There is no need
// to provide other options. Should be handled as part of a different Issue/PR.
func getAllowedLeaderElectionResourceLocks() []string {
	return []string{
		"endpoints",
		"configmaps",
		"leases",
		"endpointsleases",
		"configmapsleases",
	}
}
