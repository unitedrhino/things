echo "start update front"
git submodule update --init --recursive
git submodule foreach git checkout master
git submodule foreach git pull
echo "end update front"