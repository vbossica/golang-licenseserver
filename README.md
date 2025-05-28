# Golang License Server

This repository contains a simple Golang license server that generates license keys for embedded systems. The license keys are signed with a private key and verified with a public key.

For the purpose of validating the workflow, the server and client code have been integrated into command line applications. For production use, the server code should be run as a service, while the client code is integrated into the embedded system.

1. Build the server and client command line applications:

    ```bash
    make build
    ```

1. Generate the public and private keys in the PEM format:

    ```bash
    ./license-server -generateKeys \
        -privateKeyPath private.pem \
        -publicKeyPath public.pem
    ```

    The private key is stored privately, while the public key is installed on the embedded system, alongside the license key generated in the next step.

2. Generate a license key by providing the name of the feature to support as well as the duration (in months):

    ```bash
    ./license-server -generateLicense \
        -privateKeyPath private.pem \
        -features "feature-1,feature-2" \
        -duration 12 \
        -licensePath license.txt
    ```

3. At startup, the embedded system verifies that the feature corresponding to that system is listed in the license:

    ```bash
    ./license-client -verifyFeature \
        -publicKeyPath public.pem \
        -licensePath license.txt \
        -feature feature-1
    ```
