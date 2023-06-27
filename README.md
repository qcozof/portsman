# portsman
Provide a simple http service to check ports are open.

## Samples

+ ### Listen one port
```shell
./portsman --ports 9090
```

+ ### Listen multiple ports
```shell
./portsman --ports 9090,9091,9092
```

+ ### Specify domain
```shell
./portsman --ports 9090 --domain sample.com
```

+ ### Specify web dir
```shell
./portsman --ports 9090,9091,9092 --webDir /path/to
``````

+ ### Enable ssl and use default ssl key and cert file
```shell
./portsman --ports 9090,9091,9092 --enableSsl true
`````````

+ ### Specify ssl key and cert file will enable ssl automatically
```shell
./portsman --ports 9090,9091,9092 --certFile keys/fullchain.cer --keyFile keys/sample.com.key
```

* *In the Windows operating system, replace `./portsman` with `portsman.exe` before running the command.*

Then you can visit the URL to check if the port is open.
