package crypto

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
)

// pkcs5Padding 对数据进行 PKCS5 填充。
// PKCS5 实质上是块大小为 8 字节的 PKCS7 填充（DES 的块大小为 8）。
// 如果原始数据长度不是 8 的倍数，则用缺少的字节数进行填充。
//
// 参数：
//   - data: 原始数据（字节切片）
//
// 返回值：
//   - 填充后的数据（长度为 8 的倍数）
func pkcs5Padding(data []byte) []byte {
	padding := 8 - len(data)%8 // 需要填充的字节数
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding) // 所有填充字节都等于填充长度
	}
	return append(data, padText...) // 返回原始数据 + 填充
}

// pkcs5Unpadding 去除 PKCS5 填充。
// 会验证填充是否合法，防止恶意输入或解密错误。
//
// 参数：
//   - data: 已解密的数据（包含填充）
//
// 返回值：
//   - 去除填充后的原始数据，或错误信息
func pkcs5UnPadding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("数据为空，无法去填充")
	}
	last := data[len(data)-1] // 最后一个字节表示填充长度
	padding := int(last)

	// 检查填充长度是否合法（必须在 1~8 范围内）
	if padding == 0 || padding > 8 {
		return nil, fmt.Errorf("无效的 PKCS5 填充值：%d", padding)
	}
	if len(data) < padding {
		return nil, fmt.Errorf("填充长度超过数据总长度")
	}

	// 验证所有填充字节是否都等于 padding 值
	for i := 0; i < padding; i++ {
		if data[len(data)-padding+i] != last {
			return nil, fmt.Errorf("PKCS5 填充不一致，可能解密失败")
		}
	}

	// 去除填充部分
	return data[:len(data)-padding], nil
}

// DesEncrypt 使用 DES 算法加密明文数据（CBC 模式 + PKCS5 填充）。
// 加密结果会进行 Base64 编码，便于传输和存储。
//
// 参数：
//   - plainText: 明文数据（字节切片）
//   - key: 密钥，必须是 8 字节长
//   - iv: 初始化向量，必须是 8 字节长
//
// 返回值：
//   - Base64 编码的密文字符串，以及可能的错误
func DesEncrypt(plainText, key, iv []byte) (string, error) {
	if len(key) != 8 {
		return "", fmt.Errorf("密钥长度必须为 8 字节，当前为 %d 字节", len(key))
	}
	if len(iv) != 8 {
		return "", fmt.Errorf("IV 长度必须为 8 字节，当前为 %d 字节", len(iv))
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建 DES 加密器失败：%w", err)
	}

	plainText = pkcs5Padding(plainText)

	cipherText := make([]byte, len(plainText))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, plainText) // 执行加密

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DesDecrypt 使用 DES 算法解密密文（Base64 编码 + CBC 模式 + PKCS5 去填充）。
//
// 参数：
//   - ciphertextStr: 经过 Base64 编码的密文字符串
//   - key: 解密密钥，必须是 8 字节长
//   - iv: 初始化向量，必须是 8 字节长
//
// 返回值：
//   - 解密后的明文字符串，以及可能的错误
func DesDecrypt(ciphertextStr string, key, iv []byte) (string, error) {
	if len(key) != 8 {
		return "", fmt.Errorf("密钥长度必须为 8 字节，当前为 %d 字节", len(key))
	}
	if len(iv) != 8 {
		return "", fmt.Errorf("IV 长度必须为 8 字节，当前为 %d 字节", len(iv))
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextStr)
	if err != nil {
		return "", fmt.Errorf("base64 解码失败：%w", err)
	}

	if len(ciphertext)%8 != 0 {
		return "", fmt.Errorf("密文长度不是 8 的倍数，可能已损坏")
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建 DES 解密器失败：%w", err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	unpaddedData, err := pkcs5UnPadding(decrypted)
	if err != nil {
		return "", fmt.Errorf("PKCS5 去填充失败：%w", err)
	}

	return string(unpaddedData), nil
}
