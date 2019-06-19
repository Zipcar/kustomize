# Issue: configMapGenerator can't make multiple identically named ConfigMaps with different namespaces

This folder contains files describing how to address [Issue 1155](https://github.com/kubernetes-sigs/kustomize/issues/1155)

## Setup the workspace

First, define a place to reproduce the issue

<!-- @makeWorkplace @test -->
```sh
REPRODUCE_ISSUE_HOME=$(mktemp -d)
CONTENT="https://raw.githubusercontent.com\
/keleustes/kustomize\
/allinone/examples/issues/issue_1155"
```

## Import the files into the workplace

<!-- @installResources @test -->
```sh
mkdir $REPRODUCE_ISSUE_HOME/actual
mkdir $REPRODUCE_ISSUE_HOME/expected

curl -s -o "$REPRODUCE_ISSUE_HOME/#1.yml" \
  "$CONTENT/{kustomization}.yml"

curl -s -o "$REPRODUCE_ISSUE_HOME/expected/#1.yaml" \
  "$CONTENT/expected/{~g_v1_configmap_test-t5t4md8fdm,~g_v1_secret_test-h65t9hg6kc}.yaml"
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
