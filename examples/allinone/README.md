# Issue: XXX

First, define a place to reproduce the issue

<!-- @makeWorkplace @dev @site1 @site2 @test -->
```sh
DEMO_HOME=$(mktemp -d)
CONTENT="https://raw.githubusercontent.com\
/keleustes/kustomize\
/allinone/examples/allinone"
```

## Import the files into the workplace

<!-- @installResources @dev @site1 @site2 @test -->
```sh
mkdir -p $DEMO_HOME/actual
mkdir -p $DEMO_HOME/expected
mkdir -p $DEMO_HOME/base
mkdir -p $DEMO_HOME/base/catalogues
mkdir -p $DEMO_HOME/base/crds
mkdir -p $DEMO_HOME/base/kustomizeconfig
mkdir -p $DEMO_HOME/base/mysql
mkdir -p $DEMO_HOME/base/wordpress

curl -s -o "$DEMO_HOME/base/#1.yaml" \
  "$CONTENT/base/{kustomization}.yaml"

curl -s -o "$DEMO_HOME/base/catalogues/#1.yaml" \
  "$CONTENT/base/catalogues/{common-addresses,endpoints,versions}.yaml"

curl -s -o "$DEMO_HOME/base/crds/#1.yaml" \
  "$CONTENT/base/crds/{Chart,CommonAddresses,EndpointCatalogue,SoftwareVersions}.yaml"

curl -s -o "$DEMO_HOME/base/kustomizeconfig/#1.yaml" \
  "$CONTENT/base/kustomizeconfig/{Chart,CommonAddresses,Deployment,EndpointCatalogue,Service,SoftwareVersions}.yaml"

curl -s -o "$DEMO_HOME/base/mysql/#1.yaml" \
  "$CONTENT/base/mysql/{deployment,secret,service}.yaml"

curl -s -o "$DEMO_HOME/base/wordpress/#1.yaml" \
  "$CONTENT/base/wordpress/{deployment,service}.yaml"
```

## Development Environment build and test

### Import expected files for dev environment

<!-- @importDev @dev @test -->
```sh
mkdir -p $DEMO_HOME/dev
curl -s -o "$DEMO_HOME/dev/#1.yaml" \
  "$CONTENT/dev/{common-addresses,devtools,endpoints,kustomization,passphrases,versions}.yaml"

mkdir -p $DEMO_HOME/actual/dev
mkdir -p $DEMO_HOME/expected/dev
curl -s -o "$DEMO_HOME/expected/dev/#1.yaml" \
  "$CONTENT/expected/dev/{apps_v1beta2_deployment_dev-mysql,apps_v1beta2_deployment_dev-wordpress,~g_v1_secret_dev-mysql-pass,~g_v1_service_dev-mysql,~g_v1_service_dev-wordpress,my.group.org_v1alpha1_chart_dev-mysql,my.group.org_v1alpha1_chart_dev-wordpress,my.group.org_v1alpha1_commonaddresses_dev-common-addresses,my.group.org_v1alpha1_endpointcatalogue_dev-endpoints,my.group.org_v1alpha1_softwareversions_dev-software-versions}.yaml"
```

<!-- @buildDev @dev @test -->
```sh
kustomize build $DEMO_HOME/dev -o $DEMO_HOME/actual/dev
```

### Verify that the actual output is matching the expected output dev

<!-- @verifyDev @dev @test -->
```sh
#diff -r $DEMO_HOME/actual/dev $DEMO_HOME/expected/dev > $DEMO_HOME/diffs-dev.txt
```

## Common Production Sites setup.

### Import files

<!-- @importCommon @site1 @site2 @test -->
```sh
mkdir -p $DEMO_HOME/production/common
curl -s -o "$DEMO_HOME/production/common/#1.yaml" \
  "$CONTENT/production/common/{endpoints,kustomization,versions}.yaml"
```

## Production Site1 build and test

### Import files for site1

<!-- @importSite1 @site1 @test -->
```sh
mkdir -p $DEMO_HOME/production/site1
curl -s -o "$DEMO_HOME/production/site1/#1.yaml" \
  "$CONTENT/production/site1/{common-addresses,kustomization,passphrases}.yaml"

mkdir -p $DEMO_HOME/actual/site1
mkdir -p $DEMO_HOME/expected/site1
curl -s -o "$DEMO_HOME/expected/site1/#1.yaml" \
  "$CONTENT/expected/site1{apps_v1beta2_deployment_mysql,apps_v1beta2_deployment_wordpress,~g_v1_secret_mysql-pass,~g_v1_service_mysql,~g_v1_service_wordpress,my.group.org_v1alpha1_commonaddresses_common-addresses,my.group.org_v1alpha1_endpointcatalogue_endpoints,my.group.org_v1alpha1_softwareversions_software-versions}.yaml"
```

### Build using kustomize site1

<!-- @buildSite1 @site1 @test -->
```sh
kustomize build $DEMO_HOME/production/site1 -o $DEMO_HOME/actual/site1
```

### Verify that the actual output is matching the expected output site1

<!-- @verifySite1 @site1 @test -->
```sh
#diff -r $DEMO_HOME/actual/site1 $DEMO_HOME/expected/site1 > $DEMO_HOME/diffsSite1.txt
```

## Production Site2 build and test

### Import files for site2

<!-- @importSite2 @site2 @test -->
```sh
mkdir -p $DEMO_HOME/production/site2
curl -s -o "$DEMO_HOME/production/site2/#1.yaml" \
  "$CONTENT/production/site2/{common-addresses,kustomization,passphrases}.yaml"

mkdir -p $DEMO_HOME/actual/site2
mkdir -p $DEMO_HOME/expected/site2
curl -s -o "$DEMO_HOME/expected/site2/#1.yaml" \
  "$CONTENT/expected/site2/{apps_v1beta2_deployment_mysql,apps_v1beta2_deployment_wordpress,~g_v1_secret_mysql-pass,~g_v1_service_mysql,~g_v1_service_wordpress,my.group.org_v1alpha1_commonaddresses_common-addresses,my.group.org_v1alpha1_endpointcatalogue_endpoints,my.group.org_v1alpha1_softwareversions_software-versions}.yaml"
```

### Build using kustomize site2

<!-- @buildSite2 @site2 @test -->
```sh
kustomize build $DEMO_HOME/production/site2 -o $DEMO_HOME/actual/site2
```

### Verify that the actual output is matching the expected output site2

<!-- @verifySite2 @site2 @test -->
```sh
#diff -r $DEMO_HOME/actual/site2 $DEMO_HOME/expected/site2 > $DEMO_HOME/diffsSite2.txt
```
