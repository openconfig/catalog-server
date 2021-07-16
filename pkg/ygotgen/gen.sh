git clone https://github.com/openconfig/ygot.git
cp ygotgen.sh ygot/exampleoc

# cp ygotgen.sh into ygot directory
# execute that bash script and cp back ygotgen file
cd ygot/exampleoc
bash ygotgen.sh
cp ygotgen.go ../../

cd ../..
rm -rf ygot