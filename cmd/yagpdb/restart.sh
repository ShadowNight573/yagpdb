screen -r
killall yagpdb
git fetch upstream
git merge upstream/shadow
sh build.sh && sh copytemplates.sh && ./yagpdb -all -backgroundworkers -pa -syslog