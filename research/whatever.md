## Necessary Setup

curl -sSL https://get.docker.com | sh
set video split to 256mb
sudo modprobe bcm2835-v4l2

## Worked

```bash
ffmpeg -s 320x240 -f video4linux2 -i /dev/video0 -f mpeg1video -b 800k -r 30 udp://192.168.1.112:5000
```

```bash
raspivid -o - -t 99999 | cvlc -vvv stream:///dev/stdin --sout '#standard{access=http,mux=ts,dst=:8090}' :demux=h264
```

```bash
raspivid -v -w 320 -h 240 -fps 15 -n -md 7 -ih -t 0 -o - | cvlc -vvv stream:///dev/stdin --sout '#standard{access=http,mux=ts,dst=:3717}' :demux=h264
```

=================================================

https://github.com/phoboslab/jsmpeg

raspivid -v -a 524 -a 4 -a "rpi-0 %Y-%m-%d %X" -fps 15 -n -md 2 -ih -t 0 -l -o tcp://0.0.0.0:5001

raspivid -w 320 -h 240 -fps 15 -o - -t 99999 |cvlc -vvv stream:///dev/stdin --sout '#standard{access=http,mux=ts,dst=:3717}' :demux=h264

raspivid -w 320 -h 240 -fps 15 -l -o tcp://0.0.0.0:3717

raspivid -w 320 -h 240 -fps 15 -o - | cvlc -vvv stream:///dev/stdin --sout '#standard{access=http,mux=ts,dst=:8090}' :demux=h264


ffmpeg -f v4l2 -framerate 25 -video_size 640x480 -i /dev/video0 -f mpegts -codec:v mpeg1video -s 640x480 -b:v 1000k -bf 0 
