#!/bin/bash

set -e

GREEN='\033[0;32m'
YELLOw='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m'

hostname="rover"
wifi_country="US"
wifi_ssid=""
wifi_pass=""

main() {
  ask_for_inputs
  set_configurations
  install_docker
  install_rover
  success "Successfully configured your Pi as a rover!"
}

ask_for_inputs() {
  ask "Hostname? ($hostname)"
  read hostname_in
  if [[ $hostname_in != "" ]]; then hostname=$hostname_in; fi

  ask "WiFi Country? ($wifi_country)"
  read wifi_country_in
  if [[ $wifi_country_in != "" ]]; then wifi_country=$wifi_country_in; fi

  ask "WiFi SSID?"
  read wifi_ssid

  ask "WiFi Password?"
  read wifi_pass
}

set_configurations() {
  info "Expanding file system"
  out=$(sudo raspi-config nonint do_expand_rootfs)

  info "Setting video memory to 256MB"
  out=$(sudo raspi-config nonint do_memory_split 256)

  info "Enabling camera"
  out=$(sudo raspi-config nonint do_camera 1)
  
  info "Setting hostname"
  out=$(sudo raspi-config nonint do_hostname $hostname)
  
  info "Setting WiFi country"
  out=$(sudo raspi-config nonint do_wifi_country $wifi_country)
  
  info "Setting WiFi credentials"
  set_wifi_credentials
}

install_docker() {
  info "Installing docker"
  out=$(curl -sSL https://get.docker.com | sh)
  out=$(sudo usermod -a -G docker $USER)

  info "Waiting for docker to start"
  until pgrep -f docker > /dev/null; do sleep 1; done
}

install_rover() {
  info "Installing rover docker image"
  out=$(sudo docker pull danielmconrad/rover:latest)
  out=$(sudo docker run -dit --privileged --restart unless-stopped -p 3737:3737 danielmconrad/rover:latest )
}

set_wifi_credentials() {
  sudo bash -c "cat >> /etc/wpa_supplicant/wpa_supplicant.conf" <<EOF
network={
  ssid="$wifi_ssid"
  psk="$wifi_pass"
}
EOF
}

ask() {
  printf "${MAGENTA}[ROVER]${NC} $1 "
}

info() {
  printf "${CYAN}[ROVER]${NC} $1\n"
}

success() {
  printf "${GREEN}[ROVER]${NC} $1\n"
}

main
