# Build all the binaries
export GOOS=linux

for component in ingestor stopwordfilter wordcounter; do
    echo "Compiling $component"
    cd $component
    go build
    chmod +x $component
    cd ..
done

# build dockerfiles and run them
docker-compose up
