# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

go build
git clone https://github.com/openconfig/public.git
# Copy ietf directory for backup usage.
cp -r ./public/third_party/ietf .
cd public

times = 0
# This retrives all commits in reverse chronological order to mitigate bad version numbers in older commits poisoning the data.
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
    echo "./crawl -p $path -u https://raw.githubusercontent.com/openconfig/public/$commit/"

    ./crawl -p $path -u https://raw.githubusercontent.com/openconfig/public/$commit/
    cd public
    # Track how many commits the script has crawled.
    echo $times
    if $delete; then
        rm -rf ./third_party
    fi
    # Sleep for a while to release previous connection with DB, to avoid connection failure in the next run.
    sleep 10
done

cd ..
rm -rf public
rm -rf ietf