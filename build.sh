#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${DIR}
GOPATHt="${DIR}/build"
build_path="${GOPATHt}/src/github.com/doout/cps842"
rm -rf ${build_path}
mkdir -p ${build_path}
cp -r $(ls | grep -v build) ${build_path}
export GOPATH="${GOPATHt}"
cd ${build_path}
make build
cd ${DIR}
rm -rf output
mkdir output
cp -r "${build_path}/output" .
rm -rf ${GOPATHt}
