# ArgoX Protocol Implementation - Interaction DEMO

## 1. Introduction
`ArgoX Protocol` is a decentralized exchange based on Limit Order Book. With the combination of Lightning network (State channel) and Batching mechanism to help increasing the speed, throughput and reduce the gas usage of the exchange. *This is our final university project, and it has not been published yet.*

This repository is a Go implementation of the `ArgoX Protocol` and publishes some APIs. It allows interaction from the front-end and helps demonstrate how the protocol works, we use [github.com/NguyenHiu/ArgoX-Implementation-FE](https://github.com/NguyenHiu/ArgoX-Implementation-FE) as our front-end in the demo. You can check the following video for the workflow: [youtube.com/watch?v=L6DvWfdBZxs](https://www.youtube.com/watch?v=L6DvWfdBZxs)

## 2. Execution
### 2.1. Run ganache 
```
$ ganache -a 200 -m '' -e 99999999999 --chain.chainId 1337 --p 8545
```


### 2.2. Run protocol
```
$ go run .
```

**Then, you need to start the front-end in the [github.com/NguyenHiu/ArgoX-Implementation-FE](https://github.com/NguyenHiu/ArgoX-Implementation-FE) to interact with this back-end.**

## Note
- Don't forget to restart your ganache before each run
- With each ganache blockchain, you can only run one program. So if you want to run multiple instances of the protocol at a time, you should start multiple ganache blockchains using **different ports and chain IDs**