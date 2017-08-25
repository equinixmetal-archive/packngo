# packngo
Packet Go Api Client

![](https://www.packet.net/media/labs/images/1679091c5a880faf6fb5e6087eb1b2dc/ULY7-hero.png)

Committing
----------

Before committing, it's a good idea to run `gofmt -w *.go`. ([gofmt](https://golang.org/cmd/gofmt/))

Acceptance Tests
----------------

If you want to run tests against the actual Packet API, you must set envvar `PACKET_TEST_ACTUAL_API` to non-empty string for the `go test` run, e.g.

```
$ PACKET_TEST_ACTUAL_API=1 go test -v
```
