
for pid in $(ps -fe | grep "./main stockfish" | grep -v grep | awk '{print $2}'); do
    sudo kill -s 15 "$pid"
done