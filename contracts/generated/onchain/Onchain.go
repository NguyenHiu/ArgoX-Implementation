// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package onchain

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// OnchainMetaData contains all meta data concerning the Onchain contract.
var OnchainMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"name\":\"AcceptBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"name\":\"FullfilMatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"name\":\"PartialMatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"PunishMatcher\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"name\":\"ReceivedBatchDetails\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"name\":\"RemoveBatchOutOfDate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"name\":\"RevertBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"name\":\"WrongOrders\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"batchID\",\"type\":\"bytes16\"}],\"name\":\"deleteBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRegisterFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWaitingTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_m\",\"type\":\"address\"}],\"name\":\"isMatcher\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"batchID\",\"type\":\"bytes16\"}],\"name\":\"isPending\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_m\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"batchID\",\"type\":\"bytes16\"}],\"name\":\"reportMissingDeadline\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"batchID\",\"type\":\"bytes16\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"side\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"sign\",\"type\":\"bytes\"}],\"name\":\"sendBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"batchID\",\"type\":\"bytes16\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"updateBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b50670de0b6b3a76400006002819055505f600381905550612bfa806100335f395ff3fe608060405260043610610085575f3560e01c80634420e486116100585780634420e486146101535780634dcbd09b1461016f57806378b32cf51461019957806382736cd7146101c1578063b29e6299146101e957610085565b806319f5e9fe146100895780633005d34c146100c557806332a58e79146100ed57806336ee674914610117575b5f80fd5b348015610094575f80fd5b506100af60048036038101906100aa9190612075565b610211565b6040516100bc91906120ba565b60405180910390f35b3480156100d0575f80fd5b506100eb60048036038101906100e691906122c6565b6102a9565b005b3480156100f8575f80fd5b50610101610353565b60405161010e919061237a565b60405180910390f35b348015610122575f80fd5b5061013d60048036038101906101389190612393565b61035c565b60405161014a91906120ba565b60405180910390f35b61016d60048036038101906101689190612393565b6103a5565b005b34801561017a575f80fd5b50610183610431565b604051610190919061237a565b60405180910390f35b3480156101a4575f80fd5b506101bf60048036038101906101ba9190612075565b61043a565b005b3480156101cc575f80fd5b506101e760048036038101906101e291906123be565b610733565b005b3480156101f4575f80fd5b5061020f600480360381019061020a9190612075565b610973565b005b5f8060055f846fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f2060050154141580156102a257504260035460055f856fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f20600501546102a09190612429565b105b9050919050565b7f48ecfbdd39de4d68f2a28458698ce634ea21a4083a608e845c14122c9d70bc0d866040516102d8919061246b565b60405180910390a15f6040518060e00160405280886fffffffffffffffffffffffffffffffff1916815260200187815260200186815260200185151581526020018473ffffffffffffffffffffffffffffffffffffffff1681526020018381526020015f815250905061034a81610c2a565b50505050505050565b5f600254905090565b5f8060045f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205414159050919050565b6002543410156103ea576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103e1906124de565b60405180910390fd5b60025460045f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f208190555050565b5f600354905090565b805f60055f836fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f2060050154036104b7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104ae90612546565b60405180910390fd5b4260035460055f856fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f20600501546105009190612429565b10610540576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610537906125ae565b60405180910390fd5b5f60045f60055f866fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f2060030160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055505f60055f846fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f206005018190555061066760065f846fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f205f9054906101000a900460801b610c5c565b7fef2d1181ef6c5750f7ef1076cc112a454f5bf01f65b8e40daeb7390d2b66022082604051610696919061246b565b60405180910390a17f5d03dcef971a6d5b97413cad12abae79f43e9422a6c38e8bc70592b18937ba2360055f846fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f2060030160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1660405161072791906125db565b60405180910390a15050565b813373ffffffffffffffffffffffffffffffffffffffff1660055f836fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f2060030160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146107fc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107f39061263e565b60405180910390fd5b825f73ffffffffffffffffffffffffffffffffffffffff1660055f836fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f2060030160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036108c5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108bc906126a6565b60405180910390fd5b60055f856fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f206003015f9054906101000a900460ff16156109215761091c5f8585610e80565b61092e565b61092d60018585610e80565b5b8260055f866fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f206002018190555050505050565b803373ffffffffffffffffffffffffffffffffffffffff1660055f836fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f2060030160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610a3c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a339061263e565b60405180910390fd5b815f73ffffffffffffffffffffffffffffffffffffffff1660055f836fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f2060030160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603610b05576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610afc906126a6565b60405180910390fd5b60055f846fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f206003015f9054906101000a900460ff1615610b6057610b5b5f84610f27565b610b6c565b610b6b600184610f27565b5b60055f846fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f205f8082015f6101000a8154906fffffffffffffffffffffffffffffffff0219169055600182015f9055600282015f9055600382015f6101000a81549060ff02191690556003820160016101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600482015f610c1c9190611fb7565b600582015f90555050505050565b806060015115610c4457610c3e5f826111b8565b50610c51565b610c4f6001826111b8565b505b610c59611772565b50565b5f60055f836fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f206005015403610e7d577ff335387fffae8d4c523007d8c373bffc677e94ca97652a6dabecea00f48e4dd681604051610ccc919061246b565b60405180910390a15f60055f836fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f206040518060e00160405290815f82015f9054906101000a900460801b6fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020016001820154815260200160028201548152602001600382015f9054906101000a900460ff161515151581526020016003820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600482018054610de6906126f1565b80601f0160208091040260200160405190810160405280929190818152602001828054610e12906126f1565b8015610e5d5780601f10610e3457610100808354040283529160200191610e5d565b820191905f5260205f20905b815481529060010190602001808311610e4057829003601f168201915b505050505081526020016005820154815250509050610e7b81610c2a565b505b50565b5f5b8380549050811015610f2157826fffffffffffffffffffffffffffffffff1916848281548110610eb557610eb4612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b6fffffffffffffffffffffffffffffffff191603610f145781848281548110610efc57610efb612721565b5b905f5260205f20906006020160020181905550610f21565b8080600101915050610e82565b50505050565b5f805b60018480549050610f3b919061274e565b81101561110b5781158015610fa25750826fffffffffffffffffffffffffffffffff1916848281548110610f7257610f71612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b6fffffffffffffffffffffffffffffffff1916145b15610fac57600191505b81156110fe5783600182610fc09190612429565b81548110610fd157610fd0612721565b5b905f5260205f209060060201848281548110610ff057610fef612721565b5b905f5260205f2090600602015f82015f9054906101000a900460801b815f015f6101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506001820154816001015560028201548160020155600382015f9054906101000a900460ff16816003015f6101000a81548160ff0219169083151502179055506003820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff168160030160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600482018160040190816110ef9190612933565b50600582015481600501559050505b8080600101915050610f2a565b508280548061111d5761111c612a18565b5b600190038181905f5260205f2090600602015f8082015f6101000a8154906fffffffffffffffffffffffffffffffff0219169055600182015f9055600282015f9055600382015f6101000a81549060ff02191690556003820160016101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600482015f6111a89190611fb7565b600582015f905550509055505050565b5f807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff90508260600151156111ec57600190505b5f8480549050148061124a575080846001868054905061120c919061274e565b8154811061121d5761121c612721565b5b905f5260205f209060060201600101546112379190612a4e565b8184602001516112479190612a4e565b13155b15611345578383908060018154018082558091505060019003905f5260205f2090600602015f909190919091505f820151815f015f6101000a8154816fffffffffffffffffffffffffffffffff021916908360801c021790555060208201518160010155604082015181600201556060820151816003015f6101000a81548160ff02191690831515021790555060808201518160030160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060a08201518160040190816113339190612ace565b5060c082015181600501555050611767565b5f5b8480549050811015611765578185828154811061136757611366612721565b5b905f5260205f209060060201600101546113819190612a4e565b8285602001516113919190612a4e565b1315611758578485600187805490506113aa919061274e565b815481106113bb576113ba612721565b5b905f5260205f209060060201908060018154018082558091505060019003905f5260205f2090600602015f909190919091505f82015f9054906101000a900460801b815f015f6101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506001820154816001015560028201548160020155600382015f9054906101000a900460ff16816003015f6101000a81548160ff0219169083151502179055506003820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff168160030160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600482018160040190816114e09190612933565b506005820154816005015550505f600286805490506114ff919061274e565b90505b818111156116685785600182611518919061274e565b8154811061152957611528612721565b5b905f5260205f20906006020186828154811061154857611547612721565b5b905f5260205f2090600602015f82015f9054906101000a900460801b815f015f6101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506001820154816001015560028201548160020155600382015f9054906101000a900460ff16816003015f6101000a81548160ff0219169083151502179055506003820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff168160030160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600482018160040190816116479190612933565b5060058201548160050155905050808061166090612b9d565b915050611502565b508385828154811061167d5761167c612721565b5b905f5260205f2090600602015f820151815f015f6101000a8154816fffffffffffffffffffffffffffffffff021916908360801c021790555060208201518160010155604082015181600201556060820151816003015f6101000a81548160ff02191690831515021790555060808201518160030160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060a08201518160040190816117459190612ace565b5060c08201518160050155905050611765565b8080600101915050611347565b505b600191505092915050565b5b5f80805490501415801561178c57505f60018054905014155b80156117dd575060015f815481106117a7576117a6612721565b5b905f5260205f209060060201600101545f80815481106117ca576117c9612721565b5b905f5260205f2090600602016001015410155b15611fb55760015f815481106117f6576117f5612721565b5b905f5260205f209060060201600201545f808154811061181957611818612721565b5b905f5260205f2090600602016002015411156118c75760015f8154811061184357611842612721565b5b905f5260205f209060060201600201545f808154811061186657611865612721565b5b905f5260205f2090600602016002015f828254611883919061274e565b925050819055506118c26001805f815481106118a2576118a1612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b610f27565b611fb0565b60015f815481106118db576118da612721565b5b905f5260205f209060060201600201545f80815481106118fe576118fd612721565b5b905f5260205f2090600602016002015410156119ab575f808154811061192757611926612721565b5b905f5260205f2090600602016002015460015f8154811061194b5761194a612721565b5b905f5260205f2090600602016002015f828254611968919061274e565b925050819055506119a65f805f8154811061198657611985612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b610f27565b611faf565b425f80815481106119bf576119be612721565b5b905f5260205f209060060201600501819055505f80815481106119e5576119e4612721565b5b905f5260205f20906006020160055f805f81548110611a0757611a06612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b6fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f205f82015f9054906101000a900460801b815f015f6101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506001820154816001015560028201548160020155600382015f9054906101000a900460ff16816003015f6101000a81548160ff0219169083151502179055506003820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff168160030160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060048201816004019081611b489190612933565b50600582015481600501559050504260015f81548110611b6b57611b6a612721565b5b905f5260205f2090600602016005018190555060015f81548110611b9257611b91612721565b5b905f5260205f20906006020160055f60015f81548110611bb557611bb4612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b6fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f205f82015f9054906101000a900460801b815f015f6101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055506001820154816001015560028201548160020155600382015f9054906101000a900460ff16816003015f6101000a81548160ff0219169083151502179055506003820160019054906101000a900473ffffffffffffffffffffffffffffffffffffffff168160030160016101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060048201816004019081611cf69190612933565b506005820154816005015590505060015f81548110611d1857611d17612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b60065f805f81548110611d4957611d48612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b6fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f205f6101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055505f8081548110611dd057611dcf612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b60065f60015f81548110611e0257611e01612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b6fffffffffffffffffffffffffffffffff19166fffffffffffffffffffffffffffffffff191681526020019081526020015f205f6101000a8154816fffffffffffffffffffffffffffffffff021916908360801c02179055507f57c1c352bdac3386003c09fd995913d2076b2fcd7da4e5b8c98e19c1b08ae9655f8081548110611eaa57611ea9612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b604051611ed2919061246b565b60405180910390a17f57c1c352bdac3386003c09fd995913d2076b2fcd7da4e5b8c98e19c1b08ae96560015f81548110611f0f57611f0e612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b604051611f37919061246b565b60405180910390a1611f776001805f81548110611f5757611f56612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b610f27565b611fae5f805f81548110611f8e57611f8d612721565b5b905f5260205f2090600602015f015f9054906101000a900460801b610f27565b5b5b611773565b565b508054611fc3906126f1565b5f825580601f10611fd45750611ff1565b601f0160209004905f5260205f2090810190611ff09190611ff4565b5b50565b5b8082111561200b575f815f905550600101611ff5565b5090565b5f604051905090565b5f80fd5b5f80fd5b5f7fffffffffffffffffffffffffffffffff0000000000000000000000000000000082169050919050565b61205481612020565b811461205e575f80fd5b50565b5f8135905061206f8161204b565b92915050565b5f6020828403121561208a57612089612018565b5b5f61209784828501612061565b91505092915050565b5f8115159050919050565b6120b4816120a0565b82525050565b5f6020820190506120cd5f8301846120ab565b92915050565b5f819050919050565b6120e5816120d3565b81146120ef575f80fd5b50565b5f81359050612100816120dc565b92915050565b61210f816120a0565b8114612119575f80fd5b50565b5f8135905061212a81612106565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61215982612130565b9050919050565b6121698161214f565b8114612173575f80fd5b50565b5f8135905061218481612160565b92915050565b5f80fd5b5f80fd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6121d882612192565b810181811067ffffffffffffffff821117156121f7576121f66121a2565b5b80604052505050565b5f61220961200f565b905061221582826121cf565b919050565b5f67ffffffffffffffff821115612234576122336121a2565b5b61223d82612192565b9050602081019050919050565b828183375f83830152505050565b5f61226a6122658461221a565b612200565b9050828152602081018484840111156122865761228561218e565b5b61229184828561224a565b509392505050565b5f82601f8301126122ad576122ac61218a565b5b81356122bd848260208601612258565b91505092915050565b5f805f805f8060c087890312156122e0576122df612018565b5b5f6122ed89828a01612061565b96505060206122fe89828a016120f2565b955050604061230f89828a016120f2565b945050606061232089828a0161211c565b935050608061233189828a01612176565b92505060a087013567ffffffffffffffff8111156123525761235161201c565b5b61235e89828a01612299565b9150509295509295509295565b612374816120d3565b82525050565b5f60208201905061238d5f83018461236b565b92915050565b5f602082840312156123a8576123a7612018565b5b5f6123b584828501612176565b91505092915050565b5f80604083850312156123d4576123d3612018565b5b5f6123e185828601612061565b92505060206123f2858286016120f2565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f612433826120d3565b915061243e836120d3565b9250828201905080821115612456576124556123fc565b5b92915050565b61246581612020565b82525050565b5f60208201905061247e5f83018461245c565b92915050565b5f82825260208201905092915050565b7f726567697374657220666565206973206e6f7420656e6f7567680000000000005f82015250565b5f6124c8601a83612484565b91506124d382612494565b602082019050919050565b5f6020820190508181035f8301526124f5816124bc565b9050919050565b7f746865206261746368206973206e6f742070656e64696e6700000000000000005f82015250565b5f612530601883612484565b915061253b826124fc565b602082019050919050565b5f6020820190508181035f83015261255d81612524565b9050919050565b7f746865206261746368206973206e6f74206f75742d6f662d64617465000000005f82015250565b5f612598601c83612484565b91506125a382612564565b602082019050919050565b5f6020820190508181035f8301526125c58161258c565b9050919050565b6125d58161214f565b82525050565b5f6020820190506125ee5f8301846125cc565b92915050565b7f726571756972652062617463682773206f0000000000000000000000000000005f82015250565b5f612628601183612484565b9150612633826125f4565b602082019050919050565b5f6020820190508181035f8301526126558161261c565b9050919050565b7f74686520626174636820646f6573206e6f7420657869737400000000000000005f82015250565b5f612690601883612484565b915061269b8261265c565b602082019050919050565b5f6020820190508181035f8301526126bd81612684565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061270857607f821691505b60208210810361271b5761271a6126c4565b5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f612758826120d3565b9150612763836120d3565b925082820390508181111561277b5761277a6123fc565b5b92915050565b5f8154905061278f816126f1565b9050919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026127f27fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826127b7565b6127fc86836127b7565b95508019841693508086168417925050509392505050565b5f819050919050565b5f61283761283261282d846120d3565b612814565b6120d3565b9050919050565b5f819050919050565b6128508361281d565b61286461285c8261283e565b8484546127c3565b825550505050565b5f90565b61287861286c565b612883818484612847565b505050565b5b818110156128a65761289b5f82612870565b600181019050612889565b5050565b601f8211156128eb576128bc81612796565b6128c5846127a8565b810160208510156128d4578190505b6128e86128e0856127a8565b830182612888565b50505b505050565b5f82821c905092915050565b5f61290b5f19846008026128f0565b1980831691505092915050565b5f61292383836128fc565b9150826002028217905092915050565b818103612941575050612a16565b61294a82612781565b67ffffffffffffffff811115612963576129626121a2565b5b61296d82546126f1565b6129788282856128aa565b5f601f8311600181146129a5575f8415612993578287015490505b61299d8582612918565b865550612a0f565b601f1984166129b387612796565b96506129be86612796565b5f5b828110156129e5578489015482556001820191506001850194506020810190506129c0565b86831015612a0257848901546129fe601f8916826128fc565b8355505b6001600288020188555050505b5050505050505b565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603160045260245ffd5b5f819050919050565b5f612a5882612a45565b9150612a6383612a45565b9250828202612a7181612a45565b91507f800000000000000000000000000000000000000000000000000000000000000084145f84121615612aa857612aa76123fc565b5b8282058414831517612abd57612abc6123fc565b5b5092915050565b5f81519050919050565b612ad782612ac4565b67ffffffffffffffff811115612af057612aef6121a2565b5b612afa82546126f1565b612b058282856128aa565b5f60209050601f831160018114612b36575f8415612b24578287015190505b612b2e8582612918565b865550612b95565b601f198416612b4486612796565b5f5b82811015612b6b57848901518255600182019150602085019450602081019050612b46565b86831015612b885784890151612b84601f8916826128fc565b8355505b6001600288020188555050505b505050505050565b5f612ba7826120d3565b91505f8203612bb957612bb86123fc565b5b60018203905091905056fea2646970667358221220f885a9bdad160d34ae4f34f4a0d642affb3059edf3fdd4dd1ee65b2216cf2db764736f6c63430008180033",
}

