# sol

sol is curl for calling eth contracts.

Quick & dirty prototype that:

* Looks up Contract ABI from etherscan.io
* ABI encodes JSON array into calldata
* ABI decodes returndata

# TODO

* hmm. maybe better to treat arguments as individual cli arguments rather than JSON
* Cache contract ABI somehow
* Listen to log events
* Be able to send transactions?

# Install

```
go get github.com/hayeah/sol
```

# Examples

* USDT address: [0xdac17f958d2ee523a2206206994597c13d831ec7](https://etherscan.io/token/0xdac17f958d2ee523a2206206994597c13d831ec7)
* AAVE Pool: [0x3dfd23a6c5e8bbcfc9581d2e864a68feb6a076d3](https://etherscan.io/address/0x3dfd23a6c5e8bbcfc9581d2e864a68feb6a076d3)


Look up the USDT balance of AAVE balance pool:

```
sol 0xdac17f958d2ee523a2206206994597c13d831ec7 \
    balanceOf '["0x3dfd23a6c5e8bbcfc9581d2e864a68feb6a076d3"]'
```

```
sol 0xdac17f958d2ee523a2206206994597c13d831ec7 \
    decimals
```