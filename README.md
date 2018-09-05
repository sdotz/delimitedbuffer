# delimitedbuffer

A library to read and write streams of delimited blobs, as described in: https://developers.google.com/protocol-buffers/docs/techniques#streaming
Can be used for protocol buffer binaries, or any other data. Blobs are preceeded by their 4-byte size (max 4294967295 bytes)

Embeds bytes.Buffer and implements the common interfaces to support chained readers/writers for to handle one full datum at a time.

## Ideas/Future plans
- Option to write a header before every datum
    - Checksum for the datum
    - Content encoding (gzip etc)
    - Byte-size of the datum to follow contained in header
