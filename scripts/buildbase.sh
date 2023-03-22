systemctl start docker
docker-compose -f ./build/docker-compose.yml  -p "k2" stop
docker-compose -f ./build/docker-compose.yml  -p "k2" down
docker-compose -f ./build/docker-compose.base.yml  -p "k2" up -d 