#!/bin/bash

export BASEDIR=.

set -e

echo Testing built-in plugins...

for i in `cat ${BASEDIR}/builtInPluginList | grep -v '#'`
do
 echo $i
 mkdir -p ${BASEDIR}/$i/expected
 kustomize build ${BASEDIR}/$i -o ${BASEDIR}/${i}/expected/result.yaml
done

echo Testing external plugins...

for i in `cat ${BASEDIR}/externalPluginList | grep -v '#'`
do
 echo $i
 mkdir -p ${BASEDIR}/$i/expected
 kustomize build ${BASEDIR}/$i -o ${BASEDIR}/${i}/expected/result.yaml --enable_alpha_plugins
done

echo All done.
