#Shutdown goplaychess server server so letsencrypt can bind to port 80 temporarily
PID=$(ps -ef | grep "sudo ./main goplaychess.com" | grep -v grep | awk '{ printf $2 }')
if [[ "" !=  "$PID" ]]; then
	sudo kill -s 15 $PID 
fi

cd ~/workspace/gochess
go build main.go
nohup sudo ./main goplaychess.com &

echo "Successfully restarted web server!"