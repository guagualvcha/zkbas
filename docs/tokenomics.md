# Tokenomics

## Supported Tokens
ZkBAS can be used to transfer BNB, BEP20 and BEP721 tokens. ZkBAS supports a maximum of 65535 BEP20 tokens and 1099511627775 BEP721 tokens. For BEP20 tokens, each supported token need to be listed on ZkBAS L1 contract first. For BEP721 tokens, there are two sources: one is to deposit BEP721 token from BSC network to ZkBAS, another is to directly mint on ZkBAS.

ZkBAS is not responsible for security of supported token contract. Please use any token on your own risk.

## List Token
ZkBAS use `AssetGovernance` contract to manage supported tokens. To list token on ZkBAS, call function `addAsset(address _assetAddress)` of AssetGovernance contract. 

Notice there is a listing fee for listing token, and fees are denominated using a specific token. The listing fee and fee token can be retrived by calling function `listingFee` and `listingFeeToken`. Make sure the sender that calls `addAsset` has enough fee token balance.

## Fee
In ZkBAS the cost of every transaction comes from two parts:

- **L2 part**: every transaction needs to read or write state storage, and ZkBAS prover needs to generate a SNARK proof for every transaction. This part depends on the use of hardware resources and is therefor invariable;

- **L1 part**: ZkBAS must pay BNB gas to commit, verify and execute L2 blocks by sending corresponding L1 transaction. The L1 fees need to be averaged per L2 transaction which is orders of magnitude cheaper than the cost of normal BNB/BEP20 transfers. In addition, for the special transaction types which need to be executed on contract such as `Withdraw` and `WithdrawNft`, there is extra gas cost to be covered.

Users can easily get fee cost of every transaction type using rpc method provided by ZkBAS, then pay transaction fees in multi fee tokens supported by ZkBAS. For example, suppose ZkBAS supports BNB/USDT, when users make a transaction, users can use BNB or USDT to pay transaction fees for their own convenience.