pid=`ps -ef|grep fc_go|grep -v grep|awk '{print $2}'`
sudo kill -9 $pid
