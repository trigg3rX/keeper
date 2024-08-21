package core

import (
	"math/big"
	"errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	cstaskmanager "github.com/Layr-Labs/incredible-squaring-avs/contracts/bindings/IncredibleSquaringTaskManager"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"
)

// this hardcodes abi.encode() for cstaskmanager.IIncredibleSquaringTaskManagerTaskResponse
// unclear why abigen doesn't provide this out of the box...
func AbiEncodeTaskResponse(h *cstaskmanager.IIncredibleSquaringTaskManagerTaskResponse) ([]byte, error) {

	// The order here has to match the field ordering of cstaskmanager.IIncredibleSquaringTaskManagerTaskResponse
	taskResponseType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "referenceTaskIndex",
			Type: "uint32",
		},
		{
			Name: "numberSquared",
			Type: "uint256",
		},
	})
	if err != nil {
		return nil, err
	}
	arguments := abi.Arguments{
		{
			Type: taskResponseType,
		},
	}

	bytes, err := arguments.Pack(h)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

type SignedTaskResponse struct {
    JobID      uint32
    SignedData []byte
}

// GetTaskResponseDigest returns the hash of the TaskResponse, which is what operators sign over
func GetTaskResponseDigest(signedResponse *SignedTaskResponse) ([32]byte, error) {
    if signedResponse == nil {
        return [32]byte{}, errors.New("signedResponse is nil")
    }

    // Combine JobID and SignedData into a single byte slice
    jobIDBytes := common.LeftPadBytes([]byte{byte(signedResponse.JobID)}, 4)
    dataToHash := append(jobIDBytes, signedResponse.SignedData...)

    var taskResponseDigest [32]byte
    hasher := sha3.NewLegacyKeccak256()
    hasher.Write(dataToHash)
    copy(taskResponseDigest[:], hasher.Sum(nil)[:32])

    return taskResponseDigest, nil
}

// BINDING UTILS - conversion from contract structs to golang structs

// BN254.sol is a library, so bindings for G1 Points and G2 Points are only generated
// in every contract that imports that library. Thus the output here will need to be
// type casted if G1Point is needed to interface with another contract (eg: BLSPublicKeyCompendium.sol)
func ConvertToBN254G1Point(input *bls.G1Point) cstaskmanager.BN254G1Point {
	output := cstaskmanager.BN254G1Point{
		X: input.X.BigInt(big.NewInt(0)),
		Y: input.Y.BigInt(big.NewInt(0)),
	}
	return output
}

func ConvertToBN254G2Point(input *bls.G2Point) cstaskmanager.BN254G2Point {
	output := cstaskmanager.BN254G2Point{
		X: [2]*big.Int{input.X.A1.BigInt(big.NewInt(0)), input.X.A0.BigInt(big.NewInt(0))},
		Y: [2]*big.Int{input.Y.A1.BigInt(big.NewInt(0)), input.Y.A0.BigInt(big.NewInt(0))},
	}
	return output
}
