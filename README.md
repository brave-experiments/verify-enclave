attest-enclave
==============

This tool attests a remotely running AWS Nitro enclave, i.e., it ensures that
the remotely running code is identical to a given local code repository.

Installation
------------

The code currently depends on a patched version of the
[nitrite](https://github.com/hf/nitrite/) library.  The file go.mod contains a
directive that tells the compiler to use a local copy of nitrite rather than the
official one:

    replace github.com/hf/nitrite => ../nitrite

The patched version of nitrite is available
[here](https://github.com/NullHypothesis/nitrite/tree/issue-1).

Usage
-----

To attest the enclave, run the following:

    make verify CODE=/path/to/ia2/ ENCLAVE=https://example.com/attest

For attestation to succeed, your version of both [Go](https://go.dev) and
[ko](https://github.com/google/ko) must be identical to the versions that have
been used to compile the remotely running enclave.
