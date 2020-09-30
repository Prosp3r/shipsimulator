# portforward

portforward will let you set up a local tcp listener at specfied ip and port, and it will forward all connections made to that listener to the specfied remote-ip:port.

```bash
./portforward --help
  -localAddress string
    ipAddress:port (default "127.0.0.1:8080")
  -remoteAddress string
    remoteAddress:port (default "erter.org:80")
```

Precompiled binary for linux amd64 architecture are provided in the repository.
