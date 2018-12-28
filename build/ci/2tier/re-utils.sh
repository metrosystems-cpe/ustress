
#!/bin/bash

set -eu -o pipefail
set +x

SHEESY_CLONE_URL=https://github.com/share-secrets-safely/getting-started.git
IDAM_HOST="https://idam.metrosystems.net"
DEPLOYMENT_SERVICE_URL=https://deploymentservice.metrosystems.net
utils_root="$PWD/.tools"
mkdir -p "$utils_root"

SY_EXE="$utils_root/sy/sy"


function alpine_init() {
  apk add -U curl git ca-certificates
}

function git_hash() {
  git rev-parse HEAD | head -c 8
}


function docker_login {
  echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin "$DOCKER_REGISTRY_HOST"
}

function local_tag {
  local service=${1:?}
  echo "${service}-run-image:$(git_hash)"
}


function tag_image () (
  set -eu
  local stage="${1:?}"
  local vertical="${2:?}"
  local service="${3:?}"
  local localTag
  local remote_url
  # duno what this does but seems important :)
  {
    get_sy
  } 1>&2

  docker_source_image_tag="${service}:latest"
  docker_target_image_tag="$(image_url "$stage" "$vertical" "$service")"
  
  {
    docker tag "$docker_source_image_tag" "$docker_target_image_tag"
  } 1>&2
  echo "$docker_target_image_tag"
)


function build_and_tag_image () (
  set -eu
  local stage="${1:?}"
  local vertical="${2:?}"
  local service="${3:?}"
  local run_dockerfile_path="${4:?}"
  local localTag
  local remote_url

  {
    get_sy
  } 1>&2

  localTag="$(local_tag "${service}")"
  remote_url="$(image_url "$stage" "$vertical" "$service")"
  {
    local tmpfile="$utils_root/tmp"
    # shellcheck disable=2064
    trap "rm -f $tmpfile" EXIT
    {
      grep -vi "^ENV.*DRP_" "$run_dockerfile_path"
      echo '{}' \
         | "$SY_EXE" process - "COMMIT_HASH=$(git_hash)" \
         | "$SY_EXE" substitute <(grep -i "^ENV.*DRP_" "$run_dockerfile_path")
    } > "$tmpfile"
    docker build --network host -f "$tmpfile" -t "$localTag" .
    docker tag "$localTag" "$remote_url"
  } 1>&2
  echo "$remote_url"
)

function push_image_to_registry {
  local remote_url="${1:?}"
  docker_login
  docker_push "$remote_url"
}

function docker_push() {
  docker push "${1:?}"
}

function image_url {
  local stage="${1:?}"
  local vertical="${2:?}"
  local service="${3:?}"
  echo "$DOCKER_REGISTRY_HOST/$stage/$vertical/$service:$(git_hash)"
}

function get_sy(){
  if [ ! -f "$SY_EXE" ]; then
    git clone $SHEESY_CLONE_URL "${SY_EXE%/*}"
    "$SY_EXE" --version >/dev/null
  fi
}

function idam_access_token() {
  local auth
  auth=$(basic_auth_string "$IDAM_CLIENT_ID" "$IDAM_SECRET")


  curl --fail -X POST \
    "$IDAM_HOST/authorize/api/oauth2/access_token" \
    -H "authorization: $auth" \
    -H 'content-type: application/x-www-form-urlencoded' \
    -d "grant_type=client_credentials&client_id=$IDAM_CLIENT_ID&realm_id=PENG_2TR_RLM" 2>/dev/null \
    | "$SY_EXE" extract .access_token
}

# takes user name and pwd and returns the basic auth string
function basic_auth_string() {
  local user=${1?User must be set as first arg.}
  local password=${2?Password must be set as second arg.}
  local auth
  auth=$(echo -n "$user:$password" | base64)
  echo "Basic $auth"
}

function deploy_to_ds () (
  set -e
  local remote_image_url="${1:?}"
  local stage="${2:?}"
  local ds_payload_file_path=${3:?}
  local ds_resources_file_path=${4:?}
  local ds_model
  local payload
  local deployment_id_path
  local deployment_result
  local auth_header

  local status
  local check_intervall=10
  get_sy

  ds_model="{ IMAGE_URL: '${remote_image_url}', COMMIT_HASH: '$(git_hash)', STAGE: '$stage'}"
  payload=$( echo '{}' \
             | "$SY_EXE" merge - "$ds_payload_file_path" --select="$stage" "$ds_resources_file_path" \
             | "$SY_EXE" substitute -d <(echo "$ds_model") )
  auth_header="authorization: Bearer $(idam_access_token)"
  deployment_result=$(curl --fail -i -s \
                  -H "$auth_header" \
                  -X POST \
                  -d "$payload" \
                  $DEPLOYMENT_SERVICE_URL/deployments)

  deployment_id_path="$(echo -n "$deployment_result" | grep "Location" | awk -F : '{print $2}' | tr -d '[:space:]')"
  [[ "${deployment_id_path}" = "" ]] && {
    echo 1>&2 "Deployment seemed to have failed with output:"
    echo 1>&2 "$deployment_result"
    return 1
  }
  echo "Deployment is created and sent to the Deployment Service."

  set +e
  while true; do
    check_deployment "$deployment_id_path" "$auth_header"
    status=$?
    if [ $status -ne 2 ]; then
      return $status
    fi
    pause_a_while $check_intervall
  done
  set -e
)

function pause_a_while(){
  local sleep_time="${1:?}"
  sleep "$sleep_time"
}

function check_deployment {
  local deployment_id_path="${1:?}"
  local auth_header="${2:?}"
  local jsonresponse
  local status
  local message

  jsonresponse=$(curl -s -H "$auth_header" "${DEPLOYMENT_SERVICE_URL}${deployment_id_path}")
  status=$(echo "$jsonresponse" | "$SY_EXE" extract status )

  if [[ "$status" == "succeeded" ]]; then
      echo "Status: $status"
      return 0
  elif [[ "$status" == "failed" ]]; then
      message=$(echo "$jsonresponse" | "$SY_EXE" extract statusMessage )
      echo "Status: $status -- Message: $message"
      return 1
  fi
  return 2
}
