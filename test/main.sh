#!/bin/bash

set -e

./care-cli --access-key 123 --secret-key 456
diff ~/.carerc.yml ./test/init.yml