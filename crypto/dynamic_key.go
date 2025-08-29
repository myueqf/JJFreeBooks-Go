package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

type KeyIv struct {
	Key []byte
	Iv  []byte
}

// DynamicDecrypt 动态密钥解密函数，匹配 JavaScript 实现
// 返回解密后的内容、密钥和IV信息
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

	// 生成密钥和IV
	// JavaScript: CryptoJS.MD5(v43 + v38).toString().slice(0, 8)
	md5Hash1 := md5.Sum([]byte(v43 + v38))
	keyStr := hex.EncodeToString(md5Hash1[:])[:8] // 取前8个字符作为UTF8字符串

	// JavaScript: CryptoJS.MD5(v38).toString().slice(0, 8)
	md5Hash2 := md5.Sum([]byte(v38))
	ivStr := hex.EncodeToString(md5Hash2[:])[:8] // 取前8个字符作为UTF8字符串

	// 转换为字节数组（UTF8编码，不是hex解码）
	key := []byte(keyStr)
	iv := []byte(ivStr)

	return dest, KeyIv{
		Key: key,
		Iv:  iv,
	}, nil
}

// DynamicDecryptWithContent 执行完整的动态解密过程
// 包括密钥生成和DES解密
func DynamicDecryptWithContent(context, accessKey, keystring string) (string, error) {
	dest, keyIv, err := DynamicDecrypt(context, accessKey, keystring)
	if err != nil {
		return "", err
	}

	// 使用生成的密钥和IV进行DES解密
	return DesDecrypt(dest, keyIv.Key, keyIv.Iv)
}
