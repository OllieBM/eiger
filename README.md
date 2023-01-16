
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
        eiger-diff --file1 <File1> --file2 <File2> [flags]

Flags:
  -b, --blocksize uint32   the size of chunks in bytes to use when matching data from the files (default 5)
      --file1 string       The first file to read
      --file2 string       The second file to read
  -h, --help               help for Diff
  -l, --loglevel string    log level to display {DEBUG|INFO|ERROR} default=ERROR (default "ERROR")
```


# Rolling Hash Algorithm
_Spec v4 (2021-03-09)_

Make a rolling hash based file diffing algorithm. When comparing original and an updated version of an input, it should return a description ("delta") which can be used to upgrade an original version of the file into the new file. The description contains the chunks which:
- Can be reused from the original file
- have been added or modified and thus would need to be synchronized

The real-world use case for this type of construct could be a distributed file storage system. This reduces the need for bandwidth and storage. If many people have the same file stored on Dropbox, for example, there's no need to upload it again.

A library that does a similar thing is [rdiff](https://linux.die.net/man/1/rdiff). You don't need to fulfill the patch part of the API, only signature and delta.

## Requirements
- Hashing function gets the data as a parameter. Separate possible filesystem operations.
- Chunk size can be fixed or dynamic, but must be split to at least two chunks on any sufficiently sized data.
- Should be able to recognize changes between chunks. Only the exact differing locations should be added to the delta.
- Well-written unit tests function well in describing the operation, no UI necessary.

## Checklist
- [ ] 1. Input/output operations are separated from the calculations
- [ ] 2. detects chunk changes and/or additions
- [ ] 3. detects chunk removals
- [ ] 4. detects additions between chunks with shifted original chunks



TODO:
    - [ ] calculation functions
    - [ ] break input into chunks
    - [ ] detect chunk addition
      - [ ] beginning
      - [ ] middle
      - [ ] end
    - [ ] 




########### 
large chunk simple hash check
localised rolling hash algorithm


2. use rolling hash to detect file differences
3. use 


https://prettydiff.com/2/guide/unrelated_diff.xhtml


1. convert a file into a signature
   1. signature 
  