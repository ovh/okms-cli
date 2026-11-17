package kmip

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ovh/kmip-go"
	"github.com/ovh/kmip-go/kmipserver"
	"github.com/ovh/kmip-go/payloads"
	"github.com/stretchr/testify/require"
)

// newAesGcmServer returns a KMIP server encrypting and decrypting with an AES-256 GCM key,
// and records the parameters of the last request.
func newAesGcmServer(t *testing.T, lastParams **kmip.CryptographicParameters) *kmipserver.BatchExecutor {
	t.Helper()
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)
	block, err := aes.NewCipher(key)
	require.NoError(t, err)
	gcm, err := cipher.NewGCM(block)
	require.NoError(t, err)

	exec := kmipserver.NewBatchExecutor()
	exec.Route(kmip.OperationEncrypt, kmipserver.HandleFunc(func(_ context.Context, req *payloads.EncryptRequestPayload) (*payloads.EncryptResponsePayload, error) {
		*lastParams = req.CryptographicParameters
		iv := req.IVCounterNonce
		if iv == nil {
			if req.CryptographicParameters == nil || req.CryptographicParameters.RandomIV == nil || !*req.CryptographicParameters.RandomIV {
				return nil, kmipserver.Errorf(kmip.ResultReasonCryptographicFailure, "Required IV is missing.")
			}
			iv = make([]byte, gcm.NonceSize())
			_, _ = rand.Read(iv)
		}
		sealed := gcm.Seal(nil, iv, req.Data, req.AuthenticatedEncryptionAdditionalData)
		tagStart := len(sealed) - gcm.Overhead()
		return &payloads.EncryptResponsePayload{
			UniqueIdentifier:           req.UniqueIdentifier,
			Data:                       sealed[:tagStart],
			IVCounterNonce:             iv,
			AuthenticatedEncryptionTag: sealed[tagStart:],
		}, nil
	}))
	exec.Route(kmip.OperationDecrypt, kmipserver.HandleFunc(func(_ context.Context, req *payloads.DecryptRequestPayload) (*payloads.DecryptResponsePayload, error) {
		*lastParams = req.CryptographicParameters
		sealed := append(append([]byte{}, req.Data...), req.AuthenticatedEncryptionTag...)
		plain, err := gcm.Open(nil, req.IVCounterNonce, sealed, req.AuthenticatedEncryptionAdditionalData)
		if err != nil {
			return nil, kmipserver.Errorf(kmip.ResultReasonCryptographicFailure, "%s", err)
		}
		return &payloads.DecryptResponsePayload{UniqueIdentifier: req.UniqueIdentifier, Data: plain}, nil
	}))
	return exec
}

func TestEncryptDecryptGCM(t *testing.T) {
	var params *kmip.CryptographicParameters
	exec := newAesGcmServer(t, &params)

	out := runCommand(t, exec, "json", encryptCommand(), "encrypt", testObjectId, "hello world", "--mode", "GCM", "--tag-length", "16", "--aad", "my context")
	require.NotNil(t, params)
	require.Equal(t, kmip.BlockCipherModeGCM, params.BlockCipherMode)
	require.Equal(t, int32(16), params.TagLength)
	require.NotNil(t, params.RandomIV)
	require.True(t, *params.RandomIV)

	var encrypted payloads.EncryptResponsePayload
	require.NoError(t, json.Unmarshal([]byte(out), &encrypted))
	require.Equal(t, testObjectId, encrypted.UniqueIdentifier)
	require.Len(t, encrypted.IVCounterNonce, 12)
	require.Len(t, encrypted.AuthenticatedEncryptionTag, 16)

	b64 := base64.StdEncoding.EncodeToString
	out = runCommand(t, exec, "text", decryptCommand(), "decrypt", testObjectId, b64(encrypted.Data),
		"--mode", "GCM", "--iv", b64(encrypted.IVCounterNonce), "--tag", b64(encrypted.AuthenticatedEncryptionTag), "--aad", "my context")
	require.Equal(t, "hello world", out)
	require.NotNil(t, params)
	require.Equal(t, kmip.BlockCipherModeGCM, params.BlockCipherMode)
	require.Nil(t, params.RandomIV, "RandomIV is only for encryption")
}

