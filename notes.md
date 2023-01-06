reading list:
https://rsync.samba.org/tech_report/node5.html
https://iq.opengenus.org/rolling-hash/



1. 'create a chunk size'



2. create a length N (min(default, min(file1.len(), file2.len()))



// for a chunk size
we compare WINDOWS/CHUNKS, if a chunk matches it is ignored
if a chunk FAILS then it is added to the list of failed chunks

failed chunks can then be rechecked with different window sizes?


https://github.com/chmduquesne/rollinghash


# rdiff
https://librsync.github.io/delta_8c_source.html

rsync
https://en.wikipedia.org/wiki/Rsync#Algorithm



1. the old string gets split into CHUNKSIZE
2. each chunk size is then stored in a map
3. the new string uses the rolling hash to check for collisions in map
4. 
