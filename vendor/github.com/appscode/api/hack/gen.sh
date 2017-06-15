#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

RETVAL=0
ROOT=$PWD
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

ALIAS="Mappscode/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/appscodeapis/appscode/api,"
ALIAS+="Mappscode/api/dtypes/types.proto=github.com/appscode/api/dtypes,"
ALIAS+="Mappscode/api/certificate/v1beta1/certificate.proto=github.com/appscode/api/certificate/v1beta1,"
ALIAS+="Mappscode/api/kubernetes/v1beta1/cluster.proto=github.com/appscode/api/kubernetes/v1beta1,"
ALIAS+="Mappscode/api/kubernetes/v1beta2/client.proto=github.com/appscode/api/kubernetes/v1beta2,"
ALIAS+="Mappscode/api/ssh/v1beta1/ssh.proto=github.com/appscode/api/ssh/v1beta1,"
ALIAS+="Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any,"
ALIAS+="Mappscode/api/attic/artifacts.proto=github.com/appscode/api/attic,"
ALIAS+="Mappscode/api/db/v1beta1/database.proto=github.com/appscode/api/db/v1beta1,"
ALIAS+="Mappscode/api/db/v1beta1/snapshot.proto=github.com/appscode/api/db/v1beta1,"
ALIAS+="Mappscode/api/namespace/v1beta1/team.proto=github.com/appscode/api/namespace/v1beta1,"
ALIAS+="Mappscode/api/alert/v1beta1/incident.proto=github.com/appscode/api/alert/v1beta1,"
ALIAS+="Mappscode/api/operation/operation.proto=github.com/appscode/api/operation,"
ALIAS+="Mappscode/api/version/version.proto=github.com/appscode/api/version"

clean() {
	(find . | grep pb.go | xargs rm) || true
	(find . | grep pb.gw.go | xargs rm) || true
	(find . | grep cors.go | xargs rm) || true
	(find . | grep gw.cors.go | xargs rm) || true
	(find . | grep gw.js | xargs rm) || true
	# Do NOT delete schema.json files as they contain handwritten validation rules.
	# contact tamal@ / sadlil@ if in doubt.
	# (find . | grep schema.json | xargs rm) || true
	(find . | grep schema.go | xargs rm) || true
	(find . | grep php | xargs rm) || true
	(find . | grep _pb2.py | xargs rm) || true
}

gen_proto() {
  if [ $(ls -1 *.proto 2>/dev/null | wc -l) = 0 ]; then
    return
  fi
  rm -rf *.pb.go
  protoc -I /usr/local/include -I . \
         -I ${GOPATH}/src/github.com \
         -I ${GOPATH}/src/github.com/googleapis/googleapis/ \
         -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/appscodeapis \
         --go_out=plugins=grpc,${ALIAS}:. *.proto
}

gen_gateway_proto() {
  if [ $(ls -1 *.proto 2>/dev/null | wc -l) = 0 ]; then
    return
  fi
  rm -rf *.pb.gw.go
  protoc -I /usr/local/include -I . \
         -I ${GOPATH}/src/github.com \
         -I ${GOPATH}/src/github.com/googleapis/googleapis/ \
         -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/appscodeapis \
         --grpc-gateway_out=logtostderr=true,${ALIAS}:. *.proto
}

gen_cors_pattern() {
  if [ $(ls -1 *.proto 2>/dev/null | wc -l) = 0 ]; then
    return
  fi
  rm -rf *.gw.cors.go
  protoc -I /usr/local/include -I . \
         -I ${GOPATH}/src/github.com \
         -I ${GOPATH}/src/github.com/googleapis/googleapis/ \
         -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/appscodeapis \
         --grpc-gateway-cors_out=logtostderr=true,${ALIAS}:. *.proto
}

gen_js_client() {
  if [ $(ls -1 *.proto 2>/dev/null | wc -l) = 0 ]; then
    return
  fi
  rm -rf *.gw.js
  protoc -I /usr/local/include -I . \
         -I ${GOPATH}/src/github.com \
         -I ${GOPATH}/src/github.com/googleapis/googleapis/ \
         -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/appscodeapis \
         --grpc-js-client_out=logtostderr=true,remove_prefix=/_appscode/api,${ALIAS}:. *.proto
}

gen_swagger_def() {
  if [ $(ls -1 *.proto 2>/dev/null | wc -l) = 0 ]; then
    return
  fi
   rm -rf *.swagger.json
   protoc -I /usr/local/include -I . \
          -I ${GOPATH}/src/github.com \
          -I ${GOPATH}/src/github.com/googleapis/googleapis/ \
          -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/appscodeapis \
          --swagger_out=logtostderr=true,${ALIAS}:. *.proto
}

