#!/bin/bash

# Check number of arguments
if [ "$#" -lt 6 ]; then
  echo "Usage: $0 <username> <password> <host> <dbname> <table_name> <num_partitions>"
  exit 1
fi

# Parse arguments
DB_USER="$1"
DB_PASSWORD="$2"
DB_HOST="$3"
DB_NAME="$4"
TABLE_NAME="$5"
NUM_PARTITIONS="$6"

# Construct MySQL connection string

# Query to retrieve min and max id as hexadecimal strings
#query="SELECT MIN(id), MAX(id) FROM ${TABLE_NAME};"
# Connect to MySQL and execute the query
query="SELECT HEX(min_id) AS min_id_hex, HEX(max_id) AS max_id_hex FROM (SELECT MIN(id) AS min_id, MAX(id) AS max_id FROM ${TABLE_NAME}) AS subquery;"

read min_hex_uuid max_hex_uuid <<< $(mysql -N -B -e "${query}" -u "${DB_USER}" -p"${DB_PASSWORD}" -h "${DB_HOST}" "${DB_NAME}" 2>/dev/null)
min_id=$(echo "ibase=16; ${min_hex_uuid^^}" | bc)
max_id=$(echo "ibase=16; ${max_hex_uuid^^}" | bc)

# Calculate partition size
# range=$(( ($max_id - $min_id) / $NUM_PARTITIONS ))
range=$(echo "scale=0; (${max_id} - ${min_id} + 1) / ${NUM_PARTITIONS}" | bc)
echo "min $min_hex_uuid, max $max_hex_uuid, range: $range"
# Initialize variables for tracking partitions
start_id=$min_id
end_id=$min_id
total_count=0

# Loop through partitions and count rows
for (( i=1; i<=$((NUM_PARTITIONS+1)); i++ )); do
  # Calculate end_id using awk for safe integer arithmetic
  #end_id=$(awk -v start="$start_id" -v range="$range" 'BEGIN {print start + range}')
  end_id=$(echo "$end_id + $range" | bc)
  start_hex_uuid=$(echo "obase=16;ibase=10;$start_id" | bc)
  end_hex_uuid=$(echo "obase=16;ibase=10;$end_id" | bc)
  # Query to count rows in current partition
  count_query="SELECT COUNT(id) FROM ${TABLE_NAME} WHERE id >= 0x$start_hex_uuid  AND id < 0x$end_hex_uuid;"

  # Execute the count query
  count=$(mysql -N -B -e "${count_query}" -u "${DB_USER}" -p"${DB_PASSWORD}" -h "${DB_HOST}" "${DB_NAME}" 2>/dev/null)

  # Store count in partition_counts array
  total_count=$((total_count+count))

  echo "$i $start_hex_uuid $end_hex_uuid $count $total_count"
  # Print partition info including range

  # Update start ID for next partition
  start_id=$end_id
done