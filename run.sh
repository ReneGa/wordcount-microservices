# Build all the binaries
export GOOS=linux

for component in ingestor recorder searches stopwordfilter windower wordcounter; do
    echo "building $component"
    cd $component
    go build
    chmod +x $component
    cd ..
done

# build dockerfiles and run them
docker-compose up
