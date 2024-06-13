### Process
- dump all ids from leger database
- partition ids to same file with same hash
- find all duplicate id
- generate new uuid and replace duplicate
- dump data with all fields(sql format)
- load data to datacenter

### dumper
https://github.com/lazhang-wish/go-mydumper/tree/jacques_vitess
### loader
https://github.com/lazhang-wish/go-mydumper/tree/jacques_vitess
### partitioner
go build -v -o partitioner ./partitioner.go
./partitioner dumper-sql/ledger_event/ledger.ledger_event.000 output/ledger_event/ 100 16
### count duplicate
./count_duplicate.sh output/ledger_event/
### generate new id
python generate_new_uuid.py input_dir output_dir