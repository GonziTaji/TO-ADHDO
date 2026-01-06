#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DB_FILE="$SCRIPT_DIR/../main.db"
SQL_FILE="$SCRIPT_DIR/../sql/create_tables.sql"

if [ -f "$DB_FILE" ] ; then
    rm "$DB_FILE"
fi

sqlite3 $DB_FILE < $SQL_FILE
