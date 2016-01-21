#!/bin/sh

nsqlookupd="nsqlookupd"
nsqd="nsqd"
nsqadmin="nsqadmin"
log_file="nsqd.log"
DEV_MODE="debug"

out_file="/dev/null"

start() {
	mkdir ".nsqd_data" 2>${out_file} 1>${out_file}
	if [ "$DEV_MODE" = "debug" ]; then
		$nsqlookupd 2>${out_file} 1>${out_file} &                                         
		$nsqadmin --lookupd-http-address=127.0.0.1:4161 2>${out_file} 1>${out_file} & 
		$nsqd  -tcp-address=0.0.0.0:4150 -http-address=0.0.0.0:4151 -data-path ".nsqd_data" -lookupd-tcp-address=127.0.0.1:4160  2>${out_file} 1>${out_file}  &      
		#$nsqd -mem-queue-size=100000 -tcp-address=0.0.0.0:4152 -http-address=0.0.0.0:4153 -data-path "/var/lib/nsq/.nsqd_data" -lookupd-tcp-address=127.0.0.1:4160  2>${out_file} 1>${out_file}  &      
		#$nsqd -mem-queue-size=100000 -tcp-address=0.0.0.0:4154 -http-address=0.0.0.0:4155 -data-path "/var/lib/nsq/.nsqd_data" -lookupd-tcp-address=127.0.0.1:4160  2>${out_file} 1>${out_file}  &      
		#$nsqd -mem-queue-size=100000 -tcp-address=0.0.0.0:4156 -http-address=0.0.0.0:4157 -data-path "/var/lib/nsq/.nsqd_data" -lookupd-tcp-address=127.0.0.1:4160  2>${out_file} 1>${out_file}  &      
		#$nsqd -mem-queue-size=100000 -tcp-address=0.0.0.0:4158 -http-address=0.0.0.0:4159 -data-path "/var/lib/nsq/.nsqd_data" -lookupd-tcp-address=127.0.0.1:4160  2>${out_file} 1>${out_file}  &      
		#$nsqd -mem-queue-size=100000 -tcp-address=0.0.0.0:4148 -http-address=0.0.0.0:4149 -data-path "/var/lib/nsq/.nsqd_data" -lookupd-tcp-address=127.0.0.1:4160  2>${out_file} 1>${out_file}  &      
		#$nsqd -mem-queue-size=100000 -tcp-address=0.0.0.0:4146 -http-address=0.0.0.0:4147 -data-path "/var/lib/nsq/.nsqd_data" -lookupd-tcp-address=127.0.0.1:4160  2>${out_file} 1>${out_file}  &      
		## log use
		#$nsqd -mem-queue-size=100000 -tcp-address=0.0.0.0:4144 -http-address=0.0.0.0:4145 -data-path "/var/lib/nsq/.nsqd_data" -lookupd-tcp-address=127.0.0.1:4160  2>${out_file} 1>${out_file}  &      
	else
		$nsqd -data-path ".nsqd_data" &        
	fi
}

stop() {
	killall "$nsqd" 2>${out_file}
	if [ "$DEV_MODE" = "debug" ]; then
		killall "$nsqlookupd" 2>${out_file}
		killall "$nsqadmin" 2>${out_file}
	fi
}

restart() {
	stop
	start
}
status_p() {
	status "$nsqlookupd" 
	status "$nsqd"       
	status "$nsqadmin"   
}

case "$1" in 
	"start")
		start
		;;
	"stop")
		stop
		;;
	"restart")
		restart
		;;
	*)
		echo "Usage: $0 {start|stop|restart}"
		exit 2
		;;
esac




