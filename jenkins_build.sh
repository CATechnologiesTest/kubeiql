#!/bin/bash
#
# In Jenkins, make a free-form job that exports the following env vars
# prior to invoking this script:
#
# THIS_CONTAINER=<name of the container for jenkins itself>
# WORKSPACE=<jenkins workspace for the build>
# PROJDIR=<location of go sources for the project, e.g., $WORKSPACE/src/auth
# PG_NAME=<name for the postgres DB container used by this job>
#
docker run -d --name $PG_NAME postgres:9.5.5-alpine
docker run --rm --link $PG_NAME:db --volumes-from=$THIS_CONTAINER yipee-tools-spoke-cos.ca.com:5000/yipee-go-builder bash -c "export GOPATH=$WORKSPACE; export POSTGRES_HOST=db; export POSTGRES_DB=postgres export POSTGRES_USER=postgres; export YIPEE_TEAM_OWNER=jenkinsbuild; cd $PROJDIR; bash -x ./gobuild.sh"
rc=$?
docker rm -fv $PG_NAME
exit $rc
