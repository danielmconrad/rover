#!/bin/sh
set -e

GREEN='\033[0;32m'
YELLOw='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m'

main() {
	set_configurations
	install_docker
	install_rover
	success "Successfully configured your Pi as a rover! Please restart the pi."
}

set_configurations() {
	info "Resizing video memory"
	out=$(sudo raspi-config nonint do_memory_split 256)

	info "Enabling camera"
	out=$(sudo raspi-config nonint do_camera 1)
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
	out=$(sudo docker run -dit --privileged --restart unless-stopped -p 3737:3737 danielmconrad/rover:latest)
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
