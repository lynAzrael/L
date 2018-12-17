#!/bin/sh

TRACESERVER_VERSION="1.0"

# MakeGrafana
function MakeGrafana() {
    go get gitlab.33.cn/ops/grafana
    cd $GOPATH/src/gitlab.33.cn/ops/grafana
    go run build.go setup
    go run build.go build
}

function MakeTraceServer() {
    go get gitlab.33.cn/ops/trace
    cd $GOPATH/src/gitlab.33.cn/ops/trace
    go build
    tar zcvf trace_${TRACESERVER_VERSION}.tar.gz trace web
}

function DeployTraceServer() {
    
}

function DeployGrafana() {
    
}