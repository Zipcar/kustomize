#!/bin/bash
rm -fr /tmp/tmp.*
rm -fr /tmp/kust*
rm -fr /tmp/mdrip*
for i in `cat readmelist`
do
   echo "=========================" $i "====================="
   mdrip --mode test $i
done
