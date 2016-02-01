#under tor 
# torify curl -L http://check.torproject.org | grep Congratulations
# torify ./gobot
daemonize -p /var/run/GoBot.pid -l /var/lock/subsys/GoBot -u nobody gobot
