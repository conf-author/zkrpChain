# zkrpChain
ZkrpChain, which integrates Bulletproofs [1] and Hyperledger Fabric [2] in a weak-coupling way, is a privacy-preserving data auditing solution for consortium blockchains. It allows the verifier peer to conduct privacy-preserving data auditing and verify the zero-knowledge range proofs, which the prover peer generates from its off-chain private data, mainly based on the chaincodes and on-chain data of Hyperledger Fabric. This project is based on Hyperledger Fabric v1.1.0.

About the one-prover-to-one-verifier scenario
----------
The part consists of 5 chaincodes, namely `GenPrfs_StdRng.go`, `GenPrfs_ArbRng.go`, `VerPrfs_StdRng.go`, `VerPrfs_ArbRng.go` and `RangeMgt.go`. The former two are used to generate standard-range and arbitrary-range proofs for prover’s private data, and the following two can be used by verifier peer to verify standard-range and arbitrary-range proofs generated by the prover peer. If the number of input private data item is greater than 1, the proof-generation chaincodes automatically aggregate all the proofs generated from the private data items. The last chaincode RangeMgt is invoked by the former four chaincodes and responsible for uploading/downloading the range values to/from on-chain public ledgers and querying all the update logs of the range values under one given range key ID.

Besides the 5 main chaincodes, client codes are also developed to show how to invoke the chaincodes and related APIs. In this part, 4 client codes are presented, namely `RangeManagement.js`, `Prover.js`, `Verifier.js` and `Service.js`, of which the former three can be used to invoke corresponding chaincode APIs and the latter is the invocation interfaces between client codes and the corresponding chaincodes APIs. 

Additionally, `Prover.js` is also responsible for retrieving the private data from off-chain database, which is MySQL and the related database file is project_VegetablesInfo.sql.

About the multiple-provers-to-one-verifier scenario
----------
The part consists of 6 chaincodes, namely `Prover.go`, `Arb_Prover.go`, `Dealer.go`, `Arb_Dealer.go`, `Verifier.go`,  and `Arb_Verifier.go`. The former two are used to generate `V_A_and_S`, `T1_and_T2`, and `OtherShare` for each prover’s private data. The following two are used to generate `Fait_y_z`, `Fait_x` and `FinalProof` by the dealer peer, and the last two can be used by verifier peer to verify standard-range and arbitrary-range proofs generated by the dealer peer. 

Similarly, in this part, it also has 4 client codes, namely `Prover.js`, `Dealer.js`, `Verifier.js` and `Service.js`. The function of client codes about the multiple-provers-to-one-verifier is the same as that about the one-prover-to-one-verifier. And the usage process of the client codes is as follows:

1.  Setup the client codes. 
```
node Prover.js 
node Dealer.js
node Verifier.js
```

2. Open the browser and input the URL. For example: http://localhost:3000?VerPrf_StandardRange.

3. Enter and start to invoke the corresponding function. In the above URL, the client codes will invoke the `Service.js APIs` to operate `Verifier chaincode` to verify the stardard-range proofs.

References
----------
[1] B. Bnz, J. Bootle, D. Boneh, A. Poelstra, P. Wuille, and G. Maxwell, “Bulletproofs: Short proofs for confidential transactions and more,” in Proc. IEEE S&P, San Francisco, CA, USA, 2018, pp. 315–334.

[2] Hyperledger. Hyperledger Fabric home. Accessed: Jul. 25, 2020. [Online]. Available: https://www.hyperledger.org/use/fabric.
