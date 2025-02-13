// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"os"
	"time"

	"github.com/openshift/generic-admission-server/pkg/cmd"
	"github.com/pingcap/tidb-operator/pkg/features"
	"github.com/pingcap/tidb-operator/pkg/version"
	"github.com/pingcap/tidb-operator/pkg/webhook/statefulset"
	"github.com/pingcap/tidb-operator/pkg/webhook/strategy"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"

	// Enable FIPS when necessary
	_ "github.com/pingcap/tidb-operator/pkg/fips"
)

var (
	printVersion         bool
	extraServiceAccounts string
	minResyncDuration    time.Duration
)

func init() {
	klog.InitFlags(nil)
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	// Define the flag "secure-port" to avoid the `flag.Parse()` reporting error
	// TODO: remove this flag after we don't use the lib "github.com/openshift/generic-admission-server"
	flag.Int("secure-port", 6443, "The port on which to serve HTTPS with authentication and authorization. If 0, don't serve HTTPS at all.")
	flag.BoolVar(&printVersion, "V", false, "Show version and quit")
	flag.BoolVar(&printVersion, "version", false, "Show version and quit")
	flag.StringVar(&extraServiceAccounts, "extraServiceAccounts", "", "comma-separated, extra Service Accounts the Webhook should control. The full pattern for each common service account is system:serviceaccount:<namespace>:<serviceaccount-name>")
	flag.DurationVar(&minResyncDuration, "min-resync-duration", 12*time.Hour, "The resync period in reflectors will be random between MinResyncPeriod and 2*MinResyncPeriod.")
	features.DefaultFeatureGate.AddFlag(flag.CommandLine)
}

func main() {

	flag.Parse()

	logs.InitLogs()
	defer logs.FlushLogs()

	if printVersion {
		version.PrintVersionInfo()
		os.Exit(0)
	}
	version.LogVersionInfo()

	flag.CommandLine.VisitAll(func(flag *flag.Flag) {
		klog.V(1).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})
	// We choose a random resync period between MinResyncPeriod and 2 *
	// MinResyncPeriod, so that our pods started at the same time don't list the apiserver simultaneously.

	ns := os.Getenv("NAMESPACE")
	if len(ns) < 1 {
		klog.Fatal("ENV NAMESPACE should be set.")
	}

	statefulSetAdmissionHook := statefulset.NewStatefulSetAdmissionControl()
	strategyAdmissionHook := strategy.NewStrategyAdmissionHook(&strategy.Registry)

	cmd.RunAdmissionServer(statefulSetAdmissionHook, strategyAdmissionHook)
}
