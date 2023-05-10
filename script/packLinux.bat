SET GOOS=linux
SET GOARCH=amd64
cd ../src/apisvr
rem "update backend"
go mod tidy
go build api.go
mkdir -p  ..\..\pack\linux
cp -rf api dist etc ../../pack/linux
cd ../../script
rem "update front"
call buildFront.bat
cp -rf ../assets/dist/* ../pack/linux/dist/front/iThingsCore
cd ../pack
rm -rf iThings-linux.tgz
tar -zcvf iThings-linux.tgz ./linux/*
cd ../script
SET GOOS=windows
SET GOARCH=amd64
