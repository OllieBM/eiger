
for an source "Hello" and a target "HelloWorld" 
and a chunksize of 5 we create a delta file
which will contain the operations to change  source -> target
e.g.
"Hello" -> "HelloWorld"
we expect an delta to 'recognize' the use of as much existing data as possible
 
 in Source "Hello"
 vs Target "HelloWorld"
 and block length 5
 we can reuse the block "Hello" situated at block offset 0
 and then we need an operation to add the missing data "World"

 so the delta output will be 
 ```
 = BLOCK_0 
 + 5 World
 ```
 `=` signifys a Match operation and a block index to the source file which can be reused
 `+ ? `  is a MissOperation where we add a continuous string of ? characters 

applying a patch: 
out of scope for this project, 
but since the block count doesn't care for whitespace characters
each INSTRUCTUION/LINE of the delta file should be treated as it should be concatenated

## Match operations "= BLOCK_<??>"
match operations should be treated as a way of saying, you can recreate this file using the data of len(block size) at offset(blocknumber * blocksize)


## Miss operations "+ S ????? "
miss operations are used when the source file is missing some information from the target file, this data needs to be transfered from the targetfile/client to the source file/client and then a new file can be constructed to match the new file
miss operations use the op identifier "+" to signal data must be added followed by a whitespace character and then a length character continuous stream of characters Length  characters

TODO: Matches could have a range, so that we could merge multiple match operations into one if they are contigious