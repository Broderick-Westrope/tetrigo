docker run \
--rm \
-e CGO_ENABLED=1 \
--env-file .release-env \
-v /var/run/docker.sock:/var/run/docker.sock \
-v `pwd`:/go/src/"$PACKAGE_NAME" \
-v `pwd`/sysroot:/sysroot \
-w /go/src/"$PACKAGE_NAME" \
ghcr.io/goreleaser/goreleaser-cross:"$GOLANG_CROSS_VERSION" \
release --clean