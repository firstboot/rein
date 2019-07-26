cd bin
del -r *
cd ..

go build -o ./bin/rein-1.0.4.exe -v ./src/

.\bin\rein-1.0.4.exe