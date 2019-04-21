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
raspivid -v -w 320 -h 240 -fps 12 -n -md 7 -ih -t 0 -o - | cvlc -vvv stream:///dev/stdin --sout '#http{mux=ts,dst=:3717}' :demux=h264
```

=================================================


raspivid -v -w 320 -h 240 -fps 12 -n -md 7 -ih -t 0 -o - | cvlc -I http -vvv stream:///dev/stdin :sout=#transcode{vcodec=theo,vb=512,scale=1}:http{mux=ogg,dst=:3717/stream.ogg} :demux=h264



raspivid -v -w 320 -h 240 -fps 12 -n -md 7 -ih -t 0 -o - | ffmpeg -i - -c:v libx264 -preset ultrafast -tune zerolatency -crf 0 -f mpegts udp://192.168.1.112:5000

ffmpeg -s 320x240 -f video4linux2 -i /dev/video0 -f mpeg1video -b 800k -r 15 udp://192.168.1.112:5000

cvlc screen:// :screen-fps=25 :screen-caching=5000 :sout=#transcode{vcodec=theo,vb=800,scale=1,width=800,height=600,acodec=none}:http{mux=ogg,dst=:3717} :no-sout-rtp-sap :no-sout-standard-sap :ttl=1 :sout-keep


raspivid -w 320 -h 240 -fps 15 -o - | cvlc -vvv stream:///dev/stdin --sout 'http{mux=ts,dst=:3717}' :demux=h264


https://github.com/phoboslab/jsmpeg

raspivid -w 320 -h 240 -fps 15 -o - | cvlc -vvv stream:///dev/stdin --sout '#standard{access=http,mux=ts,dst=:8090}' :demux=h264


raspivid -v -a 524 -a 4 -a "rpi-0 %Y-%m-%d %X" -fps 15 -n -md 2 -ih -t 0 -l -o tcp://0.0.0.0:5001

raspivid -w 320 -h 240 -fps 15 -o - -t 99999 | cvlc -vvv stream:///dev/stdin --sout '#standard{access=http,mux=ts,dst=:3717}' :demux=h264

raspivid -w 320 -h 240 -fps 15 -l -o tcp://0.0.0.0:3717

raspivid -w 320 -h 240 -fps 15 -o - | cvlc -vvv stream:///dev/stdin --sout '#standard{access=http,mux=ts,dst=:8090}' :demux=h264


ffmpeg -f v4l2 -framerate 25 -video_size 640x480 -i /dev/video0 -f mpegts -codec:v mpeg1video -s 640x480 -b:v 1000k -bf 0 

raspivid -o - -t 0 -n -w 320 -h 240 -fps 12 | cvlc -vvv stream:///dev/stdin --sout '#rtp{sdp=rtsp://:8554/}' :demux=h264
