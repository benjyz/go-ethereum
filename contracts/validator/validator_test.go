package validator

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	addr   = crypto.PubkeyToAddress(key.PublicKey)
)

func TestValidator(t *testing.T) {
	contractBackend := backends.NewSimulatedBackend(core.GenesisAlloc{addr: {Balance: big.NewInt(1000000000)}})
	transactOpts := bind.NewKeyedTransactor(key)

	_, validator, err := DeployValidator(transactOpts, contractBackend)
	if err != nil {
		t.Fatalf("can't deploy root registry: %v", err)
	}
	contractBackend.Commit()

	candidates, err := validator.GetCandidates()
	if err != nil {
		t.Fatalf("can't get candidates: %v", err)
	}
	t.Log("candidates ", len(candidates))
	for _, it := range candidates {
		cap, _ := validator.GetCandidateCap(it)
		t.Log("candidate", it.String(), "cap", cap)
		owner, _ := validator.GetCandidateOwner(it)
		t.Log("candidate", it.String(), "validator owner", owner.String())
	}
	contractBackend.Commit()

	//someaddr := new common.addr{"0xf99805B536609cC03AcBB2604dFaC11E9E54a448"}
	//signers[common.HexToAddress("0x12f588d7d03bb269b382b842fc15d874e8c055a7")] = &rewardLog{5, new(big.Int).SetUint64(0)}
	someaddr := common.HexToAddress("0x31b249fE6F267aa2396Eb2DC36E9c79351d97Ec5")
	validator.Resign(someaddr)
	ncandidates, err := validator.GetCandidates()
	t.Log("candidates ", len(ncandidates))
}
