# HELP

1) Use ./Makefile
2) Use ./registry-api.http

# TODO
https://github.com/jenkinsci/jenkins


# REGISTRY INTERACTION
https://docs.confluent.io/platform/current/schema-registry/schema_registry_onprem_tutorial.html

## API doc:
https://docs.confluent.io/platform/current/schema-registry/develop/api.html#put--mode-(string-%20subject)

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

Change compatibility setting on subject:
```
curl -X PUT -H "Content-Type: application/vnd.schemaregistry.v1+json" \
       --data '{"compatibility": "BACKWARD_TRANSITIVE"}' \
       http://localhost:8081/config/transactions-value
```

Schemas list:
```
curl --silent -X GET http://localhost:8081/schemas/ | jq .
```


## Import
1) Activate import mode:
```
curl -X PUT -H "Content-Type: application/json" "http://localhost:8081/mode/my-cool-subject?force=1" --data '{"mode": "IMPORT"}'
```
2) Actually import:
```
curl -X POST -H "Content-Type: application/json" \
--data '{"schemaType": "AVRO", "version":1, "id":24, "schema":"{\"type\":\"record\",\"name\":\"value_a1\",\"namespace\":\"com.mycorp.mynamespace\",\"fields\":[{\"name\":\"field1\",\"type\":\"string\"}]}" }' \
http://localhost:8081/subjects/my-cool-subject/versions
```
3) Return subject to normal mode
```
curl -X PUT -H "Content-Type: application/json" "http://localhost:8081/mode/my-cool-subject?force=1" --data '{"mode": "READWRITE"}'
```
