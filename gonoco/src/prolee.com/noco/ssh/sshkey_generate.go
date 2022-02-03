package ssh

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "os"
    "unsafe"

    "golang.org/x/crypto/ssh"
)

/* 根据位数，生成私钥/公钥对 */
func RSAGenerateKey(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        return nil, nil, err
    }
    return privateKey, &privateKey.PublicKey, nil

}

/* 把私钥转换为pem的[]byte格式 */
func RSAPrivateKeyEncodePemBytes(private *rsa.PrivateKey) []byte {
    privatePem := pem.EncodeToMemory(&pem.Block{
        Type: "RSA PRIVATE KEY",
        //Headers: map[string]string{}, //可选
        Bytes: x509.MarshalPKCS1PrivateKey(private),
    })
    return privatePem
}

/* 把私钥转换为pem的[]byte加密格式 */
func RSAPrivateKeyEncodeSecretPemBytes(private *rsa.PrivateKey, password string) ([]byte, error) {
    privateBlock, err := x509.EncryptPEMBlock(rand.Reader, "RSA Private Key", x509.MarshalPKCS1PrivateKey(private), []byte(password), x509.PEMCipherAES256)
    if err != nil {
        return nil, err
    }
    privatePem := pem.EncodeToMemory(privateBlock)

    return privatePem, nil

}

/* 把公钥转换为rsa pem的[]byte格式 */
func RSAPublicKeyEncodePemByte(public *rsa.PublicKey) []byte {
    publicBytes, err := x509.MarshalPKIXPublicKey(public)
    if err != nil {
        return nil
    }
    return pem.EncodeToMemory(&pem.Block{
        Type: "PUBLIC KEY",
        //Headers: map[string]string{}, //可选
        Bytes: publicBytes,
    })
}

func byte2string(b []byte) *string {
    return (*string)(unsafe.Pointer(&b))

}

/* 把公钥转换为ssh pem的[]byte格式，这个格式和普通的rsa有区别 */
func EncodeSSHKey(public *rsa.PublicKey) ([]byte, error) {
    sshPublicKey, err := ssh.NewPublicKey(public)
    if err != nil {
        return nil, err
    }
    return ssh.MarshalAuthorizedKey(sshPublicKey), nil
}

func GenerateSSHPemKey(bits int) (string, string, error) {
    privatekey, publickey, err := RSAGenerateKey(1024)
    if err != nil {
        os.Exit(1)
    }

    sshpublic, err := EncodeSSHKey(publickey)
    if err != nil {
        os.Exit(1)
    }

    return string(RSAPrivateKeyEncodePemBytes(privatekey)), string(sshpublic), nil
}

func GenerateRsaPemKey(bits int) (string, string, error) {
    privatekey, publickey, err := RSAGenerateKey(bits)
    if err != nil {
        return "", "", err
    }

    return string(RSAPrivateKeyEncodePemBytes(privatekey)), string(RSAPublicKeyFromPemString(publickey)), nil
}

//没有使用string做为参数或返回值是因为go没有提供&string(plainText)这种操作
func RSAEncrypt(plainText []byte, publicKey *rsa.PublicKey) ([]byte, error) {
    //对明文进行加密
    cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
    if err != nil {
        return nil, err
    }
    //返回密文
    return cipherText, nil
}

func RSADecrypt(cipherText []byte, privateKey *rsa.PrivateKey) ([]byte, error) {

    //对密文进行解密
    plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
    if err != nil {
        return nil, err
    }
    //返回密文
    return plainText, nil
}

func RSAPublicKeyFromPemString(buf []byte) (*rsa.PublicKey, error) {
    //pem解码，返回*pem.Block类型，也就是der格式的
    block, _ := pem.Decode(buf)
    //x509解码，返回的是interface{}
    publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }

    //类型断言，转换
    publicKey := publicKeyInterface.(*rsa.PublicKey)
    return publicKey, nil
}

func RSAPrivateKeyFromPemString(buf []byte) (*rsa.PrivateKey, error) {
    //pem解码，返回*pem.Block类型，也就是der格式的
    block, _ := pem.Decode(buf)
    //x509解码，直接返回*rsa.PrivateKey
    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }

    return privateKey, nil
}

// 如果有密码
func RSAPrivateKeyFromSecretPemString(buf []byte, password string) (*rsa.PrivateKey, error) {
    //pem解码，返回*pem.Block类型，也就是der格式的
    block, _ := pem.Decode(buf)

    privateDerByte, err := x509.DecryptPEMBlock(block, []byte(password))
    //x509解码，直接返回*rsa.PrivateKey
    privateKey, err := x509.ParsePKCS1PrivateKey(privateDerByte)
    if err != nil {
        return nil, err
    }

    return privateKey, nil
}

func RSAPublicKeyFromPemFile(path string) (*rsa.PublicKey, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    //读取文件的state，获取文件大小
    info, _ := file.Stat()

    //根据文件大小创建buf
    buf := make([]byte, info.Size())
    file.Read(buf)

    //pem解码，返回*pem.Block类型
    block, _ := pem.Decode(buf)
    //x509解码

    publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }

    //类型断言
    publicKey := publicKeyInterface.(*rsa.PublicKey)

    return publicKey, nil
}
