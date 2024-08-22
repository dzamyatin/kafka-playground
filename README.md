# HELP

1) Use ./Makefile
2) Use ./registry-api.http

# TODO
https://github.com/jenkinsci/jenkins


# REGISTRY INTERACTION
https://docs.confluent.io/platform/current/schema-registry/schema_registry_onprem_tutorial.html


To view all the subjects registered in Schema Registry
```
curl --silent -X GET http://localhost:8081/subjects/ | jq .
```

To view the latest schema for this subject
```
curl --silent -X GET http://localhost:8081/subjects/transactions-value/versions/latest | jq .
```

Based on the schema id, you can also retrieve the associated schema
```
curl --silent -X GET http://localhost:8081/schemas/ids/1 | jq .
```

Create a new schema
```
curl -X POST -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  --data '{"schema": "{\"type\":\"record\",\"name\":\"Payment\",\"namespace\":\"io.confluent.examples.clients.basicavro\",\"fields\":[{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"amount\",\"type\":\"double\"}]}"}' \
  http://localhost:8081/subjects/test-value/versions
```


