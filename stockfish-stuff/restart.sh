#Shutdown goplaychess server server so letsencrypt can bind to port 80 temporarily
PID=$(ps -ef | grep "sudo ./main stockfish" | grep -v grep | awk '{ printf $2 }')
if [[ "" !=  "$PID" ]]; then
	sudo kill -s 15 $PID 
fi

go build main.go
nohup ./main stockfish &

echo "Successfully started stockfish!"