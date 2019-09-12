#!/bin/bash

set -e

echo Generating linkable plugins...

for i in `cat examples/plugin/pluginList | grep -v '#'`
do
 echo $i
 mkdir -p examples/plugin/$i/expected
 kustomize build examples/plugin/$i --enable_alpha_plugins -o examples/plugin/${i}/expected/result.yaml
done

echo All done.

