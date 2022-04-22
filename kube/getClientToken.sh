#! /usr/bin/env bash




SUDORY_SERVER_URL=$1
CLUSTER_NAME="new-node-cluster"
CLUSTER_DESCRIPTION='TEST-CLUSTER'

user_uuid=$(curl --silent -X POST $SUDORY_SERVER_URL/server/cluster \
            -H "Content-Type: application/json"  \
            --data '{ "name": "'$CLUSTER_NAME'", "polling_option": { "additionalProp1": {} }, "summary": "'$CLUSTER_DESCRIPTION'" }' | jq -r '.uuid')


token=$(curl --silent -X POST $SUDORY_SERVER_URL/server/token/cluster \
            -H "Content-Type: application/json"  \
            --data '{ "name": "'$CLUSTER_NAME'", "uuid": "'$uuid'" }' | jq -r '.token')


export S_SERVER_URL=$SUDORY_SERVER_URL
export S_CLUSTER_ID=$user_uuid
export S_TOKEN=$token

envsubst '${S_SERVER_URL} ${S_CLUSTER_ID} ${S_TOKEN}' < kube/k8s-deploy-nexclipper-sudory-client.yaml > kube/k8s-deploy-nexclipper-sudory-client.yaml_gen 
mv kube/k8s-deploy-nexclipper-sudory-client.yaml_gen kube/k8s-deploy-nexclipper-sudory-client.yaml