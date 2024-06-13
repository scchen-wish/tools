import uuid

def increment_clock_seq(uuid_str):
    # Remove '0x' prefix if present
    if uuid_str.startswith("0x"):
        uuid_str = uuid_str[2:]

    # Parse the UUID fields
    original_uuid = uuid.UUID(uuid_str)
    time_low = original_uuid.time_low
    time_mid = original_uuid.time_mid
    time_hi_version = original_uuid.time_hi_version
    clock_seq = original_uuid.clock_seq
    node = original_uuid.node

    # Increment the clock sequence by 1
    new_clock_seq = (clock_seq + 1) & 0x3FFF  # Ensure it stays within 14 bits

    # Create a new UUID with the incremented clock sequence
    new_uuid = uuid.UUID(fields=(time_low, time_mid, time_hi_version,
                                 new_clock_seq >> 8, new_clock_seq & 0xFF, node), version=1)

    return f"0x{new_uuid.hex}"

def process_file(input_file, output_file):
    with open(input_file, 'r') as infile, open(output_file, 'w') as outfile:
        for line in infile:
            original_uuid = line.strip()
            new_uuid = increment_clock_seq(original_uuid)
            outfile.write(f"{original_uuid}\t{new_uuid}\n")

if __name__ == "__main__":
    import sys
    if len(sys.argv) != 3:
        print("Usage: python script.py input_file output_file")
        sys.exit(1)

    input_file = sys.argv[1]
    output_file = sys.argv[2]

    process_file(input_file, output_file)