func TestEncryptTextOutputWithGivenIV(t *testing.T) {
	var params *kmip.CryptographicParameters
	exec := newAesGcmServer(t, &params)
	iv := base64.StdEncoding.EncodeToString([]byte("0123456789ab"))

	out := runCommand(t, exec, "text", encryptCommand(), "encrypt", testObjectId, "hello world", "--mode", "GCM", "--iv", iv)
	require.NotNil(t, params)
	require.Nil(t, params.RandomIV, "no RandomIV when the IV is given")

	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Len(t, lines, 3)
	require.True(t, strings.HasPrefix(lines[0], "Data: "))
	require.Equal(t, "IV: "+iv, lines[1])
	require.True(t, strings.HasPrefix(lines[2], "Tag: "))
}

func TestEncryptWithoutParametersUsesKeyParameters(t *testing.T) {
	exec := kmipserver.NewBatchExecutor()
	var req *payloads.EncryptRequestPayload
	route(exec, kmip.OperationEncrypt, func(r *payloads.EncryptRequestPayload) *payloads.EncryptResponsePayload {
		req = r
		return &payloads.EncryptResponsePayload{UniqueIdentifier: r.UniqueIdentifier, Data: []byte("encrypted")}
	})

	out := runCommand(t, exec, "text", encryptCommand(), "encrypt", testObjectId, "hello world")
	require.Nil(t, req.CryptographicParameters)
	require.Equal(t, []byte("hello world"), req.Data)
	require.Equal(t, "Data: "+base64.StdEncoding.EncodeToString([]byte("encrypted"))+"\n", out)
}

func TestEncryptDecryptRsaOaepParameters(t *testing.T) {
	exec := kmipserver.NewBatchExecutor()
	var encParams, decParams *kmip.CryptographicParameters
	route(exec, kmip.OperationEncrypt, func(r *payloads.EncryptRequestPayload) *payloads.EncryptResponsePayload {
		encParams = r.CryptographicParameters
		return &payloads.EncryptResponsePayload{UniqueIdentifier: r.UniqueIdentifier, Data: []byte("encrypted")}
	})
	route(exec, kmip.OperationDecrypt, func(r *payloads.DecryptRequestPayload) *payloads.DecryptResponsePayload {
		decParams = r.CryptographicParameters
		return &payloads.DecryptResponsePayload{UniqueIdentifier: r.UniqueIdentifier, Data: []byte("hello world")}
	})

	runCommand(t, exec, "json", encryptCommand(), "encrypt", testObjectId, "hello world", "--padding", "OAEP", "--hash", "SHA256")
	out := runCommand(t, exec, "text", decryptCommand(), "decrypt", testObjectId, base64.StdEncoding.EncodeToString([]byte("encrypted")), "--padding", "OAEP", "--hash", "SHA256")
	require.Equal(t, "hello world", out)

	for _, params := range []*kmip.CryptographicParameters{encParams, decParams} {
		require.NotNil(t, params)
		require.Equal(t, kmip.PaddingMethodOAEP, params.PaddingMethod)
		require.Equal(t, kmip.HashingAlgorithmSHA_256, params.HashingAlgorithm)
		require.Equal(t, kmip.MaskGeneratorMGF1, params.MaskGenerator)
		require.Equal(t, kmip.HashingAlgorithmSHA_256, params.MaskGeneratorHashingAlgorithm)
		require.Equal(t, kmip.BlockCipherMode(0), params.BlockCipherMode)
		require.Nil(t, params.RandomIV)
	}
}

