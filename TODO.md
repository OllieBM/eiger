# Improvements
- [ ] add metadata to the diff writer so that diffs can be stored and applied later
  - [ ] add chunk size to diff metadata so that we know the size of chunks used to generate
  - [ ] add source and target file names to metadata
- [ ] Flag options for different Checksums and Hashers
- [ ] benchmark tests

# WIBNI/Easy wins
- [ ] mockgen based tests
- [ ] sanitize names (leader/follower source/target etc)
- [ ] add a verbose option to the diffWriter to add references for matching blocks
- [ ] join the strong Hash and rolling checksum into a single type
- [ ] Use Ring buffer for the rolling checksum
- [ ] refactor diffWriter.go::AddMatch