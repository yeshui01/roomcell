#!/bin/bash
TARGETBINPATH=/data/myth/deploy/sync_env/bin/
TARGETPATHCSV=/data/myth/deploy/sync_env/csv/
SERVERADDR=112.74.107.204
#echo "begin sync src to ${TARGETPATHSERVER1}"
#rsync -vr ./bin/ --exclude .git/ --exclude .gitignore --exclude .gitlab-ci.yml  myth@${SERVERADDR}:${TARGETBINPATH}
echo "begin sync csv to ${TARGETPATHCSV}"
rsync -vr ./csv/ --exclude .git/ --exclude .gitignore --exclude .gitlab-ci.yml  myth@${SERVERADDR}:${TARGETPATHCSV}
