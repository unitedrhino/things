#! /bin/bash
set -e

RESPOSITORY="https://github.com/i4de/things.git"
#
init_env() {
    if [ ! -d "./_build/" ]; then
        mkdir -p ./_build/
    else
        rm -rf ./_build/
        mkdir -p ./_build/
    fi

}
#
fetch_resource() {
    cd _build
    git clone ${RESPOSITORY} things
    cd ..
}
#
build() {
    cd _build/things
    cd ./src/dmsvr && go get -v && go build -o dmsvr
    cd ..
    cd ./src/usersvr && go get -v && go build -o usersvr
    cd ..
    cd ./src/dcsvr && go get -v && go build -o dcsvr
    cd ..
    cd ./src/apisvr && go get -v && go build -o apisvr
    cd ..
    cd ./src/ddsvr && go get -v && go build -o ddsvr
    cd ..
}
#--------------------------------------------------
# Main
#--------------------------------------------------
init_env
fetch_resource
build
