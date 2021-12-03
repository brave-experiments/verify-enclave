attest-enclave
==============

This tool attests a remotely running AWS Nitro enclave, i.e., it ensures that
the remotely running code is identical to a given local code repository.

Usage
-----

To attest the enclave, run the following:

    make verify CODE=/path/to/ia2/
