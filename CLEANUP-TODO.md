CLEANUP :

- [ ] create cmd example
- [ ] encapsulation
  - [ ] hasher
    - [ ] combine strong hasher & rolling checksum into interface
    - [ ] rolling checksum can take a stream
      -  [ ] can keep trying to read N bytes from stream
         -  [ ] Roll - read a byte and rotate array
         -  [ ] Calculate - consume another block size
   -  [ ] signature into a type with things like block len
   -  [ ] operationwriter should bean interface
    - 