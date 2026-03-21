#!/bin/bash
awslocal dynamodb create-table \
    --table-name user-auth-tokens \
    --attribute-definitions AttributeName=token_id,AttributeType=S \
    --key-schema AttributeName=token_id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
