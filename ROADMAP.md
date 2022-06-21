# Roadmap

## Wishlist

Support for multiple transports:
* TCP (len+data) via Noise[^1]
* UDP (just data) via Noise
* SCTP[^2]
* IPFS via Noise *(sounds indirect/rendezvous)*
* Email *(sounds indirect/rendezvous)*

General idea is to use Noise for security and multiple transports for connectivity


----
[^1]: see https://noiseprotocol.org/

[^2]: see https://pkg.go.dev/github.com/pion/sctp