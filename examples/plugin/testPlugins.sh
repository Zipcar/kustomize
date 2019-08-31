#!/bin/bash

set -e

echo Testing built-in plugins...

for i in `cat examples/plugin/builtInPluginList | grep -v '#'`
do
 echo $i
 mkdir -p examples/plugin/$i/expected
 kustomize build examples/plugin/$i -o examples/plugin/${i}/expected/result.yaml
done

echo Testing external plugins...

for i in `cat examples/plugin/externalPluginList | grep -v '#'`
do
 echo $i
 mkdir -p examples/plugin/$i/expected
 kustomize build examples/plugin/$i -o examples/plugin/${i}/expected/result.yaml --enable_alpha_plugins
done

echo All done.
