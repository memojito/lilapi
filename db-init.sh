#!/bin/bash

# Execute the scripts in order
echo "Creating keyspace..."
cqlsh -f db/scripts/create-keyspace.cql

echo "Creating table..."
cqlsh -f db/scripts/schema.cql

echo "CQL commands executed successfully!"
exit