go build
git clone https://github.com/openconfig/public.git
# Copy ietf directory for backup usage.
cp -r ./public/third_party/ietf .
cd public

times = 0
commits=`git rev-list --all`
for commit in $commits
do
    delete=false
    let "times=times+1"
    git checkout $commit
    path="./public/release/models,./public/third_party/ietf"

    # If `third_party/ietf` directory is not included in this commit, copy it from previous backup.
    if [ ! -d "third_party/ietf" ]; then
        mkdir third_party
        cp -r ../ietf ./third_party
        delete=true
    fi
    cd ..
    echo "./goyang -p $path -u https://github.com/openconfig/public/blob/$commit/"

    ./goyang -p $path -u https://github.com/openconfig/public/blob/$commit/
    cd public
    # Track how many commits the script has crawled.
    echo $times
    if $delete; then
        rm -rf ./third_party
    fi
    sleep 10
done

cd ..
rm -rf public
rm -rf ietf