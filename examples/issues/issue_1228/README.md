# Issue: Regression in namePrefix configuration in 2.1.0

This folder contains files describing how to address [Issue 1228](https://github.com/kubernetes-sigs/kustomize/issues/1228)

## Setup the workspace

First, define a place to reproduce the issue

<!-- @makeWorkplace @test -->
```sh
REPRODUCE_ISSUE_HOME=$(mktemp -d)
CONTENT="https://raw.githubusercontent.com\
/keleustes/kustomize\
/allinone/examples/issues/issue_1228"
```

## Import the files into the workplace

<!-- @installResources @test -->
```sh
mkdir $REPRODUCE_ISSUE_HOME/actual
mkdir $REPRODUCE_ISSUE_HOME/expected

curl -s -o "$REPRODUCE_ISSUE_HOME/#1.yaml" \
  "$CONTENT/{kustomization,config,deployment}.yaml"

curl -s -o "$REPRODUCE_ISSUE_HOME/expected/#1.yaml" \
  "$CONTENT/expected/{apps_v1_deployment_test-deployment}.yaml"
```

## Build using kustomize

<!-- @build @test -->
```sh
kustomize build $REPRODUCE_ISSUE_HOME -o $REPRODUCE_ISSUE_HOME/actual
```

## Verify that the actual output is matching the expected output

<!-- @verify @test -->
```sh
#diff -r $REPRODUCE_ISSUE_HOME/actual $REPRODUCE_ISSUE_HOME/expected > $REPRODUCE_ISSUE_HOME/diffs.txt
```
