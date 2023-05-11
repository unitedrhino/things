echo "start update front"
git submodule update --init --recursive
git submodule foreach git checkout test
git submodule foreach git pull
echo "end update front"