package constants

const (
	BID                  = true
	ASK                  = false
	STATUS_LENGTH        = 1
	NUM_PARTS            = 2
	LIGHTNING_ORDER_SIZE = 166 // 16 + 32 + 32 + 1 + 20 + 65
	LIGHTNING_MSG_SIZE   = 150 // 16 + 16 + 32 + 1 + 20 + 65
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
	KEY_ALICE          = "2d7aaa9b78d759813448eb26483284cd5e4344a17dede2ab7f062f0757113a28"
	KEY_BOB            = "0e5c6904f09186a0cfe945da201e9d9f0443e07d9e795a9d26cc5cbb882874ac"
	KEY_SUPER_MATCHER  = "7f60d75be8f8833a47311c001adbc3771784c52ea9115200a516e3f050c3bc2b"
	KEY_REPORTER       = "949dbd0607598c41478b32c27da65ab550d54246922fa8978a8c1b9e901e06a6"
	KEY_WORKER         = "e5faea48461ef5a0b78839573073e5a2f579155bf7a4cceb15e49b41963af6a3"
	SUPER_MATCHER_PORT = 8080
	NO_MATCHER         = 20
	SEND_TO            = 3
)

const (
	NO_MINTED_GVN_TOKEN = 99999999999
	NO_GVN_IN_CHANNEL   = 499999999
	NO_ETH_IN_CHANNEL   = 9999999999
	NO_ETH_DEPOSIT      = 9999999999
	NO_GVN_APPROVE      = 499999999
)

