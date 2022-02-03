package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "fmt"
    "math/big"
    "net"
    "os"
    "time"
)

//生成RSA私钥和公钥，保存到文件中
// bits 证书大小
func GenerateRSAKey(bits int) {
    //GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
    //Reader是一个全局、共享的密码用强随机数生成器
    privateKey, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        panic(err)
    }

    //保存私钥
    //通过x509标准将得到的ras私钥序列化为ASN.1的DER编码字节数组
    X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)

    //构建一个pem.Block结构体对象
    privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}

    //使用pem格式对x509输出的内容进行编码
    //创建文件保存私钥
    privateFile, err := os.Create("private.pem")
    if err != nil {
        panic(err)
    }
    defer privateFile.Close()
    //将数据保存到文件
    pem.Encode(privateFile, &privateBlock)

    //保存公钥
    //获取公钥的数据
    publicKey := privateKey.PublicKey
    //X509对公钥编码
    X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
    if err != nil {
        panic(err)
    }
    //pem格式编码
    //创建用于保存公钥的文件
    publicFile, err := os.Create("public.pem")
    if err != nil {
        panic(err)
    }
    defer publicFile.Close()
    //创建一个pem.Block结构体对象
    publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
    //保存到文件
    pem.Encode(publicFile, &publicBlock)
}

//从私钥文件（pem格式）读取为*rsa.PrivateKey，和rsa.GenerateKey生成的结构一致，path是私钥文件路径
func RSAPrivateKeyfromPemFile(path string) *rsa.PrivateKey {
    //打开path指向的文件为*os.File
    file, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    //获取文件内容
    //获取长度
    info, _ := file.Stat()
    //创建[]byte缓存，长度和文件一致
    buf := make([]byte, info.Size())
    //文件读入缓存
    file.Read(buf)

    //pem解码
    block, _ := pem.Decode(buf)

    //X509解码
    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        panic(err)
    }
    return privateKey
}

//从公钥文件（pem格式）读取为*rsa.PublicKey
func RSAPublicKeyFromPemFile(path string) *rsa.PublicKey {
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
    return publicKey
}

//sha256哈希字符串
func Sha256Sum(plainText string) []byte {
    // 第一种调用方法
    // sum := sha256.Sum256([]byte(plainText))

    // 第二种调用方法，支持一个大文件多次写入和累计hash
    // h := sha256.New()
    // h.Write([]byte(plainText))
    // sum := h.Sum(nil)

    // 第三种调用方法
    h := sha256.New()
    sum := h.Sum([]byte(plainText))

    return sum
}

//RSA加密
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
    //对明文进行加密，用pkcs1v15的包
    cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
    //只用rsa包函数加密，rsa包里加密函数只有一种，解密有两种
    //cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plainText, nil)
    if err != nil {
        panic(err)
    }
    //返回密文
    return cipherText
}

//RSA解密
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
    // 只用rsa包的函数，rsa包里加密函数只有一种，解密有两种
    // plainText, err := privateKey.Decrypt(nil, cipherText, &rsa.OAEPOptions{Hash: crypto.SHA256})
    // plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cipherText, nil)
    if err != nil {
        panic(err)
    }
    //返回明文
    return plainText
}

func GenerateCertificate() {
    //大数结构
    max := new(big.Int).Lsh(big.NewInt(1), 128)
    //随机数存入大数结构
    serialNumber, _ := rand.Int(rand.Reader, max)
    //证书的subject
    subject := pkix.Name{
        Country:            []string{"CN"},
        Province:           []string{"BeiJing"},
        Organization:       []string{"Devops"},
        OrganizationalUnit: []string{"certDevops"},
        CommonName:         "127.0.0.1",
    }

    //证书结构体，包括serialNumber、subject、IPAddress
    template := x509.Certificate{
        SerialNumber: serialNumber,
        Subject:      subject,
        NotBefore:    time.Now(),
        NotAfter:     time.Now().Add(365 * 24 * time.Hour),
        KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
    }

    //生成RSA
    pk, _ := rsa.GenerateKey(rand.Reader, 2048)

    //通过证书结构体、RSA创建密钥
    derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)

    //创建证书文件，并保存证书
    certOut, _ := os.Create("server.pem")
    pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
    certOut.Close()

    //创建密钥文件，并保存密钥
    keyOut, _ := os.Create("server.key")
    pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
    keyOut.Close()
}

func test() {
    //生成密钥对，保存到文件
    GenerateRSAKey(2048)
    //加密
    data := []byte("hello world")
    encrypt := RSA_Encrypt(data, "public.pem")
    fmt.Println(string(encrypt))

    // 解密
    decrypt := RSA_Decrypt(encrypt, "private.pem")
    fmt.Println(string(decrypt))
}
