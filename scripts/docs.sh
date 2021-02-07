function godoc() {
  if [ ! -f go.mod ]
  then
    echo "error: go.mod not found" >&2
    return
  fi

  module=$(sed -n 's/^module \(.*\)/\1/p' go.mod)
  sudo docker run \
    --rm \
    -e "GOPATH=/tmp/go" \
    -p 127.0.0.1:6060:6060 \
    -v $PWD:/tmp/go/src/$module \
    golang \
    bash -c "go get golang.org/x/tools/cmd/godoc && echo http://localhost:6060/pkg/$module && /tmp/go/bin/godoc -http=:6060"
}

godoc
