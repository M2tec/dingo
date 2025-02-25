#!/usr/bin/env python3
import traceback
import requests
import time

mutation="""mutation {
  updateGQLSchema(
    input: { set: { schema: "$my_schema"}})
  {
    gqlSchema {
      schema
      generatedSchema
    }
  }
}
"""
f = open("cardano-schema.gql", "r")
schema = f.read() #.replace('\n', '')
#print(schema) 

mutation = mutation.replace("$my_schema", schema)
mutation = mutation.replace('\n', '').replace('  ', ' ').replace('\t', ' ').replace('  ', ' ')
print(mutation)

for i in range(0, 15):
  try:
    r = requests.post('http://localhost:8080/admin', json={"query": mutation})
  except requests.exceptions.ConnectionError:
    time.sleep(1)
    print("Waiting for connection")

message = r.json()  

for i in range(0, 15):
  if "errors" in message:
    r = requests.post('http://localhost:8080/admin', json={"query": mutation})
    print("Waiting for db to be ready ... " + str(i))# , end="  ")
    time.sleep(3)


print(r)
print(r.json())


