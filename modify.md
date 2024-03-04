# 修改记录

加了一堆daclient相关的rpc配置，用于访问da层的rpc。修改就集中在以下几块：

1. op-node：rollup consensus-layer client。相当于从l1区块中存放的l2数据来重新导出l2链，可以用执行引擎（evm相关）来执行交易，计算l2 block，用于验证和读取。
2. op-batcher：打包l2交易到l1的角色。会与op-service里的txmgr交互。
3. op-memo：与其client相关的配置参数与命令行参数。会调用外部包`github.com/rollkit/go-da/proxy`，client具体提供的接口在引的包里有封装。

## DA

在 `op-memo` 里修改client的部分配置，并就地封装Memo中间件的接口。重点为 `get` 与 `submit` 功能，另外还需要根据配置完成中间件平台账户的初始化操作，保证后续功能的正常使用。

## op-node

嵌入点为 `calldata_source.go` 的 `DataFromEVMTransactions`，如果读取到特定前缀的数据，可认为是DA数据，需要调用DA的client完成`get`操作读取出真正的数据。这块如果能保持接口一致则不需要修改，否则需要做少许调整。

## op-batcher

嵌入点为 `driver.go` 的 `sendTransaction`。在发L2交易时打包前会将交易通过DA的client完成`submit`上传操作，并将返回的值作为交易的txdata封装成L1交易。这块如果能保持接口一致则不需要修改，否则需要做少许调整。

## scripts

1. `bedrock-devnet/devnet/__init__.py`里的各个启动参数。
2. `ops-bedrock/docker-compose.yml`用于启动容器的参数。

## 账号记录

The batch submitter:

- Address: `0xde3829a23df1479438622a08a116e8eb3f620bb5`
- Private key: `bf7604d9d3a1c7748642b1b7b05c2bd219c9faa91458b370f85e5a40f3b03af7`

The devnet comes with a pre-funded account you can use as a faucet:

- Address: `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`
- Private key: `ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
