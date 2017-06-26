# Build all the binaries
make all

# build dockerfiles and run them
docker-compose build
docker-compose up
