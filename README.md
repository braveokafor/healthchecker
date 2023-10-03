# Health Checker

[![Build Status][badge_build_status]][link_build_status]
[![Release Status][badge_release_status]][link_build_status]
[![Repo size][badge_repo_size]][link_repo]
[![Image size][badge_size_latest]][link_docker_hub]
[![Docker Pulls][badge_docker_pulls]][link_docker_hub]

## Overview
Health Checker is a Lightweight Utility for HTTP Health Checks in Distroless Docker Environments.  

## Motivation
The utility provides a `HEALTHCHECK` mechanism in distroless Docker images.  
Distroless images lack a shell, making standard health checks challenging; Health Checker addresses this gap.

## How It Works
Health Checker sends HTTP requests to a specified target URL and assesses the responses against expected status codes set by the user.  
It supports configuration via both command-line flags and environment variables.

## Installation
Download the appropriate binary for your operating system from the [GitHub Releases page](https://github.com/braveokafor/healthchecker/releases).

## Usage  
Run the healthchecker binary with desired flags, or set the environment variables before starting the utility.

## Options

| Flag	             | Environment Variable	     | Description	                         | Default Value    | 
|--------------------|---------------------------|---------------------------------------|------------------|
| `-url`	         | `HC_URL`	                 | Target URL for health check	         | http://localhost | 
| `-status`	         | `HC_EXPECTED_STATUS_CODE` | Expected HTTP status code from target | 200              |
| `-timeout`	     | `HC_TIMEOUT`	             | Request timeout duration	             | 2s               |

## Examples:

### CLI:
```sh
# With flags
healthchecker -url=http://localhost -status=200 -timeout=2s

# With environment variables
export HC_URL=http://localhost
export HC_EXPECTED_STATUS_CODE=200
export HC_TIMEOUT=2s
healthchecker
```

### Docker
```sh
docker run braveokafor/healthchecker:latest -url=http://localhost -status=200 -timeout=2s
```

### Dockerfile:
```Dockerfile
FROM golang:1.21 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=0 go build -o /go/bin/app


FROM cgr.dev/chainguard/static:latest
USER nobody

COPY --from=build /go/bin/app /
COPY --from=braveokafor/healthchecker:latest /usr/local/bin/healthchecker /healthchecker

EXPOSE 5000

# Healthcheck setup
HEALTHCHECK --interval=10s --timeout=2s --start-period=5s --retries=3 \
    CMD ["/healthchecker", "-url", "http://localhost:5000/healthz", "-status", "200", "-timeout", "2s"] || exit 1

CMD ["/app"]
```

## Contributing
Contributions and suggestions are welcome.  
Please open an issue for discussion or propose improvements directly through a pull request.

## Support & Issues

[![Issues][badge_issues]][link_issues]
[![Issues][badge_pulls]][link_pulls]

For support or reporting issues, please open an issue in the GitHub repository or reach out on [LinkedIn](https://www.linkedin.com/in/braveokafor/).

## License
Health Checker is under the MIT [License](https://github.com/braveokafor/healthchecker/blob/main/LICENSE).  
Feel free to use, modify, and distribute the code per the terms of the license.

## Need Assistance?
For questions or further assistance, kindly open an issue.

Thank you for using Health Checker!


[link_issues]:https://github.com/braveokafor/healthchecker/issues
[link_pulls]:https://github.com/braveokafor/healthchecker/pulls
[link_build_status]:https://github.com/braveokafor/healthchecker/actions/workflows/go.yaml
[link_build_status]:https://github.com/braveokafor/healthchecker/actions/workflows/release.yaml
[link_docker_hub]:https://hub.docker.com/r/braveokafor/healthchecker
[link_repo]:https://github.com/braveokafor/healthchecker

[badge_issues]:https://img.shields.io/github/issues-raw/braveokafor/healthchecker?style=flat-square&logo=GitHub
[badge_pulls]:https://img.shields.io/github/issues-pr/braveokafor/healthchecker?style=flat-square&logo=GitHub
[badge_build_status]:https://img.shields.io/github/actions/workflow/status/braveokafor/healthchecker/go.yaml?style=flat-square&logo=GitHub&label=build
[badge_release_status]:https://img.shields.io/github/actions/workflow/status/braveokafor/healthchecker/release.yaml?style=flat-square&logo=GitHub&label=release
[badge_size_latest]:https://img.shields.io/docker/image-size/braveokafor/healthchecker/latest?style=flat-square&logo=Docker
[badge_docker_pulls]:https://img.shields.io/docker/pulls/braveokafor/healthchecker?style=flat-square&logo=Docker
[badge_repo_size]:https://img.shields.io/github/repo-size/braveokafor/healthchecker?style=flat-square&logo=GitHub
