#!/bin/bash

if [ "$1" = "" ]
then
    echo "Missing CODE argument.  Did you run 'make verify CODE=/path/to/ia2'?" >&2
    exit 1
fi
ia2_path="$1"

ia2_image=$(cd "$ia2_path" && ko publish --local . 2>/dev/null)

cat > Dockerfile <<EOF
FROM public.ecr.aws/amazonlinux/amazonlinux:2

# See:
# https://docs.aws.amazon.com/enclaves/latest/user/nitro-enclave-cli-install.html#install-cli
RUN amazon-linux-extras install aws-nitro-enclaves-cli
RUN yum install aws-nitro-enclaves-cli-devel -y
RUN nitro-cli -V

# Now turn the local Docker image into an Enclave Image File (EIF).
CMD ["/bin/bash", "-c", \
     "nitro-cli build-enclave --docker-uri $ia2_image --output-file dummy.eif 2>/dev/null"]
EOF

verify_image=$(docker build --quiet . | cut -d ':' -f 2)
local_pcr0=$(docker run -ti -v /var/run/docker.sock:/var/run/docker.sock "$verify_image" | \
             jq --raw-output ".Measurements.PCR0")

# Request attestation document from the enclave.
remote_pcr0=$(./fetch-attestation 2>/dev/null)

if [ "$local_pcr0" = "$remote_pcr0" ]
then
    echo "Remote image is identical to local image."
else
    echo -e "WARNING: Remote image IS NOT identical to local image!\n"
    echo -e "\tExpected: $local_pcr0"
    echo -e "\tReceived: $remote_pcr0"
fi
