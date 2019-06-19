# Demo: Using complex YAML objects as variables

This folder contains files describing how to address [Issue 1190](https://github.com/kubernetes-sigs/kustomize/issues/1190)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```
DEMO_HOME=$(mktemp -d)
CONTENT="https://raw.githubusercontent.com\
/keleustes/kustomize\
/allinone/examples/issues/issue_1190"
```

## Create the deployments

First, create the `foo` deployment.

<!-- @createFoo @test -->
```
mkdir -p $DEMO_HOME/foo
cat <<'EOF' >$DEMO_HOME/foo/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: foo
  labels:
    app: foo
spec:
  selector:
    matchLabels:
      app: foo
  template:
    metadata:
      labels:
        app: foo
    spec:
      containers:
      - image: alpine
        name: foo
EOF

cat <<'EOF' >$DEMO_HOME/foo/kustomization.yaml
resources:
- deployment.yaml
EOF
```

Next, create the `bar` deployment. Note that this deployment will use a
variable rather than hard-coding the `template`.

<!-- @createBar @test -->
```
mkdir -p $DEMO_HOME/bar
cat <<'EOF' >$DEMO_HOME/bar/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: bar
  labels:
    app: bar
spec:
  selector:
    matchLabels:
      app: bar
  template: $(CUSTOM_TEMPLATE)
EOF

cat <<'EOF' >$DEMO_HOME/bar/kustomization.yaml
resources:
- deployment.yaml
EOF
```

## Create the variable references

Now we just need to hook up the `CUSTOM_TEMPLATE` variable to the field it
refers to.

<!-- @createVar @test -->
```
cat <<'EOF' >$DEMO_HOME/kustomization.yaml
bases:
  - foo
  - bar
namePrefix: inlining-example-
vars:
  - name: CUSTOM_TEMPLATE
    objref:
      kind: Deployment
      name: foo
      apiVersion: apps/v1beta2
    fieldref:
      fieldpath: spec.template
configurations:
  - transformer.yaml
EOF

cat <<'EOF' >$DEMO_HOME/transformer.yaml
varReference:
- kind: Deployment
  path: spec/template
EOF
```

# Import expected output
<!-- @installResources @test -->
```sh
mkdir $DEMO_HOME/expected

curl -s -o "$DEMO_HOME/expected/#1.yaml" \
  "$CONTENT/expected/{apps_v1beta2_deployment_inlining-example-bar,apps_v1beta2_deployment_inlining-example-foo}.yaml"
```

## Build and inspect the results

Build both deployments with the following. Results will be emitted to stdout.
Note that the `bar` deployment's template matches that of the `foo`
deployment's template.

<!-- @build @test -->
```
mkdir $DEMO_HOME/actual
kustomize build $DEMO_HOME -o $DEMO_HOME/actual
```

<!-- @verify @test -->
```
diff -r $DEMO_HOME/actual $DEMO_HOME/expected
```
