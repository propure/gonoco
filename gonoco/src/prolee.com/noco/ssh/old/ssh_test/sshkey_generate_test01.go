package gonoco/ssh_test

import (
    "fmt"
    "log"
    "os"

    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"

    "golang.org/x/crypto/ssh"
)

func main() {
    /* 根据随机数和位数，生成私钥 */
    privatekey, err := rsa.GenerateKey(rand.Reader, 1024)
    if err != nil {
        log.Fatalln(err)
        os.Exit(1)
    }
    /* 从私钥里提取公钥 */
    publickey := &privatekey.PublicKey

    /* 把私钥对象转为[]byte */
    pkey := pem.EncodeToMemory(&pem.Block{
        Bytes: x509.MarshalPKCS1PrivateKey(privatekey),
        Type:  "RSA PRIVATE KEY",
    })

    /* 把公钥对象转为[]byte，需要两步 */
    /* 第一步是转为der格式的[]Byyte */
    tmpBytesPublic, err := x509.MarshalPKIXPublicKey(publickey)
    if err != nil {
        log.Fatalln(err)
        os.Exit(1)
    }
    /* 第二步是先构造一个pem.Block结构，再用pem.EncodeToMemory转为pem的byte[]*/
    pub := pem.EncodeToMemory(&pem.Block{
        Bytes: tmpBytesPublic,
        Type:  "PUBLIC KEY",
    })

    /* 根据rsa公钥对象生成ssh公钥对象*/
    sshpublickey, err := ssh.NewPublicKey(publickey)
    if err != nil {
        log.Fatalln(err)
        os.Exit(1)
    }
    /* ssh公钥对象转为[]byte形式 */
    sshpub := ssh.MarshalAuthorizedKey(sshpublickey)

    /* 返回私钥字符串，公钥字符串 */
    fmt.Printf("%s\n", pub)
    fmt.Printf("%s\n", pkey)
    fmt.Printf("%s\n", sshpub)

}

//RSA公钥加密
// plainText 要加密的数据
// path 公钥匙文件地址
func RSA_Encrypt(plainText []byte, path string) []byte {
    //打开文件
    file, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    //读取文件的内容
    info, _ := file.Stat()
    buf := make([]byte, info.Size())
    file.Read(buf)
    //pem解码
    block, _ := pem.Decode(buf)
    //x509解码

    publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        panic(err)
    }
    //类型断言
    publicKey := publicKeyInterface.(*rsa.PublicKey)
    //对明文进行加密
    cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
    if err != nil {
        panic(err)
    }
    //返回密文
    return cipherText
}

//RSA私钥解密
// cipherText 需要解密的byte数据
// path 私钥文件路径
func RSA_Decrypt(cipherText []byte, path string) []byte {
    //打开文件
    file, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer file.Close()
    //获取文件内容
    info, _ := file.Stat()
    buf := make([]byte, info.Size())
    file.Read(buf)
    //pem解码
    block, _ := pem.Decode(buf)
    //X509解码
    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        panic(err)
    }
    //对密文进行解密
    plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
    //返回明文
    return plainText
}
