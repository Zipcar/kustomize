# Feature Test for Issue 0519


This folder contains files describing how to address [Issue 0519](https://github.com/kubernetes-sigs/kustomize/issues/0519)
Original kubernetes files have imported from [here](https://github.com/DockbitExamples/kubernetes)

This example is using either:
- skip option to select which components are changed by common transformers.
- multibase/composition to select which the components changed by transformers.

The output ends up beeing the same. 
 
## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}//home/jb447c/src/sigs.k8s.io/kustomize/examples/issues/issue_0519_c
mkdir -p ${DEMO_HOME}/using-composition
mkdir -p ${DEMO_HOME}/using-composition/composite
mkdir -p ${DEMO_HOME}/using-composition/composite/canary
mkdir -p ${DEMO_HOME}/using-composition/composite/production
mkdir -p ${DEMO_HOME}/using-composition/constant
mkdir -p ${DEMO_HOME}/using-composition/variable
mkdir -p ${DEMO_HOME}/using-composition/variable/base
mkdir -p ${DEMO_HOME}/using-composition/variable/canary
mkdir -p ${DEMO_HOME}/using-composition/variable/production
mkdir -p ${DEMO_HOME}/using-skip
mkdir -p ${DEMO_HOME}/using-skip/base
mkdir -p ${DEMO_HOME}/using-skip/base/kustomizeconfig
mkdir -p ${DEMO_HOME}/using-skip/canary
mkdir -p ${DEMO_HOME}/using-skip/production
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-composition/composite/canary/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../constant
- ../../variable/canary
- ./ingress.yaml

patchesStrategicMerge:
# - ./ingress.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-composition/composite/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../constant
- ../../variable/production
- ./ingress.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-composition/constant/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ./namespace.yaml
- ./mycrd.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-composition/variable/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: kubeapp-ns

commonLabels:
  app: kubeapp

resources:
- ./service.yaml
- ./deployment.yaml
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-composition/variable/canary/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  env: canary

nameSuffix: -canary

resources:
- ../base

images:
- name: hack4easy/kubesim_health-amd64
  newTag: 0.1.9
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-composition/variable/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  env: production

nameSuffix: -production

resources:
- ../base

images:
- name: hack4easy/kubesim_health-amd64
  newTag: 0.1.0
EOF
```


### Preparation Step KustomizationFile6

<!-- @createKustomizationFile6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: kubeapp-ns

commonLabels:
  app: kubeapp

resources:
- ./namespace.yaml
- ./mycrd.yaml
- ./ingress.yaml
- ./service.yaml
- ./deployment.yaml

configurations:
- ./kustomizeconfig/mycrd.yaml
- ./kustomizeconfig/ingress.yaml
- ./kustomizeconfig/namespace.yaml
- ./kustomizeconfig/customresourcedefinition.yaml
EOF
```


### Preparation Step KustomizationFile7

<!-- @createKustomizationFile7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/canary/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  env: canary

nameSuffix: -canary

resources:
- ../base

patchesStrategicMerge:
- ./ingress.yaml

images:
- name: hack4easy/kubesim_health-amd64
  newTag: 0.1.9
EOF
```


### Preparation Step KustomizationFile8

<!-- @createKustomizationFile8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  env: production

nameSuffix: -production

resources:
- ../base

patchesStrategicMerge:
- ./ingress.yaml

images:
- name: hack4easy/kubesim_health-amd64
  newTag: 0.1.0
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-composition/composite/canary/ingress.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  backend:
    serviceName: kubeapp-production
    servicePort: 80
  rules:
  - host: canary.foo.bar
    http:
      paths:
      - backend:
          serviceName: kubeapp-canary
          servicePort: 80
  - host: foo.bar
    http:
      paths:
      - backend:
          serviceName: kubeapp-production
          servicePort: 80
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-composition/composite/production/ingress.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  backend:
    serviceName: kubeapp-production
    servicePort: 80
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-composition/constant/mycrd.yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: mycrds.my.org
spec:
  additionalPrinterColumns:
  group: my.org
  version: v1alpha1
  names:
    kind: MyCRD
    plural: mycrds
    shortNames:
    - mycrd
  scope: Cluster
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          type: object
          properties:
            simpletext:
              type: string
            replica:
              type: integer
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  name: my-crd
spec:
  simpletext: some simple text
  replica: 123
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-composition/constant/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: kubeapp-ns
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-composition/variable/base/deployment.yaml
kind: Deployment
apiVersion: extensions/v1beta1
metadata:
  name: kubeapp
spec:
  replicas: 1
  template:
    metadata:
      name: kubeapp
      labels:
        app: kubeapp
    spec:
      containers:
      - name: kubeapp
        image: hack4easy/kubesim_health-amd64:latest
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
        ports:
        - name: kubeapp
          containerPort: 8081
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-composition/variable/base/service.yaml
kind: Service
apiVersion: v1
metadata:
  name: kubeapp
spec:
  type: LoadBalancer
  ports:
  - name: http
    port: 80
    targetPort: 8081
    protocol: TCP
  selector:
    app: kubeapp
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/base/deployment.yaml
kind: Deployment
apiVersion: extensions/v1beta1
metadata:
  name: kubeapp
spec:
  replicas: 1
  template:
    metadata:
      name: kubeapp
      labels:
        app: kubeapp
    spec:
      containers:
      - name: kubeapp
        image: hack4easy/kubesim_health-amd64:latest
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
        ports:
        - name: kubeapp
          containerPort: 8081
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/base/ingress.yaml
kind: Ingress
apiVersion: extensions/v1beta1
metadata:
  name: kubeapp
spec:
  backend:
    serviceName: kubeapp
    servicePort: 80
EOF
```


### Preparation Step Resource8

<!-- @createResource8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/base/kustomizeconfig/customresourcedefinition.yaml
commonLabels:
- path: metadata/labels
  version: v1beta1
  group: apiextensions.k8s.io
  kind: CustomResourceDefinition
  skip: true
EOF
```


### Preparation Step Resource9

<!-- @createResource9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/base/kustomizeconfig/ingress.yaml
namePrefix:
- path: metadata/name
  group: extensions
  version: v1beta1
  kind: Ingress
  skip: true

# Not implemented yet
commonLabels:
- path: metadata/labels
  group: extensions
  version: v1beta1
  kind: Ingress
  skip: true
EOF
```


### Preparation Step Resource10

<!-- @createResource10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/base/kustomizeconfig/mycrd.yaml
namespace:
- path: metadata/namespace
  version: v1alpha1
  group: my.org
  kind: MyCRD
  skip: true

namePrefix:
- path: metadata/name
  version: v1alpha1
  group: my.org
  kind: MyCRD
  skip: true

commonLabels:
- path: metadata/labels
  version: v1alpha1
  group: my.org
  kind: MyCRD
  skip: true

EOF
```


### Preparation Step Resource11

<!-- @createResource11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/base/kustomizeconfig/namespace.yaml
namePrefix:
- path: metadata/name
  version: v1
  kind: Namespace
  skip: true

# Not implemented yet
commonLabels:
- path: metadata/labels
  version: v1
  kind: Namespace
  skip: true
EOF
```


### Preparation Step Resource12

<!-- @createResource12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/base/mycrd.yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: mycrds.my.org
spec:
  additionalPrinterColumns:
  group: my.org
  version: v1alpha1
  names:
    kind: MyCRD
    plural: mycrds
    shortNames:
    - mycrd
  scope: Cluster
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          type: object
          properties:
            simpletext:
              type: string
            replica:
              type: integer
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  name: my-crd
spec:
  simpletext: some simple text
  replica: 123
EOF
```


### Preparation Step Resource13

<!-- @createResource13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/base/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: kubeapp-ns
EOF
```


### Preparation Step Resource14

<!-- @createResource14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/base/service.yaml
kind: Service
apiVersion: v1
metadata:
  name: kubeapp
spec:
  type: LoadBalancer
  ports:
  - name: http
    port: 80
    targetPort: 8081
    protocol: TCP
  selector:
    app: kubeapp
EOF
```


### Preparation Step Resource15

<!-- @createResource15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/canary/ingress.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  backend:
    serviceName: kubeapp-production
    servicePort: 80
  rules:
  - host: canary.foo.bar
    http:
      paths:
      - backend:
          serviceName: kubeapp-canary
          servicePort: 80
  - host: foo.bar
    http:
      paths:
      - backend:
          serviceName: kubeapp-production
          servicePort: 80
EOF
```


### Preparation Step Resource16

<!-- @createResource16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/using-skip/production/ingress.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
EOF
```

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
mkdir ${DEMO_HOME}/actual/using-skip
mkdir ${DEMO_HOME}/actual/using-skip/production
mkdir ${DEMO_HOME}/actual/using-skip/canary
mkdir ${DEMO_HOME}/actual/using-composition
mkdir ${DEMO_HOME}/actual/using-composition/production
mkdir ${DEMO_HOME}/actual/using-composition/canary
kustomize build ${DEMO_HOME}/using-skip/production -o ${DEMO_HOME}/actual/using-skip/production
kustomize build ${DEMO_HOME}/using-skip/canary -o ${DEMO_HOME}/actual/using-skip/canary
kustomize build ${DEMO_HOME}/using-composition/composite/production -o ${DEMO_HOME}/actual/using-composition/production
kustomize build ${DEMO_HOME}/using-composition/composite/canary -o ${DEMO_HOME}/actual/using-composition/canary
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
mkdir ${DEMO_HOME}/expected/using-skip
mkdir ${DEMO_HOME}/expected/using-skip/production
mkdir ${DEMO_HOME}/expected/using-skip/canary
mkdir ${DEMO_HOME}/expected/using-composition
mkdir ${DEMO_HOME}/expected/using-composition/production
mkdir ${DEMO_HOME}/expected/using-composition/canary
```

### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-composition/canary/apiextensions.k8s.io_v1beta1_customresourcedefinition_mycrds.my.org.yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: mycrds.my.org
spec:
  additionalPrinterColumns: null
  group: my.org
  names:
    kind: MyCRD
    plural: mycrds
    shortNames:
    - mycrd
  scope: Cluster
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            replica:
              type: integer
            simpletext:
              type: string
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-composition/canary/default_my.org_v1alpha1_mycrd_my-crd.yaml
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  name: my-crd
spec:
  replica: 123
  simpletext: some simple text
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-composition/canary/~g_v1_namespace_kubeapp-ns.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: kubeapp-ns
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-composition/canary/kubeapp-ns_extensions_v1beta1_deployment_kubeapp-canary.yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: kubeapp
    env: canary
  name: kubeapp-canary
  namespace: kubeapp-ns
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubeapp
      env: canary
  template:
    metadata:
      labels:
        app: kubeapp
        env: canary
      name: kubeapp
    spec:
      containers:
      - image: hack4easy/kubesim_health-amd64:0.1.9
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        name: kubeapp
        ports:
        - containerPort: 8081
          name: kubeapp
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-composition/canary/kubeapp-ns_extensions_v1beta1_ingress_kubeapp.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  backend:
    serviceName: kubeapp-production
    servicePort: 80
  rules:
  - host: canary.foo.bar
    http:
      paths:
      - backend:
          serviceName: kubeapp-canary
          servicePort: 80
  - host: foo.bar
    http:
      paths:
      - backend:
          serviceName: kubeapp-production
          servicePort: 80
EOF
```


### Verification Step Expected5

<!-- @createExpected5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-composition/canary/kubeapp-ns_~g_v1_service_kubeapp-canary.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: kubeapp
    env: canary
  name: kubeapp-canary
  namespace: kubeapp-ns
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8081
  selector:
    app: kubeapp
    env: canary
  type: LoadBalancer
EOF
```


### Verification Step Expected6

<!-- @createExpected6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-composition/production/apiextensions.k8s.io_v1beta1_customresourcedefinition_mycrds.my.org.yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: mycrds.my.org
spec:
  additionalPrinterColumns: null
  group: my.org
  names:
    kind: MyCRD
    plural: mycrds
    shortNames:
    - mycrd
  scope: Cluster
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            replica:
              type: integer
            simpletext:
              type: string
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
EOF
```


### Verification Step Expected7

<!-- @createExpected7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-composition/production/default_my.org_v1alpha1_mycrd_my-crd.yaml
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  name: my-crd
spec:
  replica: 123
  simpletext: some simple text
EOF
```


### Verification Step Expected8

<!-- @createExpected8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-composition/production/~g_v1_namespace_kubeapp-ns.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: kubeapp-ns
EOF
```


### Verification Step Expected9

<!-- @createExpected9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-composition/production/kubeapp-ns_extensions_v1beta1_deployment_kubeapp-production.yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: kubeapp
    env: production
  name: kubeapp-production
  namespace: kubeapp-ns
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubeapp
      env: production
  template:
    metadata:
      labels:
        app: kubeapp
        env: production
      name: kubeapp
    spec:
      containers:
      - image: hack4easy/kubesim_health-amd64:0.1.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        name: kubeapp
        ports:
        - containerPort: 8081
          name: kubeapp
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
EOF
```


### Verification Step Expected10

<!-- @createExpected10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-composition/production/kubeapp-ns_extensions_v1beta1_ingress_kubeapp.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  backend:
    serviceName: kubeapp-production
    servicePort: 80
EOF
```


### Verification Step Expected11

<!-- @createExpected11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-composition/production/kubeapp-ns_~g_v1_service_kubeapp-production.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: kubeapp
    env: production
  name: kubeapp-production
  namespace: kubeapp-ns
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8081
  selector:
    app: kubeapp
    env: production
  type: LoadBalancer
EOF
```


### Verification Step Expected12

<!-- @createExpected12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-skip/canary/apiextensions.k8s.io_v1beta1_customresourcedefinition_mycrds.my.org.yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: mycrds.my.org
spec:
  additionalPrinterColumns: null
  group: my.org
  names:
    kind: MyCRD
    plural: mycrds
    shortNames:
    - mycrd
  scope: Cluster
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            replica:
              type: integer
            simpletext:
              type: string
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
EOF
```


### Verification Step Expected13

<!-- @createExpected13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-skip/canary/default_my.org_v1alpha1_mycrd_my-crd.yaml
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  name: my-crd
spec:
  replica: 123
  simpletext: some simple text
EOF
```


### Verification Step Expected14

<!-- @createExpected14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-skip/canary/~g_v1_namespace_kubeapp-ns.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: kubeapp-ns
EOF
```


### Verification Step Expected15

<!-- @createExpected15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-skip/canary/kubeapp-ns_extensions_v1beta1_deployment_kubeapp-canary.yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: kubeapp
    env: canary
  name: kubeapp-canary
  namespace: kubeapp-ns
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubeapp
      env: canary
  template:
    metadata:
      labels:
        app: kubeapp
        env: canary
      name: kubeapp
    spec:
      containers:
      - image: hack4easy/kubesim_health-amd64:0.1.9
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        name: kubeapp
        ports:
        - containerPort: 8081
          name: kubeapp
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
EOF
```


### Verification Step Expected16

<!-- @createExpected16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-skip/canary/kubeapp-ns_extensions_v1beta1_ingress_kubeapp.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  backend:
    serviceName: kubeapp-production
    servicePort: 80
  rules:
  - host: canary.foo.bar
    http:
      paths:
      - backend:
          serviceName: kubeapp-canary
          servicePort: 80
  - host: foo.bar
    http:
      paths:
      - backend:
          serviceName: kubeapp-production
          servicePort: 80
EOF
```


### Verification Step Expected17

<!-- @createExpected17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-skip/canary/kubeapp-ns_~g_v1_service_kubeapp-canary.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: kubeapp
    env: canary
  name: kubeapp-canary
  namespace: kubeapp-ns
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8081
  selector:
    app: kubeapp
    env: canary
  type: LoadBalancer
EOF
```


### Verification Step Expected18

<!-- @createExpected18 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-skip/production/apiextensions.k8s.io_v1beta1_customresourcedefinition_mycrds.my.org.yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: mycrds.my.org
spec:
  additionalPrinterColumns: null
  group: my.org
  names:
    kind: MyCRD
    plural: mycrds
    shortNames:
    - mycrd
  scope: Cluster
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            replica:
              type: integer
            simpletext:
              type: string
          type: object
  version: v1alpha1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
EOF
```


### Verification Step Expected19

<!-- @createExpected19 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-skip/production/default_my.org_v1alpha1_mycrd_my-crd.yaml
apiVersion: my.org/v1alpha1
kind: MyCRD
metadata:
  name: my-crd
spec:
  replica: 123
  simpletext: some simple text
EOF
```


### Verification Step Expected20

<!-- @createExpected20 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-skip/production/~g_v1_namespace_kubeapp-ns.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: kubeapp-ns
EOF
```


### Verification Step Expected21

<!-- @createExpected21 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-skip/production/kubeapp-ns_extensions_v1beta1_deployment_kubeapp-production.yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: kubeapp
    env: production
  name: kubeapp-production
  namespace: kubeapp-ns
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubeapp
      env: production
  template:
    metadata:
      labels:
        app: kubeapp
        env: production
      name: kubeapp
    spec:
      containers:
      - image: hack4easy/kubesim_health-amd64:0.1.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /liveness
            port: 8081
        name: kubeapp
        ports:
        - containerPort: 8081
          name: kubeapp
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
EOF
```


### Verification Step Expected22

<!-- @createExpected22 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-skip/production/kubeapp-ns_extensions_v1beta1_ingress_kubeapp.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    app: kubeapp
  name: kubeapp
  namespace: kubeapp-ns
spec:
  backend:
    serviceName: kubeapp-production
    servicePort: 80
EOF
```


### Verification Step Expected23

<!-- @createExpected23 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/using-skip/production/kubeapp-ns_~g_v1_service_kubeapp-production.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: kubeapp
    env: production
  name: kubeapp-production
  namespace: kubeapp-ns
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 8081
  selector:
    app: kubeapp
    env: production
  type: LoadBalancer
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

