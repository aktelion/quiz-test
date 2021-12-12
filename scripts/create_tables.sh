aws dynamodb create-table \
    --table-name Places \
    --attribute-definitions \
        AttributeName=Label,AttributeType=S \
    --key-schema \
        AttributeName=Label,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=5

aws dynamodb create-table \
    --table-name Games \
    --attribute-definitions \
        AttributeName=Id,AttributeType=N \
    --key-schema \
        AttributeName=Id,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=5