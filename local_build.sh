#!/bin/bash
REGISTRY=yipee-development

export POSTGRES_HOST=localhost
export POSTGRES_SSL=disable
export POSTGRES_DB=postgres
export POSTGRES_USER=postgres
export YIPEE_TEAM_OWNER=$TOKEN_USER

docker-compose -f PostgreSQL.yml up -d
bash gobuild.sh
result=$?
docker-compose -f PostgreSQL.yml down
docker volume rm auth_data >& /dev/null
test $result -eq 0 || exit 1
IMAGE=auth
docker build -t $REGISTRY/$IMAGE .
docker tag $REGISTRY/$IMAGE yipee-tools-spoke-cos.ca.com:5000/$IMAGE:development

