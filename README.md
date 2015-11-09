This is a simple snappy application controlling a piglow
device attached to a raspberry pi2. It checks with http://juju.fail
and determines the status of the juju master branch and sets the piglow color
accordingly (GREEN -> unblocked, RED -> blocked).

To deploy this on your own raspberry pi2 running
[Ubuntu Snappy Core](https://developer.ubuntu.com/en/snappy/),
do:

```
$ make
$ snappy-remote --url=ssh://[IP-of-our-RPi2]:22 install blink_1.0.0_all.snap
```
