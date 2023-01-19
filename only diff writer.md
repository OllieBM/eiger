only diff writer
-> attempt at a unified diff
match == a checksum match at offset match

block 0
block 1
block 2
V
block 0
block 2
= 
- block_2

block 0 
block 1 
block 2
V 
block 0
block 2 
block 1

- block_1
+ @offset R block_1


block 0 
block 1 
block 2
block 3
block 4

V 
block 0
block 2
block 3
block 4
block 1

=
- BLOCK_1
+ @4*bs = BLOCK_1



block 0 
block 1 
block 2
block 3
block 4

V 
block 0
block 3
block 4

- BLOCK_1
- BLOCK_2