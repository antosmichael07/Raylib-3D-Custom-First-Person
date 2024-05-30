zip ray.zip ./*
echo "put ray.zip ray.zip" | sftp student@192.168.10.65
echo -e "unzip -o ray.zip; go build" | ssh student@192.168.10.65
echo "get rl-go program" | sftp student@192.168.10.65
./program
