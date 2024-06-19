#!/bin/bash

# Check number of arguments
if [ "$#" -ne 7 ]; then
  echo "Usage: $0 <username> <password> <host> <dbname> <table_name> <input_file> <output_file>"
  exit 1
fi

# Parse arguments
DB_USER="$1"
DB_PASSWORD="$2"
DB_HOST="$3"
DB_NAME="$4"
TABLE_NAME="$5"
INPUT_FILE="$6"
OUTPUT_FILE="$7"

# Function to execute MySQL query for each ID
function query_ids {
  local id="$1"
  local query="SELECT * FROM ${TABLE_NAME} WHERE id = ${id};"
  mysql -u"${DB_USER}" -p"${DB_PASSWORD}" -h"${DB_HOST}" -D"${DB_NAME}" -sN -e "${query}" 2>/dev/null
}

# Read IDs from input file and query MySQL for each ID
while IFS= read -r id || [[ -n "$id" ]]; do
  #echo "$id"
  result=$(query_ids "0x$id")
  echo "$result" >> "$OUTPUT_FILE"
done < "$INPUT_FILE"

# Calculate MD5 checksum of output file
md5sum "$OUTPUT_FILE"