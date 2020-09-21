#!/bin/bash

set -e

./care-cli --token 123
diff ~/.carerc.yml ./test/init.yml