echo "start build front"
call updateFront.bat
cd ../assets
call yarn install
call yarn build
cd ../script

