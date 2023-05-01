#!/bin/bash

# Execute the scripts in order
echo "Creating keyspace..."
cqlsh -f docker-entrypoint-initdb.d/db/scripts/create-keyspace.cql

echo "Creating table..."
cqlsh -f docker-entrypoint-initdb.d/db/scripts/schema.cql

echo "CQL commands executed successfully!"
exit