// OnchainABI is the input ABI used to generate the binding from.
// Deprecated: Use OnchainMetaData.ABI instead.
var OnchainABI = OnchainMetaData.ABI

// OnchainBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OnchainMetaData.Bin instead.
var OnchainBin = OnchainMetaData.Bin

// DeployOnchain deploys a new Ethereum contract, binding an instance of Onchain to it.
func DeployOnchain(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Onchain, error) {
	parsed, err := OnchainMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OnchainBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Onchain{OnchainCaller: OnchainCaller{contract: contract}, OnchainTransactor: OnchainTransactor{contract: contract}, OnchainFilterer: OnchainFilterer{contract: contract}}, nil
}

// Onchain is an auto generated Go binding around an Ethereum contract.
type Onchain struct {
	OnchainCaller     // Read-only binding to the contract
	OnchainTransactor // Write-only binding to the contract
	OnchainFilterer   // Log filterer for contract events
}

// OnchainCaller is an auto generated read-only Go binding around an Ethereum contract.
type OnchainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnchainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OnchainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnchainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OnchainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OnchainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OnchainSession struct {
	Contract     *Onchain          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OnchainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OnchainCallerSession struct {
	Contract *OnchainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// OnchainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OnchainTransactorSession struct {
	Contract     *OnchainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OnchainRaw is an auto generated low-level Go binding around an Ethereum contract.
type OnchainRaw struct {
	Contract *Onchain // Generic contract binding to access the raw methods on
}

// OnchainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OnchainCallerRaw struct {
	Contract *OnchainCaller // Generic read-only contract binding to access the raw methods on
}

// OnchainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OnchainTransactorRaw struct {
	Contract *OnchainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOnchain creates a new instance of Onchain, bound to a specific deployed contract.
func NewOnchain(address common.Address, backend bind.ContractBackend) (*Onchain, error) {
	contract, err := bindOnchain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Onchain{OnchainCaller: OnchainCaller{contract: contract}, OnchainTransactor: OnchainTransactor{contract: contract}, OnchainFilterer: OnchainFilterer{contract: contract}}, nil
}

// NewOnchainCaller creates a new read-only instance of Onchain, bound to a specific deployed contract.
func NewOnchainCaller(address common.Address, caller bind.ContractCaller) (*OnchainCaller, error) {
	contract, err := bindOnchain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OnchainCaller{contract: contract}, nil
}

// NewOnchainTransactor creates a new write-only instance of Onchain, bound to a specific deployed contract.
func NewOnchainTransactor(address common.Address, transactor bind.ContractTransactor) (*OnchainTransactor, error) {
	contract, err := bindOnchain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OnchainTransactor{contract: contract}, nil
}

// NewOnchainFilterer creates a new log filterer instance of Onchain, bound to a specific deployed contract.
func NewOnchainFilterer(address common.Address, filterer bind.ContractFilterer) (*OnchainFilterer, error) {
	contract, err := bindOnchain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OnchainFilterer{contract: contract}, nil
}

// bindOnchain binds a generic wrapper to an already deployed contract.
func bindOnchain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OnchainMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Onchain *OnchainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Onchain.Contract.OnchainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Onchain *OnchainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Onchain.Contract.OnchainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Onchain *OnchainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Onchain.Contract.OnchainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Onchain *OnchainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Onchain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Onchain *OnchainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Onchain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Onchain *OnchainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Onchain.Contract.contract.Transact(opts, method, params...)
}

