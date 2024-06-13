#!/bin/bash

# Check if enough arguments are provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <output_dir>"
    exit 1
fi

# Argument
output_dir="$1"

# Process each result file in the output directory
for file in "$output_dir"/result_*.txt; do
    awk '{ count[$0]++ } END { for (id in count) if (count[id] > 1) print id }' "$file"
done
