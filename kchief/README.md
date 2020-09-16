# Kchief simulator

This package contains 4 tools:

- protoencoder <https://github.com/RaaLabs/shipsimulator/tree/master/kchief/cmd/protoencoder>
- protodecoder <https://github.com/RaaLabs/shipsimulator/tree/master/kchief/cmd/protodecoder>
- mqttpublish <https://github.com/RaaLabs/shipsimulator/tree/master/kchief/cmd/mqttpublish>
- mqttsubscribe <https://github.com/RaaLabs/shipsimulator/tree/master/kchief/cmd/mqttsubscribe>

Each respective package are described below.

## protoencoder

Flags:

```bash
  -inJsonFile string
    specify the full path of the json file to take as input
  -outFile string
    specify the filename of the output file (default "myfile.bin")
```

protoencoder let's you generate a Khief proto file based on a JSON file describing tags and values. Example of JSON file below.

```json
{
    "dataPoints": [{
            "tagName": "MAN-TAGER-HVA-MAN-HAVER",
            "value": 1.1
        },
        {
            "tagName": "ANTAGER",
            "value": -2.2
        },
        {
            "tagName": "SINNATAG",
            "value": 1
        }
    ]
}
```

The output will be a binary file containing protobuf data.

## protodecoder

Flags:

```bash
  -fileName string
    The full path with the filename of the Kchief protobuf data file (default "./sample/sample2.bin")
```

protodecode will let you decode a binary protobuf, unmarshal it from protobuf format into ascii key/values, and print the output to console.

## mqttpublish

Flags:

```bash
  -broker string
    The ip address of the MQTT broker (default "10.0.0.26")
  -clientID string
    The client ID to use with MQTT (default "btclient1")
  -delay int
    The number of milliseconds to wait between each mqtt publish (default 300)
  -inFile value
    specify the files to use as input comma separated
  -port string
    The port where the MQTT broker listens (default "1883")
  -protocol string
    The protocol to use when connecting to the MQTT broker (default "tcp")
  -repetitions int
    specify how many repetitions to run (default 1)
  -topic string
    The name of the MQTT topic (default "CloudBoundContainer")
```

mqttpublish will let you publish data like for example protobuf binary data to a MQTT broker. The package let's you specify more infiles by repeating the use of the `--inFile` flag. Repetitions and delay between each send can be tuned with flags. Example below.

```bash
go run main.go --inFile=./ship1.bin --inFile=./ship2.bin --repetitions=1000 --delay=1
```

## mqttsubscribe

Flags:

```bash
  -broker string
    The ip address of the MQTT broker (default "10.0.0.26")
  -clientID string
    The client ID to use with MQTT (default "btclient2")
  -port string
    The port where the MQTT broker listens (default "1883")
  -protocol string
    The protocol to use when connecting to the MQTT broker (default "tcp")
  -topic string
    The name of the MQTT topic (default "CloudBoundContainer")
```

mqttsubscribe will let you subscribe to a MQTT topic on a specified broker. The messages received will be saved in it's binary format and the naming will be numbered incrementally for each package received.

`protodecoder` can then be used to decode the files of the received messages to verify their correctness.
