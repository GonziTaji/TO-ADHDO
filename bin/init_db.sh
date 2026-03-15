#!/usr/bin/env bash
set -euo pipefail

NOW=$( date '+%F_%H:%M:%S' )

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

ROOT="$SCRIPT_DIR/.."

DB_FILE="$ROOT/database/main.db"
SCHEMA_FILE="$ROOT/sql/schema.sql"
INITIAL_DATA_FILE="$ROOT/sql/test_data.sql"

if [ -f "$DB_FILE" ] ; then
    mv "$DB_FILE" $ROOT/database/bk_$NOW.db
fi

sqlite3 $DB_FILE < $SCHEMA_FILE
sqlite3 $DB_FILE < $INITIAL_DATA_FILE
