# Use distroless as minimal base image to package the zupd binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
ARG TARGETPLATFORM
LABEL PROJECT="krcrdr" \
      MAINTAINER="mbengtss@redhat.com" \
      DESCRIPTION="A Kubernetes Recorder" \
      LICENSE="Apache-2.0" \
      PLATFORM="$TARGETPLATFORM" \
      VCS_URL="github.com/cldmnky/krcrdr" \
      COMPONENT="krcrdr"
WORKDIR /
COPY ${TARGETPLATFORM}/krcrdr /krcrdr
ENTRYPOINT ["/krcrdr"]
