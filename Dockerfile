# checkov:skip=CKV_DOCKER_7: "Ensure the base image uses a non latest version tag"
FROM cgr.dev/chainguard/static:latest

ADD healthchecker /usr/local/bin/healthchecker

ENTRYPOINT ["healthchecker"]