set buildOs=%1
cd ../src/apisvr
rem "update backend"
go mod tidy
go build
mkdir  ..\..\pack\%buildOs%
cp -rf apisvr dist etc ../../pack/%buildOs%
cd ../../script
rem "update front"
call buildFront.bat
cp -rf ../assets/dist/* ../pack/%buildOs%/dist/front/iThingsCore
cd ../pack
rm -rf iThings-%buildOs%.tgz
tar -zcvf iThings-%buildOs%.tgz ./%buildOs%/*
cd ../script