// newEcdsaServer returns a KMIP server signing and verifying with an ECDSA P-256 key and ECDSA with SHA-256.
func newEcdsaServer(t *testing.T, lastSign **payloads.SignRequestPayload) *kmipserver.BatchExecutor {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	digest := func(params *kmip.CryptographicParameters, data, digested []byte) ([]byte, error) {
		if params == nil || params.DigitalSignatureAlgorithm != kmip.DigitalSignatureAlgorithmECDSAWithSHA256 || params.HashingAlgorithm != kmip.HashingAlgorithmSHA_256 {
			return nil, kmipserver.Errorf(kmip.ResultReasonPermissionDenied, "Object specified is not able to be used for signing")
		}
		if len(digested) > 0 {
			return digested, nil
		}
		h := sha256.Sum256(data)
		return h[:], nil
	}

	exec := kmipserver.NewBatchExecutor()
	exec.Route(kmip.OperationSign, kmipserver.HandleFunc(func(_ context.Context, req *payloads.SignRequestPayload) (*payloads.SignResponsePayload, error) {
		*lastSign = req
		h, err := digest(req.CryptographicParameters, req.Data, req.DigestedData)
		if err != nil {
			return nil, err
		}
		sig, err := ecdsa.SignASN1(rand.Reader, key, h)
		if err != nil {
			return nil, err
		}
		return &payloads.SignResponsePayload{UniqueIdentifier: req.UniqueIdentifier, SignatureData: sig}, nil
	}))
	exec.Route(kmip.OperationSignatureVerify, kmipserver.HandleFunc(func(_ context.Context, req *payloads.SignatureVerifyRequestPayload) (*payloads.SignatureVerifyResponsePayload, error) {
		h, err := digest(req.CryptographicParameters, req.Data, req.DigestedData)
		if err != nil {
			return nil, err
		}
		validity := kmip.ValidityIndicatorInvalid
		if ecdsa.VerifyASN1(&key.PublicKey, h, req.SignatureData) {
			validity = kmip.ValidityIndicatorValid
		}
		return &payloads.SignatureVerifyResponsePayload{UniqueIdentifier: req.UniqueIdentifier, ValidityIndicator: validity}, nil
	}))
	return exec
}

func TestSignVerify(t *testing.T) {
	var signReq *payloads.SignRequestPayload
	exec := newEcdsaServer(t, &signReq)

	signature := strings.TrimSpace(runCommand(t, exec, "text", signCommand(), "sign", testObjectId, "hello world", "--alg", "ES256"))
	require.Equal(t, []byte("hello world"), signReq.Data)
	require.Empty(t, signReq.DigestedData)
	_, err := base64.StdEncoding.DecodeString(signature)
	require.NoError(t, err)

	out := runCommand(t, exec, "text", verifyCommand(), "verify", testObjectId, "hello world", signature, "--alg", "es256")
	require.Equal(t, "Signature is valid\n", out)

	out = runCommand(t, exec, "json", verifyCommand(), "signature-verify", testObjectId, "hello world", signature, "-a", "ES256")
	var resp payloads.SignatureVerifyResponsePayload
	require.NoError(t, json.Unmarshal([]byte(out), &resp))
	require.Equal(t, kmip.ValidityIndicatorValid, resp.ValidityIndicator)
}

func TestSignVerifyDigested(t *testing.T) {
	var signReq *payloads.SignRequestPayload
	exec := newEcdsaServer(t, &signReq)
	h := sha256.Sum256([]byte("hello world"))
	digest := base64.StdEncoding.EncodeToString(h[:])

	out := runCommand(t, exec, "json", signCommand(), "sign", testObjectId, digest, "--alg", "ES256", "--digested")
	require.Equal(t, h[:], signReq.DigestedData)
	require.Empty(t, signReq.Data)
	var signed payloads.SignResponsePayload
	require.NoError(t, json.Unmarshal([]byte(out), &signed))

	// a signature of the digest is the signature of the data
	out = runCommand(t, exec, "text", verifyCommand(), "verify", testObjectId, "hello world", base64.StdEncoding.EncodeToString(signed.SignatureData), "--alg", "ES256")
	require.Equal(t, "Signature is valid\n", out)
}

// TestVerifyInvalidSignature runs in a subprocess, since a failing command exits the process.
func TestVerifyInvalidSignature(t *testing.T) {
	if os.Getenv("OKMS_TEST_SUBPROCESS") == "1" {
		var signReq *payloads.SignRequestPayload
		exec := newEcdsaServer(t, &signReq)
		signature := strings.TrimSpace(runCommand(t, exec, "text", signCommand(), "sign", testObjectId, "hello world", "--alg", "ES256"))
		runCommand(t, exec, "text", verifyCommand(), "verify", testObjectId, "tampered", signature, "--alg", "ES256")
		return
	}
	stderr := requireExitFailure(t, "TestVerifyInvalidSignature")
	require.Contains(t, stderr, "Error: Signature is invalid")
}

