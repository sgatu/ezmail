#!/bin/bash

ok=false
retries=0

function ping() {
    echo "Running ping..."
    retries=$((retries + 1))
    mysql --host=database --user=$DB_USER --password=$DB_PASS -e "SELECT 1;"
    if [ $? -eq 0 ]; then
        ok=true
    fi
}

while ! [ $ok ]; do
    echo "Checking mysql status..."
    ping
    if [ 10 -eq $retries ]; then
        exit 1
    fi
    sleep 2
done

if ! [[ $(mysql --host=database --user=$DB_USER --password=$DB_PASS -e "SHOW TABLES FROM \`${DB_NAME}\` like 'domain';" | wc -l) -gt 1 ]]; then
    mysql --host=database --user=$DB_USER --password=$DB_PASS $DB_NAME </opt/schema.sql
fi
