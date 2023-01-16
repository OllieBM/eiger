# explanation & attempted proof

using the example in ./resources/test_input/numbers_follower.txt
```
12346789
12346789
12346789
12346789
```
assuming a chunk_size of '5' the signature file will be generated to match  these sequence of characters
```

1234|6789|
123|4678|9
12|3467|89
1|2346|789 

== [1234][6789][\n123][4678][9\n12][3467][89\n1][2346][789]
```

the diff generation will then try to match blocks of 5 characters from the leader file

leader file
```
123456789
123456789
123456789
123456789```

```
// some of these values match completely so could be reduced if we are confident in our hashing uniqueness
// otherwise use a hashmap
[1234][5678][9\n12][3456][789\n][1234][5678][9\n12][3456][789]
```
so we compare the <LEADER> to the <FOLLOWER> using a rolling window of blocksize

follower:
`12346789\n12346789\n12346789\n12346789`
`[1234][6789][\n123][4678][9\n12][3467][89\n1][2346][789]`


leader iteration:
1. [1234] = match block 0
2. [5678] = miss // missing '5'
3. [6789] = match block 1
4. [\n123] = match block 2
5. [4567] = miss // no 5 exists [add 4] 
6. [5678] = miss // miss no 5 exists [add 5]
7. [6789] = match block 1
8. [\n123] = match block 2
9. [4567] = miss // no 5 exists [add 4] 
10. [5678] = miss // miss no 5 exists [add 5]
11. [6789] = match block 1
12. [\n123] = match block 2
13. [4567] = miss // no 5 exists [add 4] 
14. [5678] = miss // miss no 5 exists [add 5]
15. [6789] = match block 1

`123456789\n123456789\n123456789\n123456789`
`[1234][5678][9\n12][3456][789\n][1234][5678][9\n12][3456][789]`

with a different chunk size we can get different results



### follower to leader
follower:
12346789\n
12346789\n
12346789\n
12346789\n
leader
123456789\n
123456789\n
123456789\n
123456789\n

we want to add 1 5 to every line
then re use blocks 0, & 1?



### leader to follower
follower:
12346789\n
12346789\n
12346789\n
12346789\n
leader
123456789\n
123456789\n
123456789\n
123456789\n
// since we don't have a removal option
we want to read chunks and add in as few characters as possible
so we can use chunk 0 [1:4], add six and use another chunk (BLOCK_4) [789\n]

