#!/bin/sh

set -e

WGET=$(which wget)
BOOTSTRAP_VERSION="5.0.0-beta1"
BOOTSTRAP_OUT="bootstrap.zip"
LEAFLET_VERSION="1.6.0"
LEAFLET_OUT="leaflet.zip"
STATIC_DIR="static/3rdparty"

pre_init() {
    rm -rf "${STATIC_DIR}/*"
}

init_bootstrap() {
    "${WGET}" -q "https://github.com/twbs/bootstrap/releases/download/v${BOOTSTRAP_VERSION}/bootstrap-${BOOTSTRAP_VERSION}-dist.zip" -O "${BOOTSTRAP_OUT}"
    unzip -o -d "${STATIC_DIR}" "${BOOTSTRAP_OUT}" \
        bootstrap-${BOOTSTRAP_VERSION}-dist/css/bootstrap.min.css \
        bootstrap-${BOOTSTRAP_VERSION}-dist/css/bootstrap.min.css.map
}

init_leaflet() {
    "${WGET}" -q "http://cdn.leafletjs.com/leaflet/v${LEAFLET_VERSION}/leaflet.zip" -O "${LEAFLET_OUT}"
    mkdir -p "${STATIC_DIR}/leaflet-v${LEAFLET_VERSION}"
    unzip -o -d "${STATIC_DIR}/leaflet-v${LEAFLET_VERSION}" "${LEAFLET_OUT}"
}

cleanup() {
    rm -f \
        "${BOOTSTRAP_OUT}" \
        "${LEAFLET_OUT}"
}

pre_init
init_bootstrap
init_leaflet
cleanup
