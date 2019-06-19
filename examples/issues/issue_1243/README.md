# Issue: Kustomize generated files when using -o <folder> option should contain the namespace

This folder contains files describing how to address [Issue 1243](https://github.com/kubernetes-sigs/kustomize/issues/1243)

## Create workplace

First, define a place to reproduce the issue

<!-- @makeWorkplace @basensStagingNs @baseNsBaseNs @defaultNoNs @defaultStagingNs @noNsDefault @test -->
```sh
DEMO_HOME=$(mktemp -d)
CONTENT="https://raw.githubusercontent.com\
/keleustes/kustomize\
/allinone/examples/issues/issue_1243"
```

## Import the files into the workplace

<!-- @installResources @basensStagingNs @baseNsBaseNs @defaultNoNs @defaultStagingNs @noNsDefault @test -->
```sh
mkdir -p $DEMO_HOME/actual
mkdir -p $DEMO_HOME/expected

mkdir -p $DEMO_HOME/base-basens
mkdir -p $DEMO_HOME/base-default
mkdir -p $DEMO_HOME/base-nons

curl -s -o "$DEMO_HOME/base-basens/#1.yaml" \
  "$CONTENT/base-basens/{kustomization,deployment}.yaml"

curl -s -o "$DEMO_HOME/base-nons/#1.yaml" \
  "$CONTENT/base-nons/{kustomization,deployment}.yaml"

curl -s -o "$DEMO_HOME/base-default/#1.yaml" \
  "$CONTENT/base-default/{kustomization,deployment}.yaml"
```

## Case 1: Base NS set to "base" , Overlay NS set to "staging"

### Import expected files

<!-- @importExpectedFiles @basensStagingNs @test -->
```sh
mkdir -p $DEMO_HOME/overlay-basens-stagingns
curl -s -o "$DEMO_HOME/overlay-basens-stagingns/#1.yaml" \
  "$CONTENT/overlay-basens-stagingns/{kustomization,deployment}.yaml"

curl -s -o "$DEMO_HOME/expected/#1.yaml" \
  "$CONTENT/expected/{basens-stagingns}.yaml"

mkdir -p $DEMO_HOME/expected/basens-stagingns
curl -s -o "$DEMO_HOME/expected/basens-stagingns/#1.yaml" \
  "$CONTENT/expected/basens-stagingns/{base_apps_v1beta2_deployment_dply1,staging_apps_v1beta2_deployment_dply1}.yaml"
```

### Build basens-stagingns files

Let's build into a file

<!-- @buildAsFileStaging @basensStagingNs @test -->
```sh
kustomize build $DEMO_HOME/overlay-basens-stagingns -o $DEMO_HOME/actual/basens-stagingns.yaml
```

Let's build into a directory

<!-- @buildAsDirStaging @basensStagingNs @test -->
```sh
mkdir -p $DEMO_HOME/actual/basens-stagingns
kustomize build $DEMO_HOME/overlay-basens-stagingns -o $DEMO_HOME/actual/basens-stagingns/
```

### Verify that the actual output is matching the expected output

<!-- @verifyStaging @basensStagingNs @test -->
```sh
diff -r $DEMO_HOME/actual/basens-stagingns $DEMO_HOME/expected/basens-stagingns
```

## Case 2: Base NS set to "base" , Overlay NS set to "base"

### Import expected files

<!-- @importExpectedFiles @baseNsBaseNs @test -->
```sh
mkdir -p $DEMO_HOME/overlay-basens-basens
curl -s -o "$DEMO_HOME/overlay-basens-basens/#1.yaml" \
  "$CONTENT/overlay-basens-basens/{kustomization,deployment}.yaml"

curl -s -o "$DEMO_HOME/expected/#1.yaml" \
  "$CONTENT/expected/{basens-basens}.yaml"

mkdir -p $DEMO_HOME/expected/basens-basens
curl -s -o "$DEMO_HOME/expected/basens-basens/#1.yaml" \
  "$CONTENT/expected/basens-basens/{apps_v1beta2_deployment_dply1,apps_v1beta2_deployment_dply2}.yaml"
```

### Build basens-basens files

Let's build into a file

<!-- @buildAsFileStaging @baseNsBaseNs @test -->
```sh
kustomize build $DEMO_HOME/overlay-basens-basens -o $DEMO_HOME/actual/basens-basens.yaml
```

Let's build into a directory

<!-- @buildAsDirStaging @baseNsBaseNs @test -->
```sh
mkdir -p $DEMO_HOME/actual/basens-basens
kustomize build $DEMO_HOME/overlay-basens-basens -o $DEMO_HOME/actual/basens-basens/
```

### Verify that the actual output is matching the expected output

<!-- @verifyStaging @baseNsBaseNs @test -->
```sh
diff -r $DEMO_HOME/actual/basens-basens $DEMO_HOME/expected/basens-basens
```

## Case 3: Base NS set to "default" , Overlay NS not set

### Import expected files

<!-- @importExpectedFiles @defaultNoNs @test -->
```sh
mkdir -p $DEMO_HOME/overlay-default-nons
curl -s -o "$DEMO_HOME/overlay-default-nons/#1.yaml" \
  "$CONTENT/overlay-default-nons/{kustomization,deployment}.yaml"

curl -s -o "$DEMO_HOME/expected/#1.yaml" \
  "$CONTENT/expected/{default-nons}.yaml"

mkdir -p $DEMO_HOME/expected/default-nons
curl -s -o "$DEMO_HOME/expected/default-nons/#1.yaml" \
  "$CONTENT/expected/default-nons/{apps_v1beta2_deployment_dply1,apps_v1beta2_deployment_dply2}.yaml"
```

### Build default-nons files

Let's build into a file

<!-- @buildAsFileStaging @defaultNoNs @test -->
```sh
kustomize build $DEMO_HOME/overlay-default-nons -o $DEMO_HOME/actual/default-nons.yaml
```

Let's build into a directory

<!-- @buildAsDirStaging @defaultNoNs @test -->
```sh
mkdir -p $DEMO_HOME/actual/default-nons
kustomize build $DEMO_HOME/overlay-default-nons -o $DEMO_HOME/actual/default-nons/
```

### Verify that the actual output is matching the expected output

<!-- @verifyStaging @defaultNoNs @test -->
```sh
diff -r $DEMO_HOME/actual/default-nons $DEMO_HOME/expected/default-nons
```

## Case 4: Base NS set to "default" , Overlay NS set to staging

### Import expected files

<!-- @importExpectedFiles @defaultStagingNs @test -->
```sh
mkdir -p $DEMO_HOME/overlay-default-stagingns
curl -s -o "$DEMO_HOME/overlay-default-stagingns/#1.yaml" \
  "$CONTENT/overlay-default-stagingns/{kustomization,deployment}.yaml"

curl -s -o "$DEMO_HOME/expected/#1.yaml" \
  "$CONTENT/expected/{default-stagingns}.yaml"

mkdir -p $DEMO_HOME/expected/default-stagingns
curl -s -o "$DEMO_HOME/expected/default-stagingns/#1.yaml" \
  "$CONTENT/expected/default-stagingns/{default_apps_v1beta2_deployment_dply1,staging_apps_v1beta2_deployment_dply1}.yaml"
```

### Build default-stagingns files

Let's build into a file

<!-- @buildAsFileStaging @defaultStagingNs @test -->
```sh
kustomize build $DEMO_HOME/overlay-default-stagingns -o $DEMO_HOME/actual/default-stagingns.yaml
```

Let's build into a directory

<!-- @buildAsDirStaging @defaultStagingNs @test -->
```sh
mkdir -p $DEMO_HOME/actual/default-stagingns
kustomize build $DEMO_HOME/overlay-default-stagingns -o $DEMO_HOME/actual/default-stagingns/
```

### Verify that the actual output is matching the expected output

<!-- @verifyStaging @defaultStagingNs @test -->
```sh
diff -r $DEMO_HOME/actual/default-stagingns $DEMO_HOME/expected/default-stagingns
```

## Case 5: Base NS not set , Overlay NS set to default

### Import expected files

<!-- @importExpectedFiles @noNsDefault @test -->
```sh
mkdir -p $DEMO_HOME/overlay-nons-default
curl -s -o "$DEMO_HOME/overlay-nons-default/#1.yaml" \
  "$CONTENT/overlay-nons-default/{kustomization,deployment}.yaml"

curl -s -o "$DEMO_HOME/expected/#1.yaml" \
  "$CONTENT/expected/{nons-default}.yaml"

mkdir -p $DEMO_HOME/expected/nons-default
curl -s -o "$DEMO_HOME/expected/nons-default/#1.yaml" \
  "$CONTENT/expected/nons-default/{apps_v1beta2_deployment_dply1,apps_v1beta2_deployment_dply2}.yaml"
```

### Build nons-default files

Let's build into a file

<!-- @buildAsFileStaging @noNsDefault @test -->
```sh
kustomize build $DEMO_HOME/overlay-nons-default -o $DEMO_HOME/actual/nons-default.yaml
```

Let's build into a directory

<!-- @buildAsDirStaging @noNsDefault @test -->
```sh
mkdir -p $DEMO_HOME/actual/nons-default
kustomize build $DEMO_HOME/overlay-nons-default -o $DEMO_HOME/actual/nons-default/
```

### Verify that the actual output is matching the expected output

<!-- @verifyStaging @noNsDefault @test -->
```sh
diff -r $DEMO_HOME/actual/nons-default $DEMO_HOME/expected/nons-default
```

