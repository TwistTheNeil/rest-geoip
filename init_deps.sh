#!/bin/sh

WGET=$(which wget)
BOOTSTRAP_VERSION="4.5.0"
BOOTSTRAP_OUT="bootstrap.zip"
LEAFLET_VERSION="1.6.0"
LEAFLET_OUT="leaflet.zip"

function init_bootstrap {
    "${WGET}" "https://github.com/twbs/bootstrap/releases/download/v${BOOTSTRAP_VERSION}/bootstrap-${BOOTSTRAP_VERSION}-dist.zip" -O "${BOOTSTRAP_OUT}"
    unzip -d static "${BOOTSTRAP_OUT}" \
        bootstrap-4.5.0-dist/css/bootstrap.min.css \
        bootstrap-4.5.0-dist/css/bootstrap.min.css.map
}

function init_leaflet {
    "${WGET}" "http://cdn.leafletjs.com/leaflet/v${LEAFLET_VERSION}/leaflet.zip" -O "${LEAFLET_OUT}"
    unzip -d "static/leaflet-v${LEAFLET_VERSION}" "${LEAFLET_OUT}"
}

function cleanup {
    rm -f \
        "${BOOTSTRAP_OUT}" \
        "${LEAFLET_OUT}"
}

init_bootstrap
init_leaflet
cleanup