func TestSignatureAlgorithmParameters(t *testing.T) {
	for alg, expected := range map[string]kmip.CryptographicParameters{
		"RS256": {DigitalSignatureAlgorithm: kmip.DigitalSignatureAlgorithmSHA_256WithRSAEncryption, HashingAlgorithm: kmip.HashingAlgorithmSHA_256},
		"RS512": {DigitalSignatureAlgorithm: kmip.DigitalSignatureAlgorithmSHA_512WithRSAEncryption, HashingAlgorithm: kmip.HashingAlgorithmSHA_512},
		"PS384": {DigitalSignatureAlgorithm: kmip.DigitalSignatureAlgorithmRSASSA_PSS, HashingAlgorithm: kmip.HashingAlgorithmSHA_384},
		"ES512": {DigitalSignatureAlgorithm: kmip.DigitalSignatureAlgorithmECDSAWithSHA512, HashingAlgorithm: kmip.HashingAlgorithmSHA_512},
	} {
		exec := kmipserver.NewBatchExecutor()
		var params *kmip.CryptographicParameters
		route(exec, kmip.OperationSign, func(r *payloads.SignRequestPayload) *payloads.SignResponsePayload {
			params = r.CryptographicParameters
			return &payloads.SignResponsePayload{UniqueIdentifier: r.UniqueIdentifier, SignatureData: []byte("signature")}
		})
		runCommand(t, exec, "text", signCommand(), "sign", testObjectId, "hello world", "--alg", alg)
		require.NotNil(t, params, alg)
		require.Equal(t, expected, *params, alg)
	}
}

func TestCreateWithUsageLimits(t *testing.T) {
	exec := kmipserver.NewBatchExecutor()
	var symmetric *payloads.CreateRequestPayload
	var keyPair *payloads.CreateKeyPairRequestPayload
	route(exec, kmip.OperationCreate, func(r *payloads.CreateRequestPayload) *payloads.CreateResponsePayload {
		symmetric = r
		return &payloads.CreateResponsePayload{ObjectType: kmip.ObjectTypeSymmetricKey, UniqueIdentifier: testObjectId}
	})
	route(exec, kmip.OperationCreateKeyPair, func(r *payloads.CreateKeyPairRequestPayload) *payloads.CreateKeyPairResponsePayload {
		keyPair = r
		return &payloads.CreateKeyPairResponsePayload{PrivateKeyUniqueIdentifier: testObjectId, PublicKeyUniqueIdentifier: testObjectId}
	})
	usageLimits := func(attrs []kmip.Attribute) any {
		for _, attr := range attrs {
			if attr.AttributeName == kmip.AttributeNameUsageLimits {
				return attr.AttributeValue
			}
		}
		return nil
	}

	runCommand(t, exec, "json", createCommand(), "create", "symmetric", "--alg", "AES", "--size", "256", "--usage-limit", "10")
	ten := int64(10)
	require.Equal(t, kmip.UsageLimits{UsageLimitsTotal: 10, UsageLimitsCount: &ten, UsageLimitsUnit: kmip.UsageLimitsUnitObject},
		usageLimits(symmetric.TemplateAttribute.Attribute))

	runCommand(t, exec, "json", createCommand(), "create", "key-pair", "--alg", "ECDSA", "--curve", "P-256", "--usage-limit", "4096", "--usage-limit-unit", "byte")
	require.NotNil(t, keyPair.CommonTemplateAttribute)
	limit := int64(4096)
	require.Equal(t, kmip.UsageLimits{UsageLimitsTotal: 4096, UsageLimitsCount: &limit, UsageLimitsUnit: kmip.UsageLimitsUnitByte},
		usageLimits(keyPair.CommonTemplateAttribute.Attribute))

	runCommand(t, exec, "json", createCommand(), "create", "symmetric", "--alg", "AES", "--size", "256")
	require.Nil(t, usageLimits(symmetric.TemplateAttribute.Attribute))
}

// requireExitFailure re-runs the given test in a subprocess with OKMS_TEST_SUBPROCESS=1, checks
// that it exits with code 1, and returns its stderr.
func requireExitFailure(t *testing.T, test string) string {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^"+test+"$") //nolint:gosec // re-runs this test binary
	cmd.Env = append(os.Environ(), "OKMS_TEST_SUBPROCESS=1")
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()
	var exitErr *exec.ExitError
	require.True(t, errors.As(err, &exitErr), "expected the subprocess to fail, got %v", err)
	require.Equal(t, 1, exitErr.ExitCode())
	return stderr.String()
}
