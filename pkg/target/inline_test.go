// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package target_test

import (
	kusttest_test "sigs.k8s.io/kustomize/v3/pkg/kusttest"
	"testing"
)

func TestEnhancedInlining(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app/whatever")
	th.WriteK("/app/whatever", `
resources:
- cronjob.yaml
- deployment.yaml
- values.yaml

configurations:
- kustomizeconfig.yaml

vars:
- name : Values.shared.spec.env
  objref:
    apiVersion: v1
    kind: Values
    name: shared
  fieldref:
    fieldpath: spec.env
`)
	th.WriteF("/app/whatever/kustomizeconfig.yaml", `
varReference:
- kind: Deployment
  path: spec/template/spec/containers[]/env

- kind: CronJob
  path: spec/jobTemplate/spec/template/spec/containers[]/env
`)
	th.WriteF("/app/whatever/cronjob.yaml", `
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: wordpress-cron
  labels:
    app: wordpress
spec:
  schedule: "*/10 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - image: wordpress:4.8-apache
            name: wordpress
            command:
            - php
            args:
            - /path/to/wp-cron.php
            env: $(Values.shared.spec.env)
          restartPolicy: OnFailure
`)
	th.WriteF("/app/whatever/deployment.yaml", `
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: wordpress
  labels:
    app: wordpress
spec:
  selector:
    matchLabels:
      app: wordpress
  template:
    metadata:
      labels:
        app: wordpress
    spec:
      containers:
      - image: wordpress:4.8-apache
        name: wordpress
        ports:
        - containerPort: 80
          name: wordpress
        env: $(Values.shared.spec.env)
`)
	th.WriteF("/app/whatever/values.yaml", `
apiVersion: v1
kind: Values
metadata:
  name: shared
spec:
  env:
  - name: WORDPRESS_DB_USER
    valueFrom:
      secretKeyRef:
        name: wordpress-db-auth
        key: user
  - name: WORDPRESS_DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: wordpress-db-auth
        key: password
`)
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.AssertActualEqualsExpected(m, `
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  labels:
    app: wordpress
  name: wordpress-cron
spec:
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - args:
            - /path/to/wp-cron.php
            command:
            - php
            env:
            - name: WORDPRESS_DB_USER
              valueFrom:
                secretKeyRef:
                  key: user
                  name: wordpress-db-auth
            - name: WORDPRESS_DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  key: password
                  name: wordpress-db-auth
            image: wordpress:4.8-apache
            name: wordpress
          restartPolicy: OnFailure
  schedule: '*/10 * * * *'
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: wordpress
  name: wordpress
spec:
  selector:
    matchLabels:
      app: wordpress
  template:
    metadata:
      labels:
        app: wordpress
    spec:
      containers:
      - env:
        - name: WORDPRESS_DB_USER
          valueFrom:
            secretKeyRef:
              key: user
              name: wordpress-db-auth
        - name: WORDPRESS_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: wordpress-db-auth
        image: wordpress:4.8-apache
        name: wordpress
        ports:
        - containerPort: 80
          name: wordpress
---
apiVersion: v1
kind: Values
metadata:
  name: shared
spec:
  env:
  - name: WORDPRESS_DB_USER
    valueFrom:
      secretKeyRef:
        key: user
        name: wordpress-db-auth
  - name: WORDPRESS_DB_PASSWORD
    valueFrom:
      secretKeyRef:
        key: password
        name: wordpress-db-auth
`)
}
