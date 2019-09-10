// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"testing"

	"sigs.k8s.io/kustomize/v3/pkg/kusttest"
	plugins_test "sigs.k8s.io/kustomize/v3/pkg/plugins/test"
)

func TestKindFilterTransformer(t *testing.T) {
	tc := plugins_test.NewEnvForTest(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"kustomize.config.k8s.io", "v1beta1", "KindFilterTransformer")

	th := kusttest_test.NewKustTestPluginHarness(t, "/app")
	_ = th.LoadAndRunTransformer(`
apiVersion: kustomize.config.k8s.io/v1beta1
kind: KindFilterTransformer
metadata:
  name: notImportantHere
`, ``)

}
