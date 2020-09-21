version=`echo $GITHUB_REF | cut -d "/" -f 3`

go build -v -ldflags "-X main.version=$version" .
chmod +x ./care-cli

wget http://gosspublic.alicdn.com/ossutil/1.6.19/ossutil64
chmod 755 ossutil64
./ossutil64 config -e oss-cn-beijing.aliyuncs.com -i $ACCESSKEY -k $SECRETKEY
./ossutil64 rm -rf oss://care-bin/linux --include "*"
./ossutil64 cp -f ./care-cli "oss://care-bin/linux/$version"