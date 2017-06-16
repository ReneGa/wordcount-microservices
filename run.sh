# Build all the binaries
export GOOS=linux

cd ingestor
go build
chmod +x ingestor
cd ..

cd recorder
go build
chmod +x recorder
cd ..

cd searches
go build
chmod +x searches
cd ..

cd stopwordfilter
go build
chmod +x stopwordfilter
cd ..

cd windower
go build
chmod +x windower
cd ..

cd wordcounter
go build
chmod +x wordcounter
cd ..

# build dockerfiles and run them
docker-compose up
