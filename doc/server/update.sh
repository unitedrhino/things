cd iThings-pkg
tar -xvzf iThings-linux.tgz
rm -rf ../run/dist/* ../run/papisvr ../run/etc
cp -rf linux/apisvr linux/dist linux/etc ../run/
cd ../run
chmod 751 ./apisvr
killall apisvr
nohup ./apisvr &
tail -f nohup.out