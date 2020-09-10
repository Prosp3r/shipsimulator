# NMEA Generator

The binary in the repository is compiled for amd64 linux architecture.

The flags available are:

```bash
 -address string
    The network host and port to send to, like localhost:8888 (default "localhost:8888")
  -delay int
    The delay to wait between each send of data given in Micro Seconds. Default is 1000000 (1 Second) (default 1000000)
  -file string
    The name of the the NMEA file to read (default "./output.nmea")
  -loop
    loop over again, and again, and again, and again,...........
```
