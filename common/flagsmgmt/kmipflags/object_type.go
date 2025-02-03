package kmipflags

import (
	"errors"
	"strings"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/ttlv"
)

type ObjectType kmip.ObjectType

const (
	ObjectTypeCertificate  = ObjectType(kmip.ObjectTypeCertificate)
	ObjectTypeSymmetricKey = ObjectType(kmip.ObjectTypeSymmetricKey)
	ObjectTypePublicKey    = ObjectType(kmip.ObjectTypePublicKey)
	ObjectTypePrivateKey   = ObjectType(kmip.ObjectTypePrivateKey)
	ObjectTypeSplitKey     = ObjectType(kmip.ObjectTypeSplitKey)
	//nolint:staticcheck // Needed to support legacy templates
	ObjectTypeTemplate     = ObjectType(kmip.ObjectTypeTemplate)
	ObjectTypeSecretData   = ObjectType(kmip.ObjectTypeSecretData)
	ObjectTypeOpaqueObject = ObjectType(kmip.ObjectTypeOpaqueObject)
	ObjectTypePGPKey       = ObjectType(kmip.ObjectTypePGPKey)
)

func (e *ObjectType) String() string {
	if *e == 0 {
		return ""
	}
	return ttlv.EnumStr(kmip.ObjectType(*e))
}

func (e *ObjectType) Set(v string) error {
	switch strings.ToLower(v) {
	case "certificate":
		*e = ObjectTypeCertificate
	case "symmetrickey":
		*e = ObjectTypeSymmetricKey
	case "publickey":
		*e = ObjectTypePublicKey
	case "privatekey":
		*e = ObjectTypePrivateKey
	case "splitkey":
		*e = ObjectTypeSplitKey
	case "template":
		*e = ObjectTypeTemplate
	case "secretdata":
		*e = ObjectTypeSecretData
	case "opaqueobject":
		*e = ObjectTypeOpaqueObject
	case "pgpkey":
		*e = ObjectTypePGPKey
	default:
		return errors.New(`must be one of "Certificate","SymmetricKey","PublicKey","PrivateKey","SplitKey","Template","SecretData","OpaqueObject","PGPKey"`)
	}
	return nil
}

func (e *ObjectType) Type() string {
	return "Certificate|SymmetricKey|PublicKey|PrivateKey|SplitKey|Template|SecretData|OpaqueObject|PGPKey"
}
