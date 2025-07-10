// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package credential

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

// CredentialMetaData contains all meta data concerning the Credential contract.
var CredentialMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ArrayLengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyRevokeStatusList\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyStatusListId\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TreeNotExists\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"revokeStatusListIds\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"revokeStatusLists\",\"type\":\"bytes[]\"}],\"name\":\"BatchRevokeStatusListUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"issuers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"treeIndices\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"newRoots\",\"type\":\"bytes32[]\"}],\"name\":\"BatchTreesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"revokeStatusListId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"revokeStatusList\",\"type\":\"bytes\"}],\"name\":\"RevokeStatusListUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"treeIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"newRoot\",\"type\":\"bytes32\"}],\"name\":\"TreeUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WRITER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"revokeStatusListIds\",\"type\":\"string[]\"},{\"internalType\":\"bytes[]\",\"name\":\"revokeStatusLists\",\"type\":\"bytes[]\"}],\"name\":\"batchUpdateRevokeStatusList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"issuers\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"treeIndices\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"newRoots\",\"type\":\"bytes32[]\"}],\"name\":\"batchUpdateTreeRoots\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"treeIndex\",\"type\":\"uint256\"}],\"name\":\"getTreeRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"issuerTrees\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"treeIndex\",\"type\":\"uint256\"}],\"name\":\"treeExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"revokeStatusListId\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"revokeStatusList\",\"type\":\"bytes\"}],\"name\":\"updateRevokeStatusList\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"treeIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"newRoot\",\"type\":\"bytes32\"}],\"name\":\"updateTreeRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"treeIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"leaf\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"}],\"name\":\"verifyVC\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f80fd5b506100225f801b3361005960201b60201c565b506100537f2b8f168f361ac1393a163ed4adfa899a87be7b7c71645167bdaddd822ae453c83361005960201b60201c565b506101b8565b5f61006a838361014e60201b60201c565b6101445760015f808581526020019081526020015f205f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055506100e16101b160201b60201c565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a460019050610148565b5f90505b92915050565b5f805f8481526020019081526020015f205f015f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b5f33905090565b61202e806101c55f395ff3fe608060405234801561000f575f80fd5b50600436106100fe575f3560e01c80639b1c499311610095578063d547741f11610064578063d547741f146102ba578063ed6cf0f3146102d6578063f19d2fda14610306578063f1a4d25c14610322576100fe565b80639b1c4993146102465780639beaab7b14610262578063a217fddf14610280578063c12c77921461029e576100fe565b806336568abe116100d157806336568abe146101ae5780635e51777d146101ca5780636dea4069146101fa57806391d1485414610216576100fe565b806301ffc9a7146101025780631ce7f74a14610132578063248a9ca3146101625780632f2ff15d14610192575b5f80fd5b61011c600480360381019061011791906111f4565b610352565b6040516101299190611239565b60405180910390f35b61014c600480360381019061014791906112df565b6103cb565b6040516101599190611335565b60405180910390f35b61017c60048036038101906101779190611378565b6103eb565b6040516101899190611335565b60405180910390f35b6101ac60048036038101906101a791906113a3565b610407565b005b6101c860048036038101906101c391906113a3565b610429565b005b6101e460048036038101906101df91906112df565b6104a4565b6040516101f19190611239565b60405180910390f35b610214600480360381019061020f9190611497565b6104ff565b005b610230600480360381019061022b91906113a3565b6106af565b60405161023d9190611239565b60405180910390f35b610260600480360381019061025b91906115d2565b610712565b005b61026a61094e565b6040516102779190611335565b60405180910390f35b610288610972565b6040516102959190611335565b60405180910390f35b6102b860048036038101906102b39190611663565b610978565b005b6102d460048036038101906102cf91906113a3565b610b21565b005b6102f060048036038101906102eb9190611708565b610b43565b6040516102fd9190611239565b60405180910390f35b610320600480360381019061031b9190611836565b610c26565b005b61033c600480360381019061033791906112df565b610df8565b6040516103499190611335565b60405180910390f35b5f7f7965db0b000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614806103c457506103c382610e4e565b5b9050919050565b6001602052815f5260405f20602052805f5260405f205f91509150505481565b5f805f8381526020019081526020015f20600101549050919050565b610410826103eb565b61041981610eb7565b6104238383610ecb565b50505050565b610431610fb4565b73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614610495576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61049f8282610fbb565b505050565b5f805f1b60015f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8481526020019081526020015f20541415905092915050565b8461052a7f2b8f168f361ac1393a163ed4adfa899a87be7b7c71645167bdaddd822ae453c8336106af565b15801561056357508073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b156105c757337f2b8f168f361ac1393a163ed4adfa899a87be7b7c71645167bdaddd822ae453c86040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016105be9291906118f5565b60405180910390fd5b5f8585905003610603576040517f7e8a733100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f838390500361063f576040517f7a939b5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b848460405161064f929190611958565b60405180910390208673ffffffffffffffffffffffffffffffffffffffff167f3b6e4e245d32d84cabec6fadf45748c59b54691dd716e6557ef2b4785e1e955a858560405161069f9291906119bc565b60405180910390a3505050505050565b5f805f8481526020019081526020015f205f015f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16905092915050565b8461073d7f2b8f168f361ac1393a163ed4adfa899a87be7b7c71645167bdaddd822ae453c8336106af565b15801561077657508073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b156107da57337f2b8f168f361ac1393a163ed4adfa899a87be7b7c71645167bdaddd822ae453c86040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016107d19291906118f5565b60405180910390fd5b5f85859050905083839050811461081d576040517fa24a13a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f5b818110156108f0575f87878381811061083b5761083a6119de565b5b905060200281019061084d9190611a17565b905003610886576040517f7e8a733100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f85858381811061089a576108996119de565b5b90506020028101906108ac9190611a79565b9050036108e5576040517f7a939b5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600101905061081f565b508673ffffffffffffffffffffffffffffffffffffffff167f2acb167af461179dbce836fad730c733e326b88b8cc56c5b17f25b2e8bab8eed8787878760405161093d9493929190611d7f565b60405180910390a250505050505050565b7f2b8f168f361ac1393a163ed4adfa899a87be7b7c71645167bdaddd822ae453c881565b5f801b81565b826109a37f2b8f168f361ac1393a163ed4adfa899a87be7b7c71645167bdaddd822ae453c8336106af565b1580156109dc57508073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b15610a4057337f2b8f168f361ac1393a163ed4adfa899a87be7b7c71645167bdaddd822ae453c86040517fe2517d3f000000000000000000000000000000000000000000000000000000008152600401610a379291906118f5565b60405180910390fd5b5f801b8203610a7b576040517f53ce4ece00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8160015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8581526020019081526020015f2081905550828473ffffffffffffffffffffffffffffffffffffffff167f6359763dd97d67c7b79a119f0e38c8d995c8b2fd50f10d53b24d9949ca132fdb84604051610b139190611335565b60405180910390a350505050565b610b2a826103eb565b610b3381610eb7565b610b3d8383610fbb565b50505050565b5f8060015f8873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8781526020019081526020015f205490505f801b8103610bcf576040517ff05d647700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610c1a8484808060200260200160405190810160405280939291908181526020018383602002808284375f81840152601f19601f8201169050808301925050505050505082876110a4565b91505095945050505050565b7f2b8f168f361ac1393a163ed4adfa899a87be7b7c71645167bdaddd822ae453c8610c5081610eb7565b5f8787905090508585905081141580610c6c5750838390508114155b15610ca3576040517fa24a13a600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f600190505f5b82811015610dab575f8a8a83818110610cc657610cc56119de565b5b9050602002016020810190610cdb9190611db8565b90505f878784818110610cf157610cf06119de565b5b9050602002013590505f801b8103610d35576040517f53ce4ece00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80845f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8c8c87818110610d8657610d856119de565b5b9050602002013581526020019081526020015f20819055508260010192505050610caa565b507fdfca0e820844f18511e987f70077b11bc4543e954283a833a15d51abdbbc7cd0898989898989604051610de596959493929190611f82565b60405180910390a1505050505050505050565b5f60015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8381526020019081526020015f2054905092915050565b5f7f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b610ec881610ec3610fb4565b6110ba565b50565b5f610ed683836106af565b610faa5760015f808581526020019081526020015f205f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550610f47610fb4565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a460019050610fae565b5f90505b92915050565b5f33905090565b5f610fc683836106af565b1561109a575f805f8581526020019081526020015f205f015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff021916908315150217905550611037610fb4565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b60405160405180910390a46001905061109e565b5f90505b92915050565b5f826110b0858461110b565b1490509392505050565b6110c482826106af565b6111075780826040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526004016110fe9291906118f5565b60405180910390fd5b5050565b5f808290505f5b845181101561114e5761113f82868381518110611132576111316119de565b5b6020026020010151611159565b91508080600101915050611112565b508091505092915050565b5f8183106111705761116b8284611183565b61117b565b61117a8383611183565b5b905092915050565b5f825f528160205260405f20905092915050565b5f80fd5b5f80fd5b5f7fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b6111d38161119f565b81146111dd575f80fd5b50565b5f813590506111ee816111ca565b92915050565b5f6020828403121561120957611208611197565b5b5f611216848285016111e0565b91505092915050565b5f8115159050919050565b6112338161121f565b82525050565b5f60208201905061124c5f83018461122a565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61127b82611252565b9050919050565b61128b81611271565b8114611295575f80fd5b50565b5f813590506112a681611282565b92915050565b5f819050919050565b6112be816112ac565b81146112c8575f80fd5b50565b5f813590506112d9816112b5565b92915050565b5f80604083850312156112f5576112f4611197565b5b5f61130285828601611298565b9250506020611313858286016112cb565b9150509250929050565b5f819050919050565b61132f8161131d565b82525050565b5f6020820190506113485f830184611326565b92915050565b6113578161131d565b8114611361575f80fd5b50565b5f813590506113728161134e565b92915050565b5f6020828403121561138d5761138c611197565b5b5f61139a84828501611364565b91505092915050565b5f80604083850312156113b9576113b8611197565b5b5f6113c685828601611364565b92505060206113d785828601611298565b9150509250929050565b5f80fd5b5f80fd5b5f80fd5b5f8083601f840112611402576114016113e1565b5b8235905067ffffffffffffffff81111561141f5761141e6113e5565b5b60208301915083600182028301111561143b5761143a6113e9565b5b9250929050565b5f8083601f840112611457576114566113e1565b5b8235905067ffffffffffffffff811115611474576114736113e5565b5b6020830191508360018202830111156114905761148f6113e9565b5b9250929050565b5f805f805f606086880312156114b0576114af611197565b5b5f6114bd88828901611298565b955050602086013567ffffffffffffffff8111156114de576114dd61119b565b5b6114ea888289016113ed565b9450945050604086013567ffffffffffffffff81111561150d5761150c61119b565b5b61151988828901611442565b92509250509295509295909350565b5f8083601f84011261153d5761153c6113e1565b5b8235905067ffffffffffffffff81111561155a576115596113e5565b5b602083019150836020820283011115611576576115756113e9565b5b9250929050565b5f8083601f840112611592576115916113e1565b5b8235905067ffffffffffffffff8111156115af576115ae6113e5565b5b6020830191508360208202830111156115cb576115ca6113e9565b5b9250929050565b5f805f805f606086880312156115eb576115ea611197565b5b5f6115f888828901611298565b955050602086013567ffffffffffffffff8111156116195761161861119b565b5b61162588828901611528565b9450945050604086013567ffffffffffffffff8111156116485761164761119b565b5b6116548882890161157d565b92509250509295509295909350565b5f805f6060848603121561167a57611679611197565b5b5f61168786828701611298565b9350506020611698868287016112cb565b92505060406116a986828701611364565b9150509250925092565b5f8083601f8401126116c8576116c76113e1565b5b8235905067ffffffffffffffff8111156116e5576116e46113e5565b5b602083019150836020820283011115611701576117006113e9565b5b9250929050565b5f805f805f6080868803121561172157611720611197565b5b5f61172e88828901611298565b955050602061173f888289016112cb565b945050604061175088828901611364565b935050606086013567ffffffffffffffff8111156117715761177061119b565b5b61177d888289016116b3565b92509250509295509295909350565b5f8083601f8401126117a1576117a06113e1565b5b8235905067ffffffffffffffff8111156117be576117bd6113e5565b5b6020830191508360208202830111156117da576117d96113e9565b5b9250929050565b5f8083601f8401126117f6576117f56113e1565b5b8235905067ffffffffffffffff811115611813576118126113e5565b5b60208301915083602082028301111561182f5761182e6113e9565b5b9250929050565b5f805f805f80606087890312156118505761184f611197565b5b5f87013567ffffffffffffffff81111561186d5761186c61119b565b5b61187989828a0161178c565b9650965050602087013567ffffffffffffffff81111561189c5761189b61119b565b5b6118a889828a016117e1565b9450945050604087013567ffffffffffffffff8111156118cb576118ca61119b565b5b6118d789828a016116b3565b92509250509295509295509295565b6118ef81611271565b82525050565b5f6040820190506119085f8301856118e6565b6119156020830184611326565b9392505050565b5f81905092915050565b828183375f83830152505050565b5f61193f838561191c565b935061194c838584611926565b82840190509392505050565b5f611964828486611934565b91508190509392505050565b5f82825260208201905092915050565b5f601f19601f8301169050919050565b5f61199b8385611970565b93506119a8838584611926565b6119b183611980565b840190509392505050565b5f6020820190508181035f8301526119d5818486611990565b90509392505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f80fd5b5f80fd5b5f80fd5b5f8083356001602003843603038112611a3357611a32611a0b565b5b80840192508235915067ffffffffffffffff821115611a5557611a54611a0f565b5b602083019250600182023603831315611a7157611a70611a13565b5b509250929050565b5f8083356001602003843603038112611a9557611a94611a0b565b5b80840192508235915067ffffffffffffffff821115611ab757611ab6611a0f565b5b602083019250600182023603831315611ad357611ad2611a13565b5b509250929050565b5f82825260208201905092915050565b5f819050919050565b5f82825260208201905092915050565b5f611b0f8385611af4565b9350611b1c838584611926565b611b2583611980565b840190509392505050565b5f611b3c848484611b04565b90509392505050565b5f80fd5b5f80fd5b5f80fd5b5f8083356001602003843603038112611b6d57611b6c611b4d565b5b83810192508235915060208301925067ffffffffffffffff821115611b9557611b94611b45565b5b600182023603831315611bab57611baa611b49565b5b509250929050565b5f602082019050919050565b5f611bca8385611adb565b935083602084028501611bdc84611aeb565b805f5b87811015611c21578484038952611bf68284611b51565b611c01868284611b30565b9550611c0c84611bb3565b935060208b019a505050600181019050611bdf565b50829750879450505050509392505050565b5f82825260208201905092915050565b5f819050919050565b5f82825260208201905092915050565b5f611c678385611c4c565b9350611c74838584611926565b611c7d83611980565b840190509392505050565b5f611c94848484611c5c565b90509392505050565b5f8083356001602003843603038112611cb957611cb8611b4d565b5b83810192508235915060208301925067ffffffffffffffff821115611ce157611ce0611b45565b5b600182023603831315611cf757611cf6611b49565b5b509250929050565b5f602082019050919050565b5f611d168385611c33565b935083602084028501611d2884611c43565b805f5b87811015611d6d578484038952611d428284611c9d565b611d4d868284611c88565b9550611d5884611cff565b935060208b019a505050600181019050611d2b565b50829750879450505050509392505050565b5f6040820190508181035f830152611d98818688611bbf565b90508181036020830152611dad818486611d0b565b905095945050505050565b5f60208284031215611dcd57611dcc611197565b5b5f611dda84828501611298565b91505092915050565b5f82825260208201905092915050565b5f819050919050565b611e0581611271565b82525050565b5f611e168383611dfc565b60208301905092915050565b5f611e306020840184611298565b905092915050565b5f602082019050919050565b5f611e4f8385611de3565b9350611e5a82611df3565b805f5b85811015611e9257611e6f8284611e22565b611e798882611e0b565b9750611e8483611e38565b925050600181019050611e5d565b5085925050509392505050565b5f82825260208201905092915050565b5f80fd5b82818337505050565b5f611ec78385611e9f565b93507f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115611efa57611ef9611eaf565b5b602083029250611f0b838584611eb3565b82840190509392505050565b5f82825260208201905092915050565b5f611f328385611f17565b93507f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115611f6557611f64611eaf565b5b602083029250611f76838584611eb3565b82840190509392505050565b5f6060820190508181035f830152611f9b81888a611e44565b90508181036020830152611fb0818688611ebc565b90508181036040830152611fc5818486611f27565b905097965050505050505056fea2646970667358221220f6c40033e39b397bf95b1b0f6407c613d5fe233655d0fb4b4bf9f8f5310b775a64736f6c637828302e382e32352d646576656c6f702e323032342e322e32342b636f6d6d69742e64626137353465630059",
}

