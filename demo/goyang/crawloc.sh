go build
git clone https://github.com/openconfig/public.git
./goyang -p ./public/release/models/,./public/third_party/ietf/
rm -rf public