// GetRegisterFee is a free data retrieval call binding the contract method 0x32a58e79.
//
// Solidity: function getRegisterFee() view returns(uint256)
func (_Onchain *OnchainCaller) GetRegisterFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Onchain.contract.Call(opts, &out, "getRegisterFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRegisterFee is a free data retrieval call binding the contract method 0x32a58e79.
//
// Solidity: function getRegisterFee() view returns(uint256)
func (_Onchain *OnchainSession) GetRegisterFee() (*big.Int, error) {
	return _Onchain.Contract.GetRegisterFee(&_Onchain.CallOpts)
}

// GetRegisterFee is a free data retrieval call binding the contract method 0x32a58e79.
//
// Solidity: function getRegisterFee() view returns(uint256)
func (_Onchain *OnchainCallerSession) GetRegisterFee() (*big.Int, error) {
	return _Onchain.Contract.GetRegisterFee(&_Onchain.CallOpts)
}

// GetWaitingTime is a free data retrieval call binding the contract method 0x4dcbd09b.
//
// Solidity: function getWaitingTime() view returns(uint256)
func (_Onchain *OnchainCaller) GetWaitingTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Onchain.contract.Call(opts, &out, "getWaitingTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetWaitingTime is a free data retrieval call binding the contract method 0x4dcbd09b.
//
// Solidity: function getWaitingTime() view returns(uint256)
func (_Onchain *OnchainSession) GetWaitingTime() (*big.Int, error) {
	return _Onchain.Contract.GetWaitingTime(&_Onchain.CallOpts)
}

// GetWaitingTime is a free data retrieval call binding the contract method 0x4dcbd09b.
//
// Solidity: function getWaitingTime() view returns(uint256)
func (_Onchain *OnchainCallerSession) GetWaitingTime() (*big.Int, error) {
	return _Onchain.Contract.GetWaitingTime(&_Onchain.CallOpts)
}

// IsMatcher is a free data retrieval call binding the contract method 0x36ee6749.
//
// Solidity: function isMatcher(address _m) view returns(bool)
func (_Onchain *OnchainCaller) IsMatcher(opts *bind.CallOpts, _m common.Address) (bool, error) {
	var out []interface{}
	err := _Onchain.contract.Call(opts, &out, "isMatcher", _m)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsMatcher is a free data retrieval call binding the contract method 0x36ee6749.
//
// Solidity: function isMatcher(address _m) view returns(bool)
func (_Onchain *OnchainSession) IsMatcher(_m common.Address) (bool, error) {
	return _Onchain.Contract.IsMatcher(&_Onchain.CallOpts, _m)
}

// IsMatcher is a free data retrieval call binding the contract method 0x36ee6749.
//
// Solidity: function isMatcher(address _m) view returns(bool)
func (_Onchain *OnchainCallerSession) IsMatcher(_m common.Address) (bool, error) {
	return _Onchain.Contract.IsMatcher(&_Onchain.CallOpts, _m)
}

// IsPending is a free data retrieval call binding the contract method 0x19f5e9fe.
//
// Solidity: function isPending(bytes16 batchID) view returns(bool)
func (_Onchain *OnchainCaller) IsPending(opts *bind.CallOpts, batchID [16]byte) (bool, error) {
	var out []interface{}
	err := _Onchain.contract.Call(opts, &out, "isPending", batchID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPending is a free data retrieval call binding the contract method 0x19f5e9fe.
//
// Solidity: function isPending(bytes16 batchID) view returns(bool)
func (_Onchain *OnchainSession) IsPending(batchID [16]byte) (bool, error) {
	return _Onchain.Contract.IsPending(&_Onchain.CallOpts, batchID)
}

// IsPending is a free data retrieval call binding the contract method 0x19f5e9fe.
//
// Solidity: function isPending(bytes16 batchID) view returns(bool)
func (_Onchain *OnchainCallerSession) IsPending(batchID [16]byte) (bool, error) {
	return _Onchain.Contract.IsPending(&_Onchain.CallOpts, batchID)
}

// DeleteBatch is a paid mutator transaction binding the contract method 0xb29e6299.
//
// Solidity: function deleteBatch(bytes16 batchID) returns()
func (_Onchain *OnchainTransactor) DeleteBatch(opts *bind.TransactOpts, batchID [16]byte) (*types.Transaction, error) {
	return _Onchain.contract.Transact(opts, "deleteBatch", batchID)
}

// DeleteBatch is a paid mutator transaction binding the contract method 0xb29e6299.
//
// Solidity: function deleteBatch(bytes16 batchID) returns()
func (_Onchain *OnchainSession) DeleteBatch(batchID [16]byte) (*types.Transaction, error) {
	return _Onchain.Contract.DeleteBatch(&_Onchain.TransactOpts, batchID)
}

// DeleteBatch is a paid mutator transaction binding the contract method 0xb29e6299.
//
// Solidity: function deleteBatch(bytes16 batchID) returns()
func (_Onchain *OnchainTransactorSession) DeleteBatch(batchID [16]byte) (*types.Transaction, error) {
	return _Onchain.Contract.DeleteBatch(&_Onchain.TransactOpts, batchID)
}

// Register is a paid mutator transaction binding the contract method 0x4420e486.
//
// Solidity: function register(address _m) payable returns()
func (_Onchain *OnchainTransactor) Register(opts *bind.TransactOpts, _m common.Address) (*types.Transaction, error) {
	return _Onchain.contract.Transact(opts, "register", _m)
}

// Register is a paid mutator transaction binding the contract method 0x4420e486.
//
// Solidity: function register(address _m) payable returns()
func (_Onchain *OnchainSession) Register(_m common.Address) (*types.Transaction, error) {
	return _Onchain.Contract.Register(&_Onchain.TransactOpts, _m)
}

// Register is a paid mutator transaction binding the contract method 0x4420e486.
//
// Solidity: function register(address _m) payable returns()
func (_Onchain *OnchainTransactorSession) Register(_m common.Address) (*types.Transaction, error) {
	return _Onchain.Contract.Register(&_Onchain.TransactOpts, _m)
}

// ReportMissingDeadline is a paid mutator transaction binding the contract method 0x78b32cf5.
//
// Solidity: function reportMissingDeadline(bytes16 batchID) returns()
func (_Onchain *OnchainTransactor) ReportMissingDeadline(opts *bind.TransactOpts, batchID [16]byte) (*types.Transaction, error) {
	return _Onchain.contract.Transact(opts, "reportMissingDeadline", batchID)
}

// ReportMissingDeadline is a paid mutator transaction binding the contract method 0x78b32cf5.
//
// Solidity: function reportMissingDeadline(bytes16 batchID) returns()
func (_Onchain *OnchainSession) ReportMissingDeadline(batchID [16]byte) (*types.Transaction, error) {
	return _Onchain.Contract.ReportMissingDeadline(&_Onchain.TransactOpts, batchID)
}

// ReportMissingDeadline is a paid mutator transaction binding the contract method 0x78b32cf5.
//
// Solidity: function reportMissingDeadline(bytes16 batchID) returns()
func (_Onchain *OnchainTransactorSession) ReportMissingDeadline(batchID [16]byte) (*types.Transaction, error) {
	return _Onchain.Contract.ReportMissingDeadline(&_Onchain.TransactOpts, batchID)
}

// SendBatch is a paid mutator transaction binding the contract method 0x3005d34c.
//
// Solidity: function sendBatch(bytes16 batchID, uint256 price, uint256 amount, bool side, address owner, bytes sign) returns()
func (_Onchain *OnchainTransactor) SendBatch(opts *bind.TransactOpts, batchID [16]byte, price *big.Int, amount *big.Int, side bool, owner common.Address, sign []byte) (*types.Transaction, error) {
	return _Onchain.contract.Transact(opts, "sendBatch", batchID, price, amount, side, owner, sign)
}

// SendBatch is a paid mutator transaction binding the contract method 0x3005d34c.
//
// Solidity: function sendBatch(bytes16 batchID, uint256 price, uint256 amount, bool side, address owner, bytes sign) returns()
func (_Onchain *OnchainSession) SendBatch(batchID [16]byte, price *big.Int, amount *big.Int, side bool, owner common.Address, sign []byte) (*types.Transaction, error) {
	return _Onchain.Contract.SendBatch(&_Onchain.TransactOpts, batchID, price, amount, side, owner, sign)
}

// SendBatch is a paid mutator transaction binding the contract method 0x3005d34c.
//
// Solidity: function sendBatch(bytes16 batchID, uint256 price, uint256 amount, bool side, address owner, bytes sign) returns()
func (_Onchain *OnchainTransactorSession) SendBatch(batchID [16]byte, price *big.Int, amount *big.Int, side bool, owner common.Address, sign []byte) (*types.Transaction, error) {
	return _Onchain.Contract.SendBatch(&_Onchain.TransactOpts, batchID, price, amount, side, owner, sign)
}

// UpdateBatch is a paid mutator transaction binding the contract method 0x82736cd7.
//
// Solidity: function updateBatch(bytes16 batchID, uint256 amount) returns()
func (_Onchain *OnchainTransactor) UpdateBatch(opts *bind.TransactOpts, batchID [16]byte, amount *big.Int) (*types.Transaction, error) {
	return _Onchain.contract.Transact(opts, "updateBatch", batchID, amount)
}

// UpdateBatch is a paid mutator transaction binding the contract method 0x82736cd7.
//
// Solidity: function updateBatch(bytes16 batchID, uint256 amount) returns()
func (_Onchain *OnchainSession) UpdateBatch(batchID [16]byte, amount *big.Int) (*types.Transaction, error) {
	return _Onchain.Contract.UpdateBatch(&_Onchain.TransactOpts, batchID, amount)
}

// UpdateBatch is a paid mutator transaction binding the contract method 0x82736cd7.
//
// Solidity: function updateBatch(bytes16 batchID, uint256 amount) returns()
func (_Onchain *OnchainTransactorSession) UpdateBatch(batchID [16]byte, amount *big.Int) (*types.Transaction, error) {
	return _Onchain.Contract.UpdateBatch(&_Onchain.TransactOpts, batchID, amount)
}

// OnchainAcceptBatchIterator is returned from FilterAcceptBatch and is used to iterate over the raw logs and unpacked data for AcceptBatch events raised by the Onchain contract.
type OnchainAcceptBatchIterator struct {
	Event *OnchainAcceptBatch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OnchainAcceptBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnchainAcceptBatch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OnchainAcceptBatch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OnchainAcceptBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OnchainAcceptBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OnchainAcceptBatch represents a AcceptBatch event raised by the Onchain contract.
type OnchainAcceptBatch struct {
	Arg0 [16]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAcceptBatch is a free log retrieval operation binding the contract event 0x48ecfbdd39de4d68f2a28458698ce634ea21a4083a608e845c14122c9d70bc0d.
//
// Solidity: event AcceptBatch(bytes16 arg0)
func (_Onchain *OnchainFilterer) FilterAcceptBatch(opts *bind.FilterOpts) (*OnchainAcceptBatchIterator, error) {

	logs, sub, err := _Onchain.contract.FilterLogs(opts, "AcceptBatch")
	if err != nil {
		return nil, err
	}
	return &OnchainAcceptBatchIterator{contract: _Onchain.contract, event: "AcceptBatch", logs: logs, sub: sub}, nil
}

// WatchAcceptBatch is a free log subscription operation binding the contract event 0x48ecfbdd39de4d68f2a28458698ce634ea21a4083a608e845c14122c9d70bc0d.
//
// Solidity: event AcceptBatch(bytes16 arg0)
func (_Onchain *OnchainFilterer) WatchAcceptBatch(opts *bind.WatchOpts, sink chan<- *OnchainAcceptBatch) (event.Subscription, error) {

	logs, sub, err := _Onchain.contract.WatchLogs(opts, "AcceptBatch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OnchainAcceptBatch)
				if err := _Onchain.contract.UnpackLog(event, "AcceptBatch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAcceptBatch is a log parse operation binding the contract event 0x48ecfbdd39de4d68f2a28458698ce634ea21a4083a608e845c14122c9d70bc0d.
//
// Solidity: event AcceptBatch(bytes16 arg0)
func (_Onchain *OnchainFilterer) ParseAcceptBatch(log types.Log) (*OnchainAcceptBatch, error) {
	event := new(OnchainAcceptBatch)
	if err := _Onchain.contract.UnpackLog(event, "AcceptBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OnchainFullfilMatchIterator is returned from FilterFullfilMatch and is used to iterate over the raw logs and unpacked data for FullfilMatch events raised by the Onchain contract.
type OnchainFullfilMatchIterator struct {
	Event *OnchainFullfilMatch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OnchainFullfilMatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnchainFullfilMatch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OnchainFullfilMatch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OnchainFullfilMatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OnchainFullfilMatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OnchainFullfilMatch represents a FullfilMatch event raised by the Onchain contract.
type OnchainFullfilMatch struct {
	Arg0 [16]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterFullfilMatch is a free log retrieval operation binding the contract event 0x57c1c352bdac3386003c09fd995913d2076b2fcd7da4e5b8c98e19c1b08ae965.
//
// Solidity: event FullfilMatch(bytes16 arg0)
func (_Onchain *OnchainFilterer) FilterFullfilMatch(opts *bind.FilterOpts) (*OnchainFullfilMatchIterator, error) {

	logs, sub, err := _Onchain.contract.FilterLogs(opts, "FullfilMatch")
	if err != nil {
		return nil, err
	}
	return &OnchainFullfilMatchIterator{contract: _Onchain.contract, event: "FullfilMatch", logs: logs, sub: sub}, nil
}

// WatchFullfilMatch is a free log subscription operation binding the contract event 0x57c1c352bdac3386003c09fd995913d2076b2fcd7da4e5b8c98e19c1b08ae965.
//
// Solidity: event FullfilMatch(bytes16 arg0)
func (_Onchain *OnchainFilterer) WatchFullfilMatch(opts *bind.WatchOpts, sink chan<- *OnchainFullfilMatch) (event.Subscription, error) {

	logs, sub, err := _Onchain.contract.WatchLogs(opts, "FullfilMatch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OnchainFullfilMatch)
				if err := _Onchain.contract.UnpackLog(event, "FullfilMatch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFullfilMatch is a log parse operation binding the contract event 0x57c1c352bdac3386003c09fd995913d2076b2fcd7da4e5b8c98e19c1b08ae965.
//
// Solidity: event FullfilMatch(bytes16 arg0)
func (_Onchain *OnchainFilterer) ParseFullfilMatch(log types.Log) (*OnchainFullfilMatch, error) {
	event := new(OnchainFullfilMatch)
	if err := _Onchain.contract.UnpackLog(event, "FullfilMatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OnchainPartialMatchIterator is returned from FilterPartialMatch and is used to iterate over the raw logs and unpacked data for PartialMatch events raised by the Onchain contract.
type OnchainPartialMatchIterator struct {
	Event *OnchainPartialMatch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OnchainPartialMatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnchainPartialMatch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OnchainPartialMatch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OnchainPartialMatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OnchainPartialMatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OnchainPartialMatch represents a PartialMatch event raised by the Onchain contract.
type OnchainPartialMatch struct {
	Arg0 [16]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterPartialMatch is a free log retrieval operation binding the contract event 0xcb56a4fd10f2bad2015ad7e01fb83de3e6d71a6f46eef88ebe216fd70f25efd4.
//
// Solidity: event PartialMatch(bytes16 arg0)
func (_Onchain *OnchainFilterer) FilterPartialMatch(opts *bind.FilterOpts) (*OnchainPartialMatchIterator, error) {

	logs, sub, err := _Onchain.contract.FilterLogs(opts, "PartialMatch")
	if err != nil {
		return nil, err
	}
	return &OnchainPartialMatchIterator{contract: _Onchain.contract, event: "PartialMatch", logs: logs, sub: sub}, nil
}

// WatchPartialMatch is a free log subscription operation binding the contract event 0xcb56a4fd10f2bad2015ad7e01fb83de3e6d71a6f46eef88ebe216fd70f25efd4.
//
// Solidity: event PartialMatch(bytes16 arg0)
func (_Onchain *OnchainFilterer) WatchPartialMatch(opts *bind.WatchOpts, sink chan<- *OnchainPartialMatch) (event.Subscription, error) {

	logs, sub, err := _Onchain.contract.WatchLogs(opts, "PartialMatch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OnchainPartialMatch)
				if err := _Onchain.contract.UnpackLog(event, "PartialMatch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePartialMatch is a log parse operation binding the contract event 0xcb56a4fd10f2bad2015ad7e01fb83de3e6d71a6f46eef88ebe216fd70f25efd4.
//
// Solidity: event PartialMatch(bytes16 arg0)
func (_Onchain *OnchainFilterer) ParsePartialMatch(log types.Log) (*OnchainPartialMatch, error) {
	event := new(OnchainPartialMatch)
	if err := _Onchain.contract.UnpackLog(event, "PartialMatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OnchainPunishMatcherIterator is returned from FilterPunishMatcher and is used to iterate over the raw logs and unpacked data for PunishMatcher events raised by the Onchain contract.
type OnchainPunishMatcherIterator struct {
	Event *OnchainPunishMatcher // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OnchainPunishMatcherIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnchainPunishMatcher)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OnchainPunishMatcher)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OnchainPunishMatcherIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OnchainPunishMatcherIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OnchainPunishMatcher represents a PunishMatcher event raised by the Onchain contract.
type OnchainPunishMatcher struct {
	Arg0 common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterPunishMatcher is a free log retrieval operation binding the contract event 0x5d03dcef971a6d5b97413cad12abae79f43e9422a6c38e8bc70592b18937ba23.
//
// Solidity: event PunishMatcher(address arg0)
func (_Onchain *OnchainFilterer) FilterPunishMatcher(opts *bind.FilterOpts) (*OnchainPunishMatcherIterator, error) {

	logs, sub, err := _Onchain.contract.FilterLogs(opts, "PunishMatcher")
	if err != nil {
		return nil, err
	}
	return &OnchainPunishMatcherIterator{contract: _Onchain.contract, event: "PunishMatcher", logs: logs, sub: sub}, nil
}

// WatchPunishMatcher is a free log subscription operation binding the contract event 0x5d03dcef971a6d5b97413cad12abae79f43e9422a6c38e8bc70592b18937ba23.
//
// Solidity: event PunishMatcher(address arg0)
func (_Onchain *OnchainFilterer) WatchPunishMatcher(opts *bind.WatchOpts, sink chan<- *OnchainPunishMatcher) (event.Subscription, error) {

	logs, sub, err := _Onchain.contract.WatchLogs(opts, "PunishMatcher")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OnchainPunishMatcher)
				if err := _Onchain.contract.UnpackLog(event, "PunishMatcher", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePunishMatcher is a log parse operation binding the contract event 0x5d03dcef971a6d5b97413cad12abae79f43e9422a6c38e8bc70592b18937ba23.
//
// Solidity: event PunishMatcher(address arg0)
func (_Onchain *OnchainFilterer) ParsePunishMatcher(log types.Log) (*OnchainPunishMatcher, error) {
	event := new(OnchainPunishMatcher)
	if err := _Onchain.contract.UnpackLog(event, "PunishMatcher", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OnchainReceivedBatchDetailsIterator is returned from FilterReceivedBatchDetails and is used to iterate over the raw logs and unpacked data for ReceivedBatchDetails events raised by the Onchain contract.
type OnchainReceivedBatchDetailsIterator struct {
	Event *OnchainReceivedBatchDetails // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OnchainReceivedBatchDetailsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnchainReceivedBatchDetails)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OnchainReceivedBatchDetails)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OnchainReceivedBatchDetailsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OnchainReceivedBatchDetailsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OnchainReceivedBatchDetails represents a ReceivedBatchDetails event raised by the Onchain contract.
type OnchainReceivedBatchDetails struct {
	Arg0 [16]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterReceivedBatchDetails is a free log retrieval operation binding the contract event 0x9e82d75e1d25a2db33e754b50ce9378e4c2e505c68c05de244b398543c0e422e.
//
// Solidity: event ReceivedBatchDetails(bytes16 arg0)
func (_Onchain *OnchainFilterer) FilterReceivedBatchDetails(opts *bind.FilterOpts) (*OnchainReceivedBatchDetailsIterator, error) {

	logs, sub, err := _Onchain.contract.FilterLogs(opts, "ReceivedBatchDetails")
	if err != nil {
		return nil, err
	}
	return &OnchainReceivedBatchDetailsIterator{contract: _Onchain.contract, event: "ReceivedBatchDetails", logs: logs, sub: sub}, nil
}

// WatchReceivedBatchDetails is a free log subscription operation binding the contract event 0x9e82d75e1d25a2db33e754b50ce9378e4c2e505c68c05de244b398543c0e422e.
//
// Solidity: event ReceivedBatchDetails(bytes16 arg0)
func (_Onchain *OnchainFilterer) WatchReceivedBatchDetails(opts *bind.WatchOpts, sink chan<- *OnchainReceivedBatchDetails) (event.Subscription, error) {

	logs, sub, err := _Onchain.contract.WatchLogs(opts, "ReceivedBatchDetails")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OnchainReceivedBatchDetails)
				if err := _Onchain.contract.UnpackLog(event, "ReceivedBatchDetails", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReceivedBatchDetails is a log parse operation binding the contract event 0x9e82d75e1d25a2db33e754b50ce9378e4c2e505c68c05de244b398543c0e422e.
//
// Solidity: event ReceivedBatchDetails(bytes16 arg0)
func (_Onchain *OnchainFilterer) ParseReceivedBatchDetails(log types.Log) (*OnchainReceivedBatchDetails, error) {
	event := new(OnchainReceivedBatchDetails)
	if err := _Onchain.contract.UnpackLog(event, "ReceivedBatchDetails", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OnchainRemoveBatchOutOfDateIterator is returned from FilterRemoveBatchOutOfDate and is used to iterate over the raw logs and unpacked data for RemoveBatchOutOfDate events raised by the Onchain contract.
type OnchainRemoveBatchOutOfDateIterator struct {
	Event *OnchainRemoveBatchOutOfDate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OnchainRemoveBatchOutOfDateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnchainRemoveBatchOutOfDate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OnchainRemoveBatchOutOfDate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OnchainRemoveBatchOutOfDateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OnchainRemoveBatchOutOfDateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OnchainRemoveBatchOutOfDate represents a RemoveBatchOutOfDate event raised by the Onchain contract.
type OnchainRemoveBatchOutOfDate struct {
	Arg0 [16]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRemoveBatchOutOfDate is a free log retrieval operation binding the contract event 0xef2d1181ef6c5750f7ef1076cc112a454f5bf01f65b8e40daeb7390d2b660220.
//
// Solidity: event RemoveBatchOutOfDate(bytes16 arg0)
func (_Onchain *OnchainFilterer) FilterRemoveBatchOutOfDate(opts *bind.FilterOpts) (*OnchainRemoveBatchOutOfDateIterator, error) {

	logs, sub, err := _Onchain.contract.FilterLogs(opts, "RemoveBatchOutOfDate")
	if err != nil {
		return nil, err
	}
	return &OnchainRemoveBatchOutOfDateIterator{contract: _Onchain.contract, event: "RemoveBatchOutOfDate", logs: logs, sub: sub}, nil
}

// WatchRemoveBatchOutOfDate is a free log subscription operation binding the contract event 0xef2d1181ef6c5750f7ef1076cc112a454f5bf01f65b8e40daeb7390d2b660220.
//
// Solidity: event RemoveBatchOutOfDate(bytes16 arg0)
func (_Onchain *OnchainFilterer) WatchRemoveBatchOutOfDate(opts *bind.WatchOpts, sink chan<- *OnchainRemoveBatchOutOfDate) (event.Subscription, error) {

	logs, sub, err := _Onchain.contract.WatchLogs(opts, "RemoveBatchOutOfDate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OnchainRemoveBatchOutOfDate)
				if err := _Onchain.contract.UnpackLog(event, "RemoveBatchOutOfDate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRemoveBatchOutOfDate is a log parse operation binding the contract event 0xef2d1181ef6c5750f7ef1076cc112a454f5bf01f65b8e40daeb7390d2b660220.
//
// Solidity: event RemoveBatchOutOfDate(bytes16 arg0)
func (_Onchain *OnchainFilterer) ParseRemoveBatchOutOfDate(log types.Log) (*OnchainRemoveBatchOutOfDate, error) {
	event := new(OnchainRemoveBatchOutOfDate)
	if err := _Onchain.contract.UnpackLog(event, "RemoveBatchOutOfDate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OnchainRevertBatchIterator is returned from FilterRevertBatch and is used to iterate over the raw logs and unpacked data for RevertBatch events raised by the Onchain contract.
type OnchainRevertBatchIterator struct {
	Event *OnchainRevertBatch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OnchainRevertBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnchainRevertBatch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OnchainRevertBatch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OnchainRevertBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OnchainRevertBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OnchainRevertBatch represents a RevertBatch event raised by the Onchain contract.
type OnchainRevertBatch struct {
	Arg0 [16]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRevertBatch is a free log retrieval operation binding the contract event 0xf335387fffae8d4c523007d8c373bffc677e94ca97652a6dabecea00f48e4dd6.
//
// Solidity: event RevertBatch(bytes16 arg0)
func (_Onchain *OnchainFilterer) FilterRevertBatch(opts *bind.FilterOpts) (*OnchainRevertBatchIterator, error) {

	logs, sub, err := _Onchain.contract.FilterLogs(opts, "RevertBatch")
	if err != nil {
		return nil, err
	}
	return &OnchainRevertBatchIterator{contract: _Onchain.contract, event: "RevertBatch", logs: logs, sub: sub}, nil
}

// WatchRevertBatch is a free log subscription operation binding the contract event 0xf335387fffae8d4c523007d8c373bffc677e94ca97652a6dabecea00f48e4dd6.
//
// Solidity: event RevertBatch(bytes16 arg0)
func (_Onchain *OnchainFilterer) WatchRevertBatch(opts *bind.WatchOpts, sink chan<- *OnchainRevertBatch) (event.Subscription, error) {

	logs, sub, err := _Onchain.contract.WatchLogs(opts, "RevertBatch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OnchainRevertBatch)
				if err := _Onchain.contract.UnpackLog(event, "RevertBatch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRevertBatch is a log parse operation binding the contract event 0xf335387fffae8d4c523007d8c373bffc677e94ca97652a6dabecea00f48e4dd6.
//
// Solidity: event RevertBatch(bytes16 arg0)
func (_Onchain *OnchainFilterer) ParseRevertBatch(log types.Log) (*OnchainRevertBatch, error) {
	event := new(OnchainRevertBatch)
	if err := _Onchain.contract.UnpackLog(event, "RevertBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OnchainWrongOrdersIterator is returned from FilterWrongOrders and is used to iterate over the raw logs and unpacked data for WrongOrders events raised by the Onchain contract.
type OnchainWrongOrdersIterator struct {
	Event *OnchainWrongOrders // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OnchainWrongOrdersIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OnchainWrongOrders)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OnchainWrongOrders)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OnchainWrongOrdersIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OnchainWrongOrdersIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OnchainWrongOrders represents a WrongOrders event raised by the Onchain contract.
type OnchainWrongOrders struct {
	Arg0 [16]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterWrongOrders is a free log retrieval operation binding the contract event 0x021a37001a82d22d35e305649d0bd569decbff74168cba17de2ec0ba197f2d82.
//
// Solidity: event WrongOrders(bytes16 arg0)
func (_Onchain *OnchainFilterer) FilterWrongOrders(opts *bind.FilterOpts) (*OnchainWrongOrdersIterator, error) {

	logs, sub, err := _Onchain.contract.FilterLogs(opts, "WrongOrders")
	if err != nil {
		return nil, err
	}
	return &OnchainWrongOrdersIterator{contract: _Onchain.contract, event: "WrongOrders", logs: logs, sub: sub}, nil
}

// WatchWrongOrders is a free log subscription operation binding the contract event 0x021a37001a82d22d35e305649d0bd569decbff74168cba17de2ec0ba197f2d82.
//
// Solidity: event WrongOrders(bytes16 arg0)
func (_Onchain *OnchainFilterer) WatchWrongOrders(opts *bind.WatchOpts, sink chan<- *OnchainWrongOrders) (event.Subscription, error) {

	logs, sub, err := _Onchain.contract.WatchLogs(opts, "WrongOrders")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OnchainWrongOrders)
				if err := _Onchain.contract.UnpackLog(event, "WrongOrders", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWrongOrders is a log parse operation binding the contract event 0x021a37001a82d22d35e305649d0bd569decbff74168cba17de2ec0ba197f2d82.
//
// Solidity: event WrongOrders(bytes16 arg0)
func (_Onchain *OnchainFilterer) ParseWrongOrders(log types.Log) (*OnchainWrongOrders, error) {
	event := new(OnchainWrongOrders)
	if err := _Onchain.contract.UnpackLog(event, "WrongOrders", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