// CredentialABI is the input ABI used to generate the binding from.
// Deprecated: Use CredentialMetaData.ABI instead.
var CredentialABI = CredentialMetaData.ABI

// CredentialBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CredentialMetaData.Bin instead.
var CredentialBin = CredentialMetaData.Bin

// DeployCredential deploys a new Ethereum contract, binding an instance of Credential to it.
func DeployCredential(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Credential, error) {
	parsed, err := CredentialMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CredentialBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Credential{CredentialCaller: CredentialCaller{contract: contract}, CredentialTransactor: CredentialTransactor{contract: contract}, CredentialFilterer: CredentialFilterer{contract: contract}}, nil
}

// Credential is an auto generated Go binding around an Ethereum contract.
type Credential struct {
	CredentialCaller     // Read-only binding to the contract
	CredentialTransactor // Write-only binding to the contract
	CredentialFilterer   // Log filterer for contract events
}

// CredentialCaller is an auto generated read-only Go binding around an Ethereum contract.
type CredentialCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CredentialTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CredentialTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CredentialFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CredentialFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CredentialSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CredentialSession struct {
	Contract     *Credential       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CredentialCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CredentialCallerSession struct {
	Contract *CredentialCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// CredentialTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CredentialTransactorSession struct {
	Contract     *CredentialTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// CredentialRaw is an auto generated low-level Go binding around an Ethereum contract.
type CredentialRaw struct {
	Contract *Credential // Generic contract binding to access the raw methods on
}

// CredentialCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CredentialCallerRaw struct {
	Contract *CredentialCaller // Generic read-only contract binding to access the raw methods on
}

// CredentialTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CredentialTransactorRaw struct {
	Contract *CredentialTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCredential creates a new instance of Credential, bound to a specific deployed contract.
func NewCredential(address common.Address, backend bind.ContractBackend) (*Credential, error) {
	contract, err := bindCredential(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Credential{CredentialCaller: CredentialCaller{contract: contract}, CredentialTransactor: CredentialTransactor{contract: contract}, CredentialFilterer: CredentialFilterer{contract: contract}}, nil
}

// NewCredentialCaller creates a new read-only instance of Credential, bound to a specific deployed contract.
func NewCredentialCaller(address common.Address, caller bind.ContractCaller) (*CredentialCaller, error) {
	contract, err := bindCredential(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CredentialCaller{contract: contract}, nil
}

// NewCredentialTransactor creates a new write-only instance of Credential, bound to a specific deployed contract.
func NewCredentialTransactor(address common.Address, transactor bind.ContractTransactor) (*CredentialTransactor, error) {
	contract, err := bindCredential(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CredentialTransactor{contract: contract}, nil
}

// NewCredentialFilterer creates a new log filterer instance of Credential, bound to a specific deployed contract.
func NewCredentialFilterer(address common.Address, filterer bind.ContractFilterer) (*CredentialFilterer, error) {
	contract, err := bindCredential(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CredentialFilterer{contract: contract}, nil
}

// bindCredential binds a generic wrapper to an already deployed contract.
func bindCredential(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CredentialMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Credential *CredentialRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Credential.Contract.CredentialCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Credential *CredentialRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Credential.Contract.CredentialTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Credential *CredentialRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Credential.Contract.CredentialTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Credential *CredentialCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Credential.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Credential *CredentialTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Credential.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Credential *CredentialTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Credential.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Credential *CredentialCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Credential.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Credential *CredentialSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Credential.Contract.DEFAULTADMINROLE(&_Credential.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Credential *CredentialCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Credential.Contract.DEFAULTADMINROLE(&_Credential.CallOpts)
}

// WRITERROLE is a free data retrieval call binding the contract method 0x9beaab7b.
//
// Solidity: function WRITER_ROLE() view returns(bytes32)
func (_Credential *CredentialCaller) WRITERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Credential.contract.Call(opts, &out, "WRITER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WRITERROLE is a free data retrieval call binding the contract method 0x9beaab7b.
//
// Solidity: function WRITER_ROLE() view returns(bytes32)
func (_Credential *CredentialSession) WRITERROLE() ([32]byte, error) {
	return _Credential.Contract.WRITERROLE(&_Credential.CallOpts)
}

// WRITERROLE is a free data retrieval call binding the contract method 0x9beaab7b.
//
// Solidity: function WRITER_ROLE() view returns(bytes32)
func (_Credential *CredentialCallerSession) WRITERROLE() ([32]byte, error) {
	return _Credential.Contract.WRITERROLE(&_Credential.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Credential *CredentialCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Credential.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Credential *CredentialSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Credential.Contract.GetRoleAdmin(&_Credential.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Credential *CredentialCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Credential.Contract.GetRoleAdmin(&_Credential.CallOpts, role)
}

// GetTreeRoot is a free data retrieval call binding the contract method 0xf1a4d25c.
//
// Solidity: function getTreeRoot(address issuer, uint256 treeIndex) view returns(bytes32 root)
func (_Credential *CredentialCaller) GetTreeRoot(opts *bind.CallOpts, issuer common.Address, treeIndex *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Credential.contract.Call(opts, &out, "getTreeRoot", issuer, treeIndex)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetTreeRoot is a free data retrieval call binding the contract method 0xf1a4d25c.
//
// Solidity: function getTreeRoot(address issuer, uint256 treeIndex) view returns(bytes32 root)
func (_Credential *CredentialSession) GetTreeRoot(issuer common.Address, treeIndex *big.Int) ([32]byte, error) {
	return _Credential.Contract.GetTreeRoot(&_Credential.CallOpts, issuer, treeIndex)
}

// GetTreeRoot is a free data retrieval call binding the contract method 0xf1a4d25c.
//
// Solidity: function getTreeRoot(address issuer, uint256 treeIndex) view returns(bytes32 root)
func (_Credential *CredentialCallerSession) GetTreeRoot(issuer common.Address, treeIndex *big.Int) ([32]byte, error) {
	return _Credential.Contract.GetTreeRoot(&_Credential.CallOpts, issuer, treeIndex)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Credential *CredentialCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Credential.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Credential *CredentialSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Credential.Contract.HasRole(&_Credential.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Credential *CredentialCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Credential.Contract.HasRole(&_Credential.CallOpts, role, account)
}

// IssuerTrees is a free data retrieval call binding the contract method 0x1ce7f74a.
//
// Solidity: function issuerTrees(address , uint256 ) view returns(bytes32)
func (_Credential *CredentialCaller) IssuerTrees(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Credential.contract.Call(opts, &out, "issuerTrees", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// IssuerTrees is a free data retrieval call binding the contract method 0x1ce7f74a.
//
// Solidity: function issuerTrees(address , uint256 ) view returns(bytes32)
func (_Credential *CredentialSession) IssuerTrees(arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	return _Credential.Contract.IssuerTrees(&_Credential.CallOpts, arg0, arg1)
}

// IssuerTrees is a free data retrieval call binding the contract method 0x1ce7f74a.
//
// Solidity: function issuerTrees(address , uint256 ) view returns(bytes32)
func (_Credential *CredentialCallerSession) IssuerTrees(arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	return _Credential.Contract.IssuerTrees(&_Credential.CallOpts, arg0, arg1)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Credential *CredentialCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Credential.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Credential *CredentialSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Credential.Contract.SupportsInterface(&_Credential.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Credential *CredentialCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Credential.Contract.SupportsInterface(&_Credential.CallOpts, interfaceId)
}

// TreeExists is a free data retrieval call binding the contract method 0x5e51777d.
//
// Solidity: function treeExists(address issuer, uint256 treeIndex) view returns(bool)
func (_Credential *CredentialCaller) TreeExists(opts *bind.CallOpts, issuer common.Address, treeIndex *big.Int) (bool, error) {
	var out []interface{}
	err := _Credential.contract.Call(opts, &out, "treeExists", issuer, treeIndex)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// TreeExists is a free data retrieval call binding the contract method 0x5e51777d.
//
// Solidity: function treeExists(address issuer, uint256 treeIndex) view returns(bool)
func (_Credential *CredentialSession) TreeExists(issuer common.Address, treeIndex *big.Int) (bool, error) {
	return _Credential.Contract.TreeExists(&_Credential.CallOpts, issuer, treeIndex)
}

// TreeExists is a free data retrieval call binding the contract method 0x5e51777d.
//
// Solidity: function treeExists(address issuer, uint256 treeIndex) view returns(bool)
func (_Credential *CredentialCallerSession) TreeExists(issuer common.Address, treeIndex *big.Int) (bool, error) {
	return _Credential.Contract.TreeExists(&_Credential.CallOpts, issuer, treeIndex)
}

// VerifyVC is a free data retrieval call binding the contract method 0xed6cf0f3.
//
// Solidity: function verifyVC(address issuer, uint256 treeIndex, bytes32 leaf, bytes32[] proof) view returns(bool)
func (_Credential *CredentialCaller) VerifyVC(opts *bind.CallOpts, issuer common.Address, treeIndex *big.Int, leaf [32]byte, proof [][32]byte) (bool, error) {
	var out []interface{}
	err := _Credential.contract.Call(opts, &out, "verifyVC", issuer, treeIndex, leaf, proof)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyVC is a free data retrieval call binding the contract method 0xed6cf0f3.
//
// Solidity: function verifyVC(address issuer, uint256 treeIndex, bytes32 leaf, bytes32[] proof) view returns(bool)
func (_Credential *CredentialSession) VerifyVC(issuer common.Address, treeIndex *big.Int, leaf [32]byte, proof [][32]byte) (bool, error) {
	return _Credential.Contract.VerifyVC(&_Credential.CallOpts, issuer, treeIndex, leaf, proof)
}

// VerifyVC is a free data retrieval call binding the contract method 0xed6cf0f3.
//
// Solidity: function verifyVC(address issuer, uint256 treeIndex, bytes32 leaf, bytes32[] proof) view returns(bool)
func (_Credential *CredentialCallerSession) VerifyVC(issuer common.Address, treeIndex *big.Int, leaf [32]byte, proof [][32]byte) (bool, error) {
	return _Credential.Contract.VerifyVC(&_Credential.CallOpts, issuer, treeIndex, leaf, proof)
}

// BatchUpdateRevokeStatusList is a paid mutator transaction binding the contract method 0x9b1c4993.
//
// Solidity: function batchUpdateRevokeStatusList(address issuer, string[] revokeStatusListIds, bytes[] revokeStatusLists) returns()
func (_Credential *CredentialTransactor) BatchUpdateRevokeStatusList(opts *bind.TransactOpts, issuer common.Address, revokeStatusListIds []string, revokeStatusLists [][]byte) (*types.Transaction, error) {
	return _Credential.contract.Transact(opts, "batchUpdateRevokeStatusList", issuer, revokeStatusListIds, revokeStatusLists)
}

// BatchUpdateRevokeStatusList is a paid mutator transaction binding the contract method 0x9b1c4993.
//
// Solidity: function batchUpdateRevokeStatusList(address issuer, string[] revokeStatusListIds, bytes[] revokeStatusLists) returns()
func (_Credential *CredentialSession) BatchUpdateRevokeStatusList(issuer common.Address, revokeStatusListIds []string, revokeStatusLists [][]byte) (*types.Transaction, error) {
	return _Credential.Contract.BatchUpdateRevokeStatusList(&_Credential.TransactOpts, issuer, revokeStatusListIds, revokeStatusLists)
}

// BatchUpdateRevokeStatusList is a paid mutator transaction binding the contract method 0x9b1c4993.
//
// Solidity: function batchUpdateRevokeStatusList(address issuer, string[] revokeStatusListIds, bytes[] revokeStatusLists) returns()
func (_Credential *CredentialTransactorSession) BatchUpdateRevokeStatusList(issuer common.Address, revokeStatusListIds []string, revokeStatusLists [][]byte) (*types.Transaction, error) {
	return _Credential.Contract.BatchUpdateRevokeStatusList(&_Credential.TransactOpts, issuer, revokeStatusListIds, revokeStatusLists)
}

// BatchUpdateTreeRoots is a paid mutator transaction binding the contract method 0xf19d2fda.
//
// Solidity: function batchUpdateTreeRoots(address[] issuers, uint256[] treeIndices, bytes32[] newRoots) returns()
func (_Credential *CredentialTransactor) BatchUpdateTreeRoots(opts *bind.TransactOpts, issuers []common.Address, treeIndices []*big.Int, newRoots [][32]byte) (*types.Transaction, error) {
	return _Credential.contract.Transact(opts, "batchUpdateTreeRoots", issuers, treeIndices, newRoots)
}

// BatchUpdateTreeRoots is a paid mutator transaction binding the contract method 0xf19d2fda.
//
// Solidity: function batchUpdateTreeRoots(address[] issuers, uint256[] treeIndices, bytes32[] newRoots) returns()
func (_Credential *CredentialSession) BatchUpdateTreeRoots(issuers []common.Address, treeIndices []*big.Int, newRoots [][32]byte) (*types.Transaction, error) {
	return _Credential.Contract.BatchUpdateTreeRoots(&_Credential.TransactOpts, issuers, treeIndices, newRoots)
}

// BatchUpdateTreeRoots is a paid mutator transaction binding the contract method 0xf19d2fda.
//
// Solidity: function batchUpdateTreeRoots(address[] issuers, uint256[] treeIndices, bytes32[] newRoots) returns()
func (_Credential *CredentialTransactorSession) BatchUpdateTreeRoots(issuers []common.Address, treeIndices []*big.Int, newRoots [][32]byte) (*types.Transaction, error) {
	return _Credential.Contract.BatchUpdateTreeRoots(&_Credential.TransactOpts, issuers, treeIndices, newRoots)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Credential *CredentialTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Credential.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Credential *CredentialSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Credential.Contract.GrantRole(&_Credential.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Credential *CredentialTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Credential.Contract.GrantRole(&_Credential.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Credential *CredentialTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Credential.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Credential *CredentialSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Credential.Contract.RenounceRole(&_Credential.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Credential *CredentialTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Credential.Contract.RenounceRole(&_Credential.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Credential *CredentialTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Credential.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Credential *CredentialSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Credential.Contract.RevokeRole(&_Credential.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Credential *CredentialTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Credential.Contract.RevokeRole(&_Credential.TransactOpts, role, account)
}

// UpdateRevokeStatusList is a paid mutator transaction binding the contract method 0x6dea4069.
//
// Solidity: function updateRevokeStatusList(address issuer, string revokeStatusListId, bytes revokeStatusList) returns()
func (_Credential *CredentialTransactor) UpdateRevokeStatusList(opts *bind.TransactOpts, issuer common.Address, revokeStatusListId string, revokeStatusList []byte) (*types.Transaction, error) {
	return _Credential.contract.Transact(opts, "updateRevokeStatusList", issuer, revokeStatusListId, revokeStatusList)
}

// UpdateRevokeStatusList is a paid mutator transaction binding the contract method 0x6dea4069.
//
// Solidity: function updateRevokeStatusList(address issuer, string revokeStatusListId, bytes revokeStatusList) returns()
func (_Credential *CredentialSession) UpdateRevokeStatusList(issuer common.Address, revokeStatusListId string, revokeStatusList []byte) (*types.Transaction, error) {
	return _Credential.Contract.UpdateRevokeStatusList(&_Credential.TransactOpts, issuer, revokeStatusListId, revokeStatusList)
}

// UpdateRevokeStatusList is a paid mutator transaction binding the contract method 0x6dea4069.
//
// Solidity: function updateRevokeStatusList(address issuer, string revokeStatusListId, bytes revokeStatusList) returns()
func (_Credential *CredentialTransactorSession) UpdateRevokeStatusList(issuer common.Address, revokeStatusListId string, revokeStatusList []byte) (*types.Transaction, error) {
	return _Credential.Contract.UpdateRevokeStatusList(&_Credential.TransactOpts, issuer, revokeStatusListId, revokeStatusList)
}

// UpdateTreeRoot is a paid mutator transaction binding the contract method 0xc12c7792.
//
// Solidity: function updateTreeRoot(address issuer, uint256 treeIndex, bytes32 newRoot) returns()
func (_Credential *CredentialTransactor) UpdateTreeRoot(opts *bind.TransactOpts, issuer common.Address, treeIndex *big.Int, newRoot [32]byte) (*types.Transaction, error) {
	return _Credential.contract.Transact(opts, "updateTreeRoot", issuer, treeIndex, newRoot)
}

// UpdateTreeRoot is a paid mutator transaction binding the contract method 0xc12c7792.
//
// Solidity: function updateTreeRoot(address issuer, uint256 treeIndex, bytes32 newRoot) returns()
func (_Credential *CredentialSession) UpdateTreeRoot(issuer common.Address, treeIndex *big.Int, newRoot [32]byte) (*types.Transaction, error) {
	return _Credential.Contract.UpdateTreeRoot(&_Credential.TransactOpts, issuer, treeIndex, newRoot)
}

// UpdateTreeRoot is a paid mutator transaction binding the contract method 0xc12c7792.
//
// Solidity: function updateTreeRoot(address issuer, uint256 treeIndex, bytes32 newRoot) returns()
func (_Credential *CredentialTransactorSession) UpdateTreeRoot(issuer common.Address, treeIndex *big.Int, newRoot [32]byte) (*types.Transaction, error) {
	return _Credential.Contract.UpdateTreeRoot(&_Credential.TransactOpts, issuer, treeIndex, newRoot)
}

// CredentialBatchRevokeStatusListUpdatedIterator is returned from FilterBatchRevokeStatusListUpdated and is used to iterate over the raw logs and unpacked data for BatchRevokeStatusListUpdated events raised by the Credential contract.
type CredentialBatchRevokeStatusListUpdatedIterator struct {
	Event *CredentialBatchRevokeStatusListUpdated // Event containing the contract specifics and raw log

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
func (it *CredentialBatchRevokeStatusListUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CredentialBatchRevokeStatusListUpdated)
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
		it.Event = new(CredentialBatchRevokeStatusListUpdated)
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
func (it *CredentialBatchRevokeStatusListUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CredentialBatchRevokeStatusListUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CredentialBatchRevokeStatusListUpdated represents a BatchRevokeStatusListUpdated event raised by the Credential contract.
type CredentialBatchRevokeStatusListUpdated struct {
	Issuer              common.Address
	RevokeStatusListIds []string
	RevokeStatusLists   [][]byte
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterBatchRevokeStatusListUpdated is a free log retrieval operation binding the contract event 0x2acb167af461179dbce836fad730c733e326b88b8cc56c5b17f25b2e8bab8eed.
//
// Solidity: event BatchRevokeStatusListUpdated(address indexed issuer, string[] revokeStatusListIds, bytes[] revokeStatusLists)
func (_Credential *CredentialFilterer) FilterBatchRevokeStatusListUpdated(opts *bind.FilterOpts, issuer []common.Address) (*CredentialBatchRevokeStatusListUpdatedIterator, error) {

	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}

	logs, sub, err := _Credential.contract.FilterLogs(opts, "BatchRevokeStatusListUpdated", issuerRule)
	if err != nil {
		return nil, err
	}
	return &CredentialBatchRevokeStatusListUpdatedIterator{contract: _Credential.contract, event: "BatchRevokeStatusListUpdated", logs: logs, sub: sub}, nil
}

// WatchBatchRevokeStatusListUpdated is a free log subscription operation binding the contract event 0x2acb167af461179dbce836fad730c733e326b88b8cc56c5b17f25b2e8bab8eed.
//
// Solidity: event BatchRevokeStatusListUpdated(address indexed issuer, string[] revokeStatusListIds, bytes[] revokeStatusLists)
func (_Credential *CredentialFilterer) WatchBatchRevokeStatusListUpdated(opts *bind.WatchOpts, sink chan<- *CredentialBatchRevokeStatusListUpdated, issuer []common.Address) (event.Subscription, error) {

	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}

	logs, sub, err := _Credential.contract.WatchLogs(opts, "BatchRevokeStatusListUpdated", issuerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CredentialBatchRevokeStatusListUpdated)
				if err := _Credential.contract.UnpackLog(event, "BatchRevokeStatusListUpdated", log); err != nil {
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

// ParseBatchRevokeStatusListUpdated is a log parse operation binding the contract event 0x2acb167af461179dbce836fad730c733e326b88b8cc56c5b17f25b2e8bab8eed.
//
// Solidity: event BatchRevokeStatusListUpdated(address indexed issuer, string[] revokeStatusListIds, bytes[] revokeStatusLists)
func (_Credential *CredentialFilterer) ParseBatchRevokeStatusListUpdated(log types.Log) (*CredentialBatchRevokeStatusListUpdated, error) {
	event := new(CredentialBatchRevokeStatusListUpdated)
	if err := _Credential.contract.UnpackLog(event, "BatchRevokeStatusListUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CredentialBatchTreesUpdatedIterator is returned from FilterBatchTreesUpdated and is used to iterate over the raw logs and unpacked data for BatchTreesUpdated events raised by the Credential contract.
type CredentialBatchTreesUpdatedIterator struct {
	Event *CredentialBatchTreesUpdated // Event containing the contract specifics and raw log

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
func (it *CredentialBatchTreesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CredentialBatchTreesUpdated)
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
		it.Event = new(CredentialBatchTreesUpdated)
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
func (it *CredentialBatchTreesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CredentialBatchTreesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CredentialBatchTreesUpdated represents a BatchTreesUpdated event raised by the Credential contract.
type CredentialBatchTreesUpdated struct {
	Issuers     []common.Address
	TreeIndices []*big.Int
	NewRoots    [][32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBatchTreesUpdated is a free log retrieval operation binding the contract event 0xdfca0e820844f18511e987f70077b11bc4543e954283a833a15d51abdbbc7cd0.
//
// Solidity: event BatchTreesUpdated(address[] issuers, uint256[] treeIndices, bytes32[] newRoots)
func (_Credential *CredentialFilterer) FilterBatchTreesUpdated(opts *bind.FilterOpts) (*CredentialBatchTreesUpdatedIterator, error) {

	logs, sub, err := _Credential.contract.FilterLogs(opts, "BatchTreesUpdated")
	if err != nil {
		return nil, err
	}
	return &CredentialBatchTreesUpdatedIterator{contract: _Credential.contract, event: "BatchTreesUpdated", logs: logs, sub: sub}, nil
}

// WatchBatchTreesUpdated is a free log subscription operation binding the contract event 0xdfca0e820844f18511e987f70077b11bc4543e954283a833a15d51abdbbc7cd0.
//
// Solidity: event BatchTreesUpdated(address[] issuers, uint256[] treeIndices, bytes32[] newRoots)
func (_Credential *CredentialFilterer) WatchBatchTreesUpdated(opts *bind.WatchOpts, sink chan<- *CredentialBatchTreesUpdated) (event.Subscription, error) {

	logs, sub, err := _Credential.contract.WatchLogs(opts, "BatchTreesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CredentialBatchTreesUpdated)
				if err := _Credential.contract.UnpackLog(event, "BatchTreesUpdated", log); err != nil {
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

// ParseBatchTreesUpdated is a log parse operation binding the contract event 0xdfca0e820844f18511e987f70077b11bc4543e954283a833a15d51abdbbc7cd0.
//
// Solidity: event BatchTreesUpdated(address[] issuers, uint256[] treeIndices, bytes32[] newRoots)
func (_Credential *CredentialFilterer) ParseBatchTreesUpdated(log types.Log) (*CredentialBatchTreesUpdated, error) {
	event := new(CredentialBatchTreesUpdated)
	if err := _Credential.contract.UnpackLog(event, "BatchTreesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CredentialRevokeStatusListUpdatedIterator is returned from FilterRevokeStatusListUpdated and is used to iterate over the raw logs and unpacked data for RevokeStatusListUpdated events raised by the Credential contract.
type CredentialRevokeStatusListUpdatedIterator struct {
	Event *CredentialRevokeStatusListUpdated // Event containing the contract specifics and raw log

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
func (it *CredentialRevokeStatusListUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CredentialRevokeStatusListUpdated)
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
		it.Event = new(CredentialRevokeStatusListUpdated)
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
func (it *CredentialRevokeStatusListUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CredentialRevokeStatusListUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CredentialRevokeStatusListUpdated represents a RevokeStatusListUpdated event raised by the Credential contract.
type CredentialRevokeStatusListUpdated struct {
	Issuer             common.Address
	RevokeStatusListId common.Hash
	RevokeStatusList   []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterRevokeStatusListUpdated is a free log retrieval operation binding the contract event 0x3b6e4e245d32d84cabec6fadf45748c59b54691dd716e6557ef2b4785e1e955a.
//
// Solidity: event RevokeStatusListUpdated(address indexed issuer, string indexed revokeStatusListId, bytes revokeStatusList)
func (_Credential *CredentialFilterer) FilterRevokeStatusListUpdated(opts *bind.FilterOpts, issuer []common.Address, revokeStatusListId []string) (*CredentialRevokeStatusListUpdatedIterator, error) {

	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}
	var revokeStatusListIdRule []interface{}
	for _, revokeStatusListIdItem := range revokeStatusListId {
		revokeStatusListIdRule = append(revokeStatusListIdRule, revokeStatusListIdItem)
	}

	logs, sub, err := _Credential.contract.FilterLogs(opts, "RevokeStatusListUpdated", issuerRule, revokeStatusListIdRule)
	if err != nil {
		return nil, err
	}
	return &CredentialRevokeStatusListUpdatedIterator{contract: _Credential.contract, event: "RevokeStatusListUpdated", logs: logs, sub: sub}, nil
}

// WatchRevokeStatusListUpdated is a free log subscription operation binding the contract event 0x3b6e4e245d32d84cabec6fadf45748c59b54691dd716e6557ef2b4785e1e955a.
//
// Solidity: event RevokeStatusListUpdated(address indexed issuer, string indexed revokeStatusListId, bytes revokeStatusList)
func (_Credential *CredentialFilterer) WatchRevokeStatusListUpdated(opts *bind.WatchOpts, sink chan<- *CredentialRevokeStatusListUpdated, issuer []common.Address, revokeStatusListId []string) (event.Subscription, error) {

	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}
	var revokeStatusListIdRule []interface{}
	for _, revokeStatusListIdItem := range revokeStatusListId {
		revokeStatusListIdRule = append(revokeStatusListIdRule, revokeStatusListIdItem)
	}

	logs, sub, err := _Credential.contract.WatchLogs(opts, "RevokeStatusListUpdated", issuerRule, revokeStatusListIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CredentialRevokeStatusListUpdated)
				if err := _Credential.contract.UnpackLog(event, "RevokeStatusListUpdated", log); err != nil {
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

// ParseRevokeStatusListUpdated is a log parse operation binding the contract event 0x3b6e4e245d32d84cabec6fadf45748c59b54691dd716e6557ef2b4785e1e955a.
//
// Solidity: event RevokeStatusListUpdated(address indexed issuer, string indexed revokeStatusListId, bytes revokeStatusList)
func (_Credential *CredentialFilterer) ParseRevokeStatusListUpdated(log types.Log) (*CredentialRevokeStatusListUpdated, error) {
	event := new(CredentialRevokeStatusListUpdated)
	if err := _Credential.contract.UnpackLog(event, "RevokeStatusListUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CredentialRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Credential contract.
type CredentialRoleAdminChangedIterator struct {
	Event *CredentialRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *CredentialRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CredentialRoleAdminChanged)
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
		it.Event = new(CredentialRoleAdminChanged)
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
func (it *CredentialRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CredentialRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CredentialRoleAdminChanged represents a RoleAdminChanged event raised by the Credential contract.
type CredentialRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Credential *CredentialFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*CredentialRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Credential.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &CredentialRoleAdminChangedIterator{contract: _Credential.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Credential *CredentialFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *CredentialRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Credential.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CredentialRoleAdminChanged)
				if err := _Credential.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Credential *CredentialFilterer) ParseRoleAdminChanged(log types.Log) (*CredentialRoleAdminChanged, error) {
	event := new(CredentialRoleAdminChanged)
	if err := _Credential.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CredentialRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Credential contract.
type CredentialRoleGrantedIterator struct {
	Event *CredentialRoleGranted // Event containing the contract specifics and raw log

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
func (it *CredentialRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CredentialRoleGranted)
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
		it.Event = new(CredentialRoleGranted)
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
func (it *CredentialRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CredentialRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CredentialRoleGranted represents a RoleGranted event raised by the Credential contract.
type CredentialRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Credential *CredentialFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*CredentialRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Credential.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &CredentialRoleGrantedIterator{contract: _Credential.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Credential *CredentialFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *CredentialRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Credential.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CredentialRoleGranted)
				if err := _Credential.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Credential *CredentialFilterer) ParseRoleGranted(log types.Log) (*CredentialRoleGranted, error) {
	event := new(CredentialRoleGranted)
	if err := _Credential.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CredentialRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Credential contract.
type CredentialRoleRevokedIterator struct {
	Event *CredentialRoleRevoked // Event containing the contract specifics and raw log

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
func (it *CredentialRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CredentialRoleRevoked)
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
		it.Event = new(CredentialRoleRevoked)
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
func (it *CredentialRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CredentialRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CredentialRoleRevoked represents a RoleRevoked event raised by the Credential contract.
type CredentialRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Credential *CredentialFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*CredentialRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Credential.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &CredentialRoleRevokedIterator{contract: _Credential.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Credential *CredentialFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *CredentialRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Credential.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CredentialRoleRevoked)
				if err := _Credential.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Credential *CredentialFilterer) ParseRoleRevoked(log types.Log) (*CredentialRoleRevoked, error) {
	event := new(CredentialRoleRevoked)
	if err := _Credential.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CredentialTreeUpdatedIterator is returned from FilterTreeUpdated and is used to iterate over the raw logs and unpacked data for TreeUpdated events raised by the Credential contract.
type CredentialTreeUpdatedIterator struct {
	Event *CredentialTreeUpdated // Event containing the contract specifics and raw log

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
func (it *CredentialTreeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CredentialTreeUpdated)
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
		it.Event = new(CredentialTreeUpdated)
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
func (it *CredentialTreeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CredentialTreeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CredentialTreeUpdated represents a TreeUpdated event raised by the Credential contract.
type CredentialTreeUpdated struct {
	Issuer    common.Address
	TreeIndex *big.Int
	NewRoot   [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTreeUpdated is a free log retrieval operation binding the contract event 0x6359763dd97d67c7b79a119f0e38c8d995c8b2fd50f10d53b24d9949ca132fdb.
//
// Solidity: event TreeUpdated(address indexed issuer, uint256 indexed treeIndex, bytes32 newRoot)
func (_Credential *CredentialFilterer) FilterTreeUpdated(opts *bind.FilterOpts, issuer []common.Address, treeIndex []*big.Int) (*CredentialTreeUpdatedIterator, error) {

	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}
	var treeIndexRule []interface{}
	for _, treeIndexItem := range treeIndex {
		treeIndexRule = append(treeIndexRule, treeIndexItem)
	}

	logs, sub, err := _Credential.contract.FilterLogs(opts, "TreeUpdated", issuerRule, treeIndexRule)
	if err != nil {
		return nil, err
	}
	return &CredentialTreeUpdatedIterator{contract: _Credential.contract, event: "TreeUpdated", logs: logs, sub: sub}, nil
}

// WatchTreeUpdated is a free log subscription operation binding the contract event 0x6359763dd97d67c7b79a119f0e38c8d995c8b2fd50f10d53b24d9949ca132fdb.
//
// Solidity: event TreeUpdated(address indexed issuer, uint256 indexed treeIndex, bytes32 newRoot)
func (_Credential *CredentialFilterer) WatchTreeUpdated(opts *bind.WatchOpts, sink chan<- *CredentialTreeUpdated, issuer []common.Address, treeIndex []*big.Int) (event.Subscription, error) {

	var issuerRule []interface{}
	for _, issuerItem := range issuer {
		issuerRule = append(issuerRule, issuerItem)
	}
	var treeIndexRule []interface{}
	for _, treeIndexItem := range treeIndex {
		treeIndexRule = append(treeIndexRule, treeIndexItem)
	}

	logs, sub, err := _Credential.contract.WatchLogs(opts, "TreeUpdated", issuerRule, treeIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CredentialTreeUpdated)
				if err := _Credential.contract.UnpackLog(event, "TreeUpdated", log); err != nil {
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

// ParseTreeUpdated is a log parse operation binding the contract event 0x6359763dd97d67c7b79a119f0e38c8d995c8b2fd50f10d53b24d9949ca132fdb.
//
// Solidity: event TreeUpdated(address indexed issuer, uint256 indexed treeIndex, bytes32 newRoot)
func (_Credential *CredentialFilterer) ParseTreeUpdated(log types.Log) (*CredentialTreeUpdated, error) {
	event := new(CredentialTreeUpdated)
	if err := _Credential.contract.UnpackLog(event, "TreeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
