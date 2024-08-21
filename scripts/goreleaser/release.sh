if [ ! -f ".release-env" ]; then
  printf "\033[91m.release-env is required for release\033[0m";
  exit 1;
fi

docker run \
--rm \
-e CGO_ENABLED=1 \
-e GITHUB_TOKEN="$GITHUB_TOKEN" \
-v /var/run/docker.sock:/var/run/docker.sock \
-v `pwd`:/go/src/"$PACKAGE_NAME" \
-v `pwd`/sysroot:/sysroot \
-w /go/src/"$PACKAGE_NAME" \
ghcr.io/goreleaser/goreleaser-cross:"$GOLANG_CROSS_VERSION" \
release --clean