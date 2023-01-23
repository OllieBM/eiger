TODO: add chunksize to delta file (only inmportant for the applying of deltas so skipping for now)

# Eiger-Diff
A diff application using a Rolling Hash Algorithm to create an instruction file with the transformations needed to convert one file to another.
I used the basic rolling checksum implementation explained in the [rsync algorithm](https://rsync.samba.org/tech_report/node3.html) whitepaper.
there are some more performant examples based closer on the Adler32, but I didn't feel the improvements would be worth spending more time on this 
example. 

## Build and Run
```
#run using go
go run main.go ...<usage>

# build an executable 
go build -o bin/eiger-diff

./bin/eiger-diff ...<usage>
```

### Usage
```
Diff File1 against File2 creating a diff file with instructions on how to transform File1 into File2

Usage:
  diff File1 File2 [flags]

Flags:
  -b, --blocksize uint32   the size of chunks in bytes to use when matching data from the files max is 0 < b <=5552 (default 4)
  -h, --help               help for diff
  -l, --loglevel string    log level to display {DEBUG|INFO|ERROR} default=ERROR (default "ERROR")
  -o, --output string      optional file to write output to
```

## Requirements
- Hashing function gets the data as a parameter. Separate possible filesystem operations.
- Chunk size can be fixed or dynamic, but must be split to at least two chunks on any sufficiently sized data.
- Should be able to recognize changes between chunks. Only the exact differing locations should be added to the delta.( this would mean not adding references where it is identical positions)
- Well-written unit tests function well in describing the operation, no UI necessary.

## Checklist
- [X] Input/output operations are separated from the calculations
- [X] detects chunk changes and/or additions
  - moved chunks, and additions will be in the delta
  - TODO: link to testcase
- [X] detects chunk removals
  - any chunk in source that is not in target will not be included in the delta
- [X] detects additions between chunks with shifted original chunks
  - delta will contain instructions on creating target using chunks, and additional data
  - TODO: link to testcase