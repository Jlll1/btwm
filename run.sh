XEPHYR=$(whereis -b Xephyr | awk '{print $2}')
xinit ./xinitrc -- "$XEPHYR" :100 -ac -screen 800x600 -host-cursor