var MATCHER_PRVKEYS = []string{
	"a5d318db535709319d40d0d80a3fbd9f30d7ba3b9efc6ed074f939837d937a92",
	"15f29a7247bd4783f2ca940959c353ec2f1d60f816d698dfe339944ff45c047e",
	"0ef19d6ee5fd8022d2d36c8a7c1eb92e9cdfe69459053ab0629aed3d7c617e47",
	"6f38a43f6771ef6d65c72ab50a10cb29d7c2f94bfdc8c8ff8a0dfb888d006225",
	"245892a5a56effd2f79b4d11e147714dbe8fbfadea43ce36a4899d0e022a6fa6",
	"e47d1a7c32e936fc785b9a618033dbed650a232f97146e23298ce77a02fdb43a",
	"8a0cb7f34f7f6bb5ddf9ddb1cc9c1ca03d6246a24d696f566b4894938ecd2f9e",
	"097d045511862b5d39b6847069e6d5c9c20dfc5f72cfe706fa5b294ddc01cc70",
	"5efe41961b93fc4b525c27736bce9be780d9f402abec023aa27e88699fa4c83b",
	"3001263b69ae16852aaeb82f30e1dffb34616ddbd664508c9fe10ed5946ffcb3",
	"8f9353f716ba28f7b2ed249d17ce5c67b219b34deac7ff8bf0a8afb73fc87ebf",
	"957b5ee488b5cc70d2b06156619fbad48adfb2952ed26a28fcba9055853f1cfa",
	"559aeb4596137352deaffd76ea976d9b6340b4ab63d8c9ee28275032df63d94a",
	"51f0e56a53de42293e34c3ed1268459bbb05fa82a994ff62eecf2cc91fbb2602",
	"12be643ab0c8c1ccd3d253c8fac49098704cecc96bb4c1502ae99e2a8b1cd7f7",
	"cb08c7733dbb182e070e77c952a38abe26e079a30ae71ab3dd9322c30158096e",
	"d49f3c77d79d076432af27ba4a184ec8928fa28ed43e9ecff78b7dfd1945b781",
	"45b14162b70d4cc3e11d27a204bf65f73d2b8cf5ade6244141da55f411a27a9e",
	"ae97c20c99b036ff4aaeed38c446f06fe7973cdf1bd3a700817c9880a231f30e",
	"b831a65628de30da2d706faf434200a184be2b0bc1261f4c33b73f33c7e54442",
	"e95be32a7cc2622e1d5dee879bd904f45b1745f7685998a4fb1b23ddaf671dd4",
	"662243510b6e4d8aed80b12d129499384e88f1a7be439d7b1a569e033f8f0a13",
	"74a72dd348c12b682b20fafa814e57a55079d634af559a59acbd3f4260806755",
	"60463b794425d0c2e1e46323c584ba9d556fe472f58b0099e18a6c730dd67512",
	"3508a2339e2fda4be1bc0780ac01e18ccd60ed07b367062d5e6a1b0ee79cb8ee",
	"9ef7d3ccc4330685c3d53e4ad407dc9eaa2eda2cb9b2756d7fb8b28dd9f22e1f",
	"6a60c1bc596336d78045ba45e74c1941f65b695a4149794bc2a1b792280159b7",
	"a82afbdad113124208533241741feb0c408f69ee80e17e64aeed43638ddf8f84",
	"6ab82dd41bf31674d809fea7bba2a5cecb2173e419c312169fb02b899e5d0127",
	"3521eb4f1c92129769e66d46994cf896dab7c821f01b407f19cb5d274d34bae8",
	"92317bfad11e141e81a17de916bedcb4531b240713562904547dac69a7951952",
	"832db43293e7e6aa6436d54f060a011f47a88de0e13c6007f271fba0df308645",
	"efffe85955ccfa80945e2e6bf8a63ee70f61207fb8e8667bdcfdb5cdc1bd8c69",
	"e5fcc3d3d1c72ae910cb1f095715964260ecec7a3b67d5bfaeb786bd16def4f0",
	"7539420e9b155084992143b31e4dc142e5e2fb3b88945802a77099b569fa2474",
	"4389eee713fed7b10cacbbff2f86cd1d923006a2deca54108a07cc347ce36a11",
	"7a95eda3a07248c18d72047347f70f05b003634fc151cc474547c0cc4acbf312",
	"08af21a96000c175e2ffd08898812af5fce9524b6cf0e4ce3dbc1a99456adba0",
	"4855b10c7616a7b84aa6d6dda290fbe84c9096634b3979d05a90782289fa1d1a",
	"e843254ee085c7ca441aac224f9d13deeaf5dba40f7c42fbd991cf1171e23c63",
	"df7ac83c60edf62dff9526aaab5856340f5f397f3fcdbaf052962e9dea2afe86",
	"7cb67ef63da5625175923200f4b7c0fef51c1d0a3b5917bb61897e4fa2db8850",
	"ac5b4516b29aa4db7e9d6912fb0d1bcc8cc0385966d132d4cd814f2f2e494b9f",
	"994eea215b9c18e1f195af482d32a0ab1d8e6f2342772abe0c2d7debb66dac55",
	"8580509506b8ce5fdfc177e31a5b2dbf1de63c8c47dbcd7d41c3a591562984e6",
	"f7df2caafd38d1b3ed87c7b0f78fed2834be8acbf3f7b6c11dae489228aa6e16",
	"0b7dcf72286ece742c8e59a8d125b2f3e8606433389eefc27024493e0890b2cc",
	"136372a15eb0de4f6fcfd8911e3299cde4ccb11186e81fba1752134590fce4d9",
	"cea341edb2b73a524005b7790e39ae1d1bc9a1ea7ed15aebe01ba6025748af03",
	"9cf48844bf5479b1c02cef3e52857cc72a0907d193dc5519d86aa4c88218f2f1",
	"3929ab27262cd8b5882c05d0b932cfc18655ae961afb92ddabe4567482aad192",
	"bdfc7c91723c3fbe3b8898983a86b3a01ee9fd87567c83c18a7112ac92af003e",
	"eaff8561fa2cdae5756440eed9a846564d2028986350e9e310fac22ff859fbad",
	"49ded061f3a931cdb65440c5d5e78f7805cd52ac190a76d88b9506ec9d62ab02",
	"11c605565a8af533070654f62b30cd85da6a7a8ec589bafce640c3d9a5462d20",
	"fc3399d334e063aa5b73a5306158cbc7511b0f59723c9edb5f9603a652ea7311",
	"dbd02a16ec5707226c8fdb3637981d3d80ea98dd6023f3d423cd7b2d4ea75d09",
	"59e6c471e6467e4188a40a7caf1730a13cc6e9011caaba6dbd6993a1604940b6",
	"ef454d6037b8dce79065d910cd4e5f752fa540395d032ed71dd76694d8d76785",
	"06cfd09f3ae6030949110f430d35ef28736bb3b01801ec21a931673d507d9632",
	"c78e4a913db9798dc349e2fcad4e6bff22575de0dc5a4489997e771d0223a2cb",
	"b4dcabe360b18b1b096d639cba23564ccd0c3b22201a0cfe975dcbd66d11d570",
	"650e73130ba53e7fda7508181993c43a74a138c6cf5de25201eb8a1fc3d3995a",
	"6b65f87c20f192bcbc534853de725b8b75f746984f7626f0bc41d981cf3f62e9",
	"a72b073bd252b30db7883ed829b849989d147d9eb163b7fd6de76fe409ca7002",
	"25c966a9726c6e0180597e5c5e547967f500c967183471a9c3eda4020cb09a5a",
	"2bb6fac40443b75231526ffebcf0ef9aaa28fa7910be6fe23c4e203f6f8f9fdf",
	"41150877920e11fdad1f7376c7cb6155de27263085d5ec32b0bc50925aa67ec1",
	"9b4d672591b1a52fbf5377583e8c1ca08b4fdd637dbb882207b7d8dffb8aa4c3",
	"10e0256aa5acd7cfc67142491a4c3f4ffe4ad136994a4ad68d8520ce8bbf000b",
	"48327d35d9376ae565ea7cf758039ba4d9fe5a351b82ca8a9c527732525a6d5f",
	"5a8c90b04b07f66550d37aeea89a4145e9f28bfd8afbbbcc47491a129c3a8bee",
	"7005f4ae7364992b18b7a1639d1894171ef21dbcb72a1d69de514425416f0f4b",
	"099154bd404c9f21c9c225219b1fa8aa11a591d6ff8530c28b9a92c374d82c77",
	"cbd0be1dfc8d82d7a4555bb6013dfb987009c685c211d9536d40af36053a8d0f",
	"dd261b6d61b60e23657a94a7098f3b18435e38e344b6172cc769e4043f7bd660",
	"9c8a7d1b7e87e191e0fe0dde71502b4d7e0fdac0a42bf7318dc976f2735ed13a",
	"732e5d0633f26f88dfb195da08e0da3a2812fb81efe20fd59a6f7a2e707877e5",
	"45568eb62eb841bce07eddb3262ed777c434e2a2d130ab3d759cca6641eabe85",
	"1b93fbccd139e6bd996ce1eb1329088e6b86238a5483c589c245e0ca30e8d8fa",
	"fb6c80453810cbe8892bdd1b7e6a94d782a662207119d96af398f028ea165b89",
	"51abec6d296766676e6e0a382c81e4bccd1b2f5b4c936ed159dd638e52f16122",
	"a0d1e381d05107c5929bdc0c36b9f811df5a6574d76027ec44337e7fab41e922",
	"9f3dcd069b0a49f7eafe2a43eb30c83d60131b2b3db5dee9b19e6f8031c36862",
	"549818746c9058fab5ccfb0a74c91005e2b6a3316490b579bf90cb05442ef3cf",
	"482381b7e31d2d35419d72351abd2103752ee76ede26b870845fbaebcf390efd",
	"3054fef55399a493c7b00030070b2672c88f4adf75b6c633f56b7f9487374b09",
	"57ae39e85f77f371d84732df8ba5d86d336a3b12fce3d467fc026da9b7687115",
	"91abbb69f5cfa93e936f56f8d04fcca7ee1861cc42378edb6b53bf50ee8f204f",
	"a96afdb875edb1cb36ffb3ebeb0b96da4e1564f597c2768b8dbef7de588e828b",
	"13a1774ad6d20678013feba710c8609e72c7a9484fd84699b640a5cc7dc6ae1d",
	"3f6d06b1797f2dcc1fce5ddab7df1b8b38790a98e1aa350d609f0229da6109c2",
	"098d18d7c2f64e2de4a2f296f1c1a0a5cbf018d3e1dcecced7cb3f4bdc5d5332",
	"ab0224161c3e5207e021eb50c08136e62f9baa33af78b531036c7a1d8c7c2dd0",
	"3df6d6bde9d22933ed6a1f032d12f9c22f4d3b0fe49da1c6ad3e4607b5e68272",
	"dcee34ae68a8b7c4580eef398e08b4632664be3972a9de99f682d6315ca173cc",
	"d61c0c3fdd2f123568d60407a51da943c5034d02e95da60fdd64db56ec0bf961",
	"4b1ca7ac7633fc52b999fdde539634122fc13ef3194e84c7104ae40ec7a28e29",
	"6c3099a69ef03766590868ae8bbb02e2bf264e5f2e5b87fa4f78d371d1d73d44",
	"2d74d57316d48e49894880dd1429e5e07c5ad11c7d5691247daea3cf3fdbd0e8",
	"fca722c2690d20f02b8ffb9b5cfd2d2362d12f5d93e91dc98fbbc38a271c7bee",
	"417b5d5b15a2f3105ed0da915f836fdad4cbbeecafe5a4340434fcb0a96ab3fe",
	"9626dc9e81e93f032d0f27009c6be1eb67161751d53e3c010aa6aac5926ac7aa",
	"7fa67da2567af9985ad8f3423bdd7769874908a2507c7112ad9540574509c58d",
	"6f8617a07b66d47c23d736ea0c8e861409b19489e3ed93e28ca755823de592c9",
	"e9e523e907bab7c9211402d40453b9aae1851f4b82bba0c9e9b741948615e30c",
	"c8fe411c70f7026bba128bb73fa9495e6df77ea772e954e3684edb55ecb8ae38",
	"e5d92934af3aed47bc9004d3d17930f3d1ea294408cfc2b423588897be7644d2",
	"964e400e337eabcea995b2ab4e3c2a9cd75f3d681ab560137528fc22e07a896a",
	"c81dec1f4912fbe0157f4d28239ffbf4aa792f26348ea862eeeeedbbba2bd81f",
	"666bb707209c03edd5990deecf1eaacbbe93a387d1760c475a378447650201bb",
	"be754d62543babd107f2cffec969db4da79ca4f3516a2e47bba2ef7e6f870efb",
	"948422a7fac8773b33fec49d8bc08f3ee4d474d1fb7ca1048b00a3523c580527",
	"24142de39fc1047d1375eae4d397c624edafad2baca368f227068bfefef94257",
	"948e7752a0b5819f2b9ce630ef53ea1f4852af677737d12dcaec5d45803b5d2d",
	"4907b3c6d51ea55322f524fb9025551a2f75a429fa361ef739feb5b67a87704b",
	"09ab70c0437bb456b0c4760c8d4d8a9d1569a1b65d2449225b96eda41c2bf86b",
	"0d03186871fea55ce6d277dd7fadbe90345bade63a13810c83fc4da7b8135cbc",
	"8cd503f089878f62fa7e0e58cc5a828b1a707cb79967c7d717ae7ebe380db208",
	"c66d00974a6bcab50be2aef5e30a87bdfa04c103f31a7b7fac7894456688f895",
	"a4df3da1ca5b9ae9e3cc1a3273a8107eb31454b5fc3f73e2abc8a902104998de",
	"4cf744319195e7479e6c66ec37f4dd46d59d208ba660c4049a5a3f050043e3b5",
	"c14b67849f5b9e2cedf300cd61cd8b7bba02ba7bcd486441a4807044a2e38c8d",
	"0ffc4a5d2eb0428a4440864dd9d9f6728f4078fb7f907535cd8574f33c51d33e",
	"8baa0ee8761cafe031c011e48d1f22417baa21071cd3c88c5edcbaa344689cee",
	"dce53693b5cabda37abb88d953f0d041f7bf004e9e6ad4511313dd75bb84cd77",
	"bb9731fed488393ba0a787e372542fb5f11904a317410719564080ffa6c17cb0",
	"d8f61ebb9f28fdab474c336992b8d01f0ec26027e0195ab14bf901e09c0a75f6",
	"c431e7576f90173f4b01663111859c9cd06e0d7c731ba97eaf884f611e5f39d5",
	"c8ed891c5dcb77bbef6408c18f0a5e99db7963cf253ba64c782c8dce8a759020",
	"b2140145b57865c9576bbf413f982144f37fac56c6782fa109018b2a62614466",
	"605f08931bcc892501d8e93e1cb841ee65937c166db1699eb2cc99f697d74f63",
	"d1a5729c386d1885fb3e0b4af2dbae96a505810520776951b1966cb20203e409",
	"e97f9b0349fdd1c3bb217fe4b663f9c74b8dbbe99011850c519516ad87ee25d1",
	"7271af695db9b7577f8c1b22b94279af2a06946b4354e85e375878f74b8473f9",
	"62426dae4c97c5d1237c2c1548b38f7f242b3af5c93d418eac4446dea840a24f",
	"20c4c89607d3afa689843d6e16ea7c225e9e4a6a3e0d5e317c14b6e9535973d0",
	"0820fccbff0243c9c774f4c96ff9fc47e0fec209ac39887d3732d7aae08182e9",
	"f700f3b556ba0082df55a0c969f7f9ff9798ee8873c2c7ef88b36319cf51b692",
	"137a1d747176ae6cd3ba05f28d867a24627d81a31fd3a9570a7a7fde9ef23bec",
	"956e91f1d3961d4ba0c8513cb101475bc2d795ea9c24ae8eecb15c21f5a39831",
	"5ff7397750004a17afbdad2d08c58729b20a5ef9757a75e6952c514610aa0dea",
	"d449927e674a85a0860ef0f58814c5676f5a7b3f4ca9ef3c01c4108919fa2c35",
	"ac95eefd3e8a82a9e7d68907ac14eac9c250d0e70f1ad81e3ca370810c31a386",
	"f89873ba0d3d020d4752d8d8f01f3ac7931aebb0b963baa6ab534e39d6c6601c",
	"73bfe5934422fc21bce75a40112f3798d7e05c969d63ad3676700578ad3bfc23",
	"95f80f250af28d206d012f2a68ecfae3587ea88fbcf461da39a47e168f29da82",
	"ff3741ed89d08d34966978c26b7d1b8a5c44e0499287b31a4ee1156859dc9853",
	"c454546dc7dbf7abf0d0f6de1ce28ce079355c50a486fb836f5e2a221b33dd27",
	"ac1f362ced6be6af1514957eb1daf3dd93d288240382acedd669a8fda1255e40",
}
