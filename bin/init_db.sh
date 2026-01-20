#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DB_FILE="$SCRIPT_DIR/../database/main.db"
SQL_FILE="$SCRIPT_DIR/../sql/schema.sql"

if [ -f "$DB_FILE" ] ; then
    rm "$DB_FILE"
fi

sqlite3 $DB_FILE < $SQL_FILE
