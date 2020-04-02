#!/bin/bash
for i in 1 2 3 4 5 6 7 8 9 10
do
   echo "Welcome $i times"
   ./socket_timeout 10000
done
