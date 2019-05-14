## Profile Setup

1. Update & Upgrade
1. Set video ram to highest split
1. Set hostname
1. Expand filesystem
1. Enable GPIO
1. Enable Camera (enable v4l2?)
1. Enable WiFi w/ creds

```
sudo apt-get update && sudo apt-get upgrade
sudo raspi-config nonint do_memory_split 256
sudo raspi-config nonint do_hostname marv
sudo raspi-config nonint do_camera 1
sudo raspi-config nonint do_expand_rootfs
sudo raspi-config nonint do_wifi_country US
```

1. Install Docker (curl -sSL https://get.docker.com | sh)

## Hardware Config

1. Fritzing


## Running

```bash
docker pull danielmconrad/rover:latest
docker run -dit --privileged --restart unless-stopped -p 3737:3737 danielmconrad/rover:latest 
```
