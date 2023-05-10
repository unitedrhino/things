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
tar -cvf ../pack/iThings-linux.tgz ../pack/linux
SET GOOS=windows
SET GOARCH=amd64
