package constants

const (
	BID                  = true
	ASK                  = false
	STATUS_LENGTH        = 1
	NUM_PARTS            = 2
	ORDER_SIZE           = 199 // 16 + 32 + 32 + 1 + 20 + 65 + 1 + 32
	ETH                  = 0
	GVN                  = 1
	TRADER               = 0
	MATCHER              = 1
	TX_FINALITY_DEPTH    = 1
	NO_BATCHES_EACH_TIME = 1
)

const (
	CHAIN_URL          = "ws://127.0.0.1:8545"
	CHAIN_ID           = 1337
	KEY_DEPLOYER       = "abf82ff96b463e9d82b83cb9bb450fe87e6166d4db6d7021d0c71d7e960d5abe"
	KEY_MATCHER        = "dcb7118c9946a39cd40b661e0d368e4afcc3cc48d21aa750d8164ca2e44961c4"
	KEY_ALICE          = "2d7aaa9b78d759813448eb26483284cd5e4344a17dede2ab7f062f0757113a28"
	KEY_BOB            = "0e5c6904f09186a0cfe945da201e9d9f0443e07d9e795a9d26cc5cbb882874ac"
	KEY_SUPER_MATCHER  = "7f60d75be8f8833a47311c001adbc3771784c52ea9115200a516e3f050c3bc2b"
	SUPER_MATCHER_PORT = 8080
)
