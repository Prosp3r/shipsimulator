# NMEA simulator

NMEA packet generator to be used to simulate an NMEA device.

The generator takes a file with a predefined NMEA route as input, and sends the data parsed from the file to the specified network receiver.

## Main program / exexutables

Can be found in this repository located under\
<https://github.com/RaaLabs/shipsimulator/tree/master/nmea/cmd/nmeagenerator>

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

### Run and generate some output locally for testing

Start a TCP network listener on port 8888

`nc -l localhost 8888`

Start a generator that will generate and send packages every 500ms

`go run ./cmd/nmeagenerator/main.go --delay=500000 --file=../cmd/nmeagenerator/output.nmea`



## References

To generate routes use the tool on the web page below.

<https://nmeagen.org/>
