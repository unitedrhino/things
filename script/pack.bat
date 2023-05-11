cd ../src/apisvr
rem "update backend"
go mod tidy
go build
mkdir  ../../pack
mkdir  ../../pack/linux
cp -rf apisvr dist etc ../../pack/linux
cd ../../script
rem "update front"
call buildFront.bat
cp -rf ../assets/dist/* ../pack/linux/dist/front/iThingsCore
cd ../pack
rm -rf iThings-linux.tgz
tar -zcvf iThings-linux.tgz ./linux/*
cd ../script
