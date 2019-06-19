# Issue: Set variable from command line

This folder contains files describing how to address [Issue 1113](https://github.com/kubernetes-sigs/kustomize/issues/1113)

## Setup the workspace

First, define a place to reproduce the issue

<!-- @makeWorkplace @test -->
```sh
REPRODUCE_ISSUE_HOME=$(mktemp -d)
CONTENT="https://raw.githubusercontent.com\
/keleustes/kustomize\
/allinone/examples/issues/issue_1113"
```

## Import the files into the workplace

<!-- @installResources @test -->
```sh
mkdir $REPRODUCE_ISSUE_HOME/actual
mkdir $REPRODUCE_ISSUE_HOME/expected

curl -s -o "$REPRODUCE_ISSUE_HOME/#1.yaml" \
  "$CONTENT/{kustomization,deployment,service,values}.yaml"

curl -s -o "$REPRODUCE_ISSUE_HOME/expected/#1.yaml" \
  "$CONTENT/expected/{apps_v1beta2_mydeployment_mysql,~g_v1_myservice_mysql,~g_v1_values_file1}.yaml"
```

## Build using kustomize

<!-- @build @test -->
```sh
kustomize build $REPRODUCE_ISSUE_HOME -o $REPRODUCE_ISSUE_HOME/actual
```

## Verify that the actual output is matching the expected output

<!-- @verify @test -->
```sh
diff -r $REPRODUCE_ISSUE_HOME/actual $REPRODUCE_ISSUE_HOME/expected
```
