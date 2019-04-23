set -e

CURDIR=$(cd $(dirname $0); pwd)

export CGO_ENABLED="0"

cd $CURDIR && cd write_api && echo build `pwd` && go build

cd $CURDIR && cd read_api && echo build `pwd` && go build

cd $CURDIR && cd consumer/checksum && echo build `pwd` && go build

cd $CURDIR && cd consumer/delete_event && echo build `pwd` && go build

cd $CURDIR && cd consumer/post_event && echo build `pwd` && go build
