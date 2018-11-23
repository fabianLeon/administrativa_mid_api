#!/usr/bin/env bash

set -e
set -u
set -o pipefail

if [ -n "${PARAMETER_STORE:-}" ]; then
  export ADMINISTRATIVA_MID_API__PGUSER="$(aws ssm get-parameter --name /${PARAMETER_STORE}/administrativa_mid_api/db/username --output text --query Parameter.Value)"
  export ADMINISTRATIVA_MID_API__PGPASS="$(aws ssm get-parameter --with-decryption --name /${PARAMETER_STORE}/administrativa_mid_api/db/password --output text --query Parameter.Value)"
  
  export ADMINISTRATIVA_MID_API_AGORA_USER="$(aws ssm get-parameter --name /${PARAMETER_STORE}/agora_mid/db/username --output text --query Parameter.Value)"
  export ADMINISTRATIVA_MID_API_AGORA_PASS="$(aws ssm get-parameter --with-decryption --name /${PARAMETER_STORE}/agora_mid/db/password --output text --query Parameter.Value)"
  
fi

exec ./main "$@"
