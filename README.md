# go

my go lib


## database/pqto 

`database/pqto` is an extension package of  [pq](github.com/lib/pq) (a pure Go postgres driver for Go's database/sql package).It introduces read and write timeouts in database soruce connection string as follows:

```
postgres://192.168.15.130:5432/postgres?sslmode=disable&connect_timeout=2&read_timeout=2s&write_timeout=2s
```

* read_timeout: It uses a format similar to Go's Duration to describe the read timeout.

* write_timeout: It uses a format similar to Go's Duration to describe the writetimeout.

## encoding

`encoding` provides a fast JSON parsing solution that is different from the one in the standard library.

## gmsm

`gmsm` provides encryption methods compliant with national cryptographic standards.

## log

`log` provides a simple logging wrapper.

## time2

`time2 provides`Go's Duration json format

## timeout

`timeout` provided a connection with read and write timeouts.