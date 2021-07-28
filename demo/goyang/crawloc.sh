go build
git clone https://github.com/openconfig/public.git
cd public

times = 0
commits=`git rev-list --all`
for commit in $commits
do
    let "times=times+1"
    git checkout $commit
    path=""
    if [ -d "release/models" ]; then
        path="$path./public/release/models"
    fi
    if [ -d "third_party/ietf" ]; then
        if [ $path != "" ]; then
            path="$path,"
        fi
        path="$path./public/third_party/ietf"
    fi
    cd ..
    echo "./goyang -p $path -u https://github.com/openconfig/public/blob/$commit/"
    pwd
    ./goyang -p $path -u https://github.com/openconfig/public/blob/$commit/
    cd public
    echo $times
    sleep 10
done

#./goyang -p ./public/release/models/,./public/third_party/ietf/
cd ..
# rm -rf public