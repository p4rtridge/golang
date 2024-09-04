#!/bin/bash

echo "##### Starting to execute script #####"

echo "##### Sleeping for 10 seconds #####"
sleep 10

while ! psql -d "${DATABASE}" -U "${USER_NAME}"; do
  echo "Waiting for main instance to be ready..."
  sleep 5
done

for sql_file in ./migrations/*.sql; do
  psql -d "${DATABASE}" -U "${USER_NAME}" -f "${sql_file}"
  echo "Script ""${sql_file}"" executed!!!"
done

echo "##### Execution script is finished! #####"
echo "##### Stopping temporary instance! #####"
