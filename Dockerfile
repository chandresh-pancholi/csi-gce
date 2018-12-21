# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


FROM golang:1.11.1-stretch as builder
WORKDIR /go/src/github.com/chandresh-pancholi/csi-gce
ADD . .
RUN make


FROM ubuntu:16.04
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates wget autofs \
        && echo "deb http://packages.cloud.google.com/apt cloud-sdk-xenial main" | tee /etc/apt/sources.list.d/google-cloud.sdk.list \
        && echo "deb http://packages.cloud.google.com/apt gcsfuse-xenial main" | tee /etc/apt/sources.list.d/gcsfuse.list \
        && wget -qO- https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - \
        && apt-get update && apt-get install -y --no-install-recommends google-cloud-sdk gcsfuse make \
        && mkdir -p /etc/autofs && touch /etc/autofs/auto.gcsfuse && rm -rf /var/lib/apt/lists




COPY --from=builder /go/src/github.com/chandresh-pancholi/csi-gce/bin/csi-gce oneconcern/csi-gce

#FROM scratch
#RUN apt-get install build-essential ca-certificates e2fsprogs util-linux -y
#COPY --from=builder /go/src/github.com/oneconcern/csi-gce/bin/csi-gce oneconcern/csi-gce

ENTRYPOINT ["/bin/csi-gce"]
