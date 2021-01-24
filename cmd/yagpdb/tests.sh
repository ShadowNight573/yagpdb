git fetch upstream
git merge upstream/upPrep
sh build.sh && sh copytemplates.sh && ./yagpdb -all -backgroundworkers -pa -syslog