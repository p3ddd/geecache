trap "rm cache_server;kill 0" EXIT

go build -o cache_server.exe
./cache_server.exe -port=8001 &
./cache_server.exe -port=8002 &
./cache_server.exe -port=8003 -api=1&

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &

wait