gen_server_protos() {
	echo "Generating server protobuf files"
    for d in */ ; do
      pushd ${d}
      gen_proto
      if dirs=$( ls -1 -F | grep "^v.*/" | tr -d "/" ); then
        for dd in $dirs ; do
          pushd ${dd}
          gen_proto
          popd
        done
      fi
      popd
    done
}

gen_proxy_protos() {
    echo "Generating gateway protobuf files"
    for d in */ ; do
      pushd ${d}
      gen_gateway_proto
      if dirs=$( ls -1 -F | grep "^v.*/" | tr -d "/" ); then
        for dd in $dirs ; do
          pushd ${dd}
          gen_gateway_proto
          popd
        done
      fi
      popd
    done
}

gen_cors_patterns() {
    echo "Generating gateway cors files"
    for d in */ ; do
      pushd ${d}
      gen_cors_pattern
      if dirs=$( ls -1 -F | grep "^v.*/" | tr -d "/" ); then
        for dd in $dirs ; do
          pushd ${dd}
          gen_cors_pattern
          popd
        done
      fi
      popd
    done
}

gen_js_clients() {
    echo "Generating protobuf js client"
    for d in */ ; do
      pushd ${d}
      gen_js_client
      if dirs=$( ls -1 -F | grep "^v.*/" | tr -d "/" ); then
        for dd in $dirs ; do
          pushd ${dd}
          gen_js_client
          popd
        done
      fi
      popd
    done
}

gen_swagger_defs() {
    echo "Generating swagger api definition files"
    for d in */ ; do
      pushd ${d}
      gen_swagger_def
      if dirs=$( ls -1 -F | grep "^v.*/" | tr -d "/" ); then
        for dd in $dirs ; do
          pushd ${dd}
          gen_swagger_def
          popd
        done
      fi
      popd
    done
    # fix host, schemes
    python $DIR/schema.py fix_swagger_schema
}

gen_py() {
  if [ $(ls -1 *.proto 2>/dev/null | wc -l) = 0 ]; then
    return
  fi
  rm -rf *.py
  python -m grpc.tools.protoc \
         -I /usr/local/include -I . \
         -I ${GOPATH}/src/github.com \
         -I ${GOPATH}/src/github.com/googleapis/googleapis/ \
         -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/appscodeapis \
         --python_out=':.' --grpc_python_out=':.' *.proto
}

gen_python_protos() {
  gen_py
  for d in */ ; do
    pushd ${d}
    gen_py
    if dirs=$( ls -1 -F | grep "^v.*/" | tr -d "/" ); then
      for dd in $dirs ; do
        pushd ${dd}
        gen_py
        popd
      done
    fi
    popd
  done
}

gen_php() {
  if [ $(ls -1 *.proto 2>/dev/null | wc -l) = 0 ]; then
    return
  fi
  rm -rf *.php
  protoc -I /usr/local/include -I . \
         -I ${GOPATH}/src/github.com \
         -I ${GOPATH}/src/github.com/googleapis/googleapis/ \
         -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/appscodeapis \
         --plugin=protoc-gen-grpc="$(which grpc_php_plugin)" \
         --php_out=':.' --grpc_out=':.' *.proto
}

gen_php_protos() {
  for d in */ ; do
    pushd ${d}
    gen_php
    if dirs=$( ls -1 -F | grep "^v.*/" | tr -d "/" ); then
      for dd in $dirs ; do
        pushd ${dd}
        gen_php
        popd
      done
    fi
    popd
  done
}

compile() {
  echo "compiling files"
  go install ./...
}

gen_protos() {
  clean
  python $DIR/schema.py gen_assets
  gen_server_protos
  gen_proxy_protos
  gen_cors_patterns
  # gen_js_clients
  gen_swagger_defs
  python $DIR/schema.py
  # gen_python_protos
  # gen_php_protos
  compile
}

if [ $# -eq 0 ]; then
	gen_protos
	exit $RETVAL
fi

case "$1" in
  compile)
      compile
      ;;
	server)
		gen_server_protos
		;;
	proxy)
		gen_proxy_protos
		gen_cors_patterns
		;;
	js)
		gen_js_clients
		;;
	swagger)
		gen_swagger_defs
		;;
	all)
		gen_protos
		;;
	clean)
	  clean
	  ;;
	py)
	  gen_python_protos
	  ;;
	php)
	  gen_php_protos
	  ;;
	*)  echo $"Usage: $0 {compile|server|proxy|js|swagger|json-schema|all|clean}"
		RETVAL=1
		;;
esac
exit $RETVAL
