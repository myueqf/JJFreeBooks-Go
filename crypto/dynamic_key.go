package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

type KeyIv struct {
	Key []byte
	Iv  []byte
}

func DynamicDecrypt(context, accessKey, keystring string) (string, KeyIv, error) {
	accessKeyLen := len(accessKey)
	v6 := accessKey[accessKeyLen-1]

	var v9 int
	for i := 0; i < accessKeyLen; i++ {
		v9 += int(accessKey[i])
	}
	v15 := v9 % len(keystring)
	v17 := v9 / 65
	v18 := len(keystring)

	end := v17 + v15
	if end > v18 {
		end = v18
	}
	v43 := keystring[v15:end]

	var v38, dest string
	if (v6 & 1) != 0 {
		v38 = context[len(context)-12:]
		dest = context[:len(context)-12]
	} else {
		v38 = context[:12]
		dest = context[12:]
	}
	// 生成 key = MD5(v43 + v38).hex[0:16] -> 取前 8 字节（16 hex 字符）
	md5Hash1 := md5.Sum([]byte(v43 + v38))
	keyHex := hex.EncodeToString(md5Hash1[:])[:16] // 取前 16 个 hex 字符（8 字节）

	// 生成 iv = MD5(v38).hex[0:16] -> 取前 8 字节
	md5Hash2 := md5.Sum([]byte(v38))
	ivHex := hex.EncodeToString(md5Hash2[:])[:16]

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", KeyIv{}, err
	}
	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return "", KeyIv{}, err
	}
	return dest, KeyIv{
		Key: key,
		Iv:  iv,
	}, err
}
