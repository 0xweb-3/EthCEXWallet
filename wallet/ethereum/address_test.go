package ethereum

import (
	"context"
	"crypto/ecdsa"
	"github.com/0xweb-3/EthCEXWallet/wallet/types"
	"github.com/ethereum/go-ethereum/crypto"
	"reflect"
	"testing"
)

func TestCreateAddress(t *testing.T) {
	addressResult, err := CreateAddress(context.Background())
	if err != nil {
		t.Error(err)
	} else {
		t.Log(addressResult)
	}
}

func TestCreateAddressByPrivateKey(t *testing.T) {
	// 已知数据
	expectedPublicKey := "ef66846f883b4aab5f5ccafd01ef6c4bf1180a2b127f3e7efa4dfa81fe77b3472abd6050d641ce4ca2c14603beba5ec2549cbba66db8b2ac70dea73ba7be95d0"
	expectedAddress := "0x01C9E6bdb351AD536236b508092f49eDEe5be0e6"

	// 使用一个已知私钥
	privateKeyHex := "100f4ac24a3eadf569b03628a858bc4c873573d31ecb91e5a49abe5d6d205ae9"
	privateKey, err := crypto.HexToECDSA(privateKeyHex)

	if err != nil {
		t.Fatalf("Failed to convert hex to private key: %v", err)
	}

	type args struct {
		ctx        context.Context
		privateKey *ecdsa.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    *types.EthAddress
		wantErr bool
	}{
		{
			name: "TestCreateAddressByPrivateKey",
			args: args{
				ctx:        context.Background(),
				privateKey: privateKey,
			},
			want: &types.EthAddress{
				PrivateKey: privateKeyHex,
				PublicKey:  expectedPublicKey,
				Address:    expectedAddress,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateAddressByPrivateKey(tt.args.ctx, tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAddressByPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 检查生成的地址和公钥
			if !reflect.DeepEqual(got.PublicKey, tt.want.PublicKey) {
				t.Errorf("CreateAddressByPrivateKey() public key = %v, want %v", got.PublicKey, tt.want.PublicKey)
			}
			if !reflect.DeepEqual(got.Address, tt.want.Address) {
				t.Errorf("CreateAddressByPrivateKey() address = %v, want %v", got.Address, tt.want.Address)
			}
		})
	}
}

/*
*
privateKey 17a01d2d0862c190dd3d286f5233039938c0522da31fd7d580569cdc07e642f4
publicKey a2e624a9cba7cb4b6b2814c8535b11a130f4818027f125c2ecdf23da303491bb822726e209563a1046b20d7dccfebd3857d4b5e9e3d079dcfbf4c9737fc06d18
address 0xEB80a127b2b763C631D8ADCeBb0976b190C8C227
*/
func TestGetAddressByPublicKey(t *testing.T) {
	addressStr, err := GetAddressByPublicKey(context.Background(), "2192ad7a4e8df85b6252f3d15e0b594dd09fe03c0885b5f81a420d5fe234fc76bddf2b27a4ef909e28ce4279bdacd20e83227a2231d6687bff8e340a0ecfea8a")
	if err != nil {
		t.Error(err)
	}
	t.Log(addressStr)
}

// {
//7b7ee3b7da1b07755293ffe455adbaeae5681cd203c931b30b0379d0b3bb223f
//2192ad7a4e8df85b6252f3d15e0b594dd09fe03c0885b5f81a420d5fe234fc76bddf2b27a4ef909e28ce4279bdacd20e83227a2231d6687bff8e340a0ecfea8a
//0xbe56653861fe159d046a95508105EC36F24fdA1D
//}
