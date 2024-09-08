aws dynamodb create-table \
    --table-name urls \
    --attribute-definitions AttributeName=short,AttributeType=S \
    --key-schema AttributeName=short,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=40000,WriteCapacityUnits=40000 \
    --endpoint-url http://localhost:8000 \
    --region us-west-2 \
    --cli-input-json file://urls-table.json
