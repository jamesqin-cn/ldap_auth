# ldap_auth
将ldap认证能力封装成微服务，以便身份验证服务调用

## 启动
首先要配置LDAP服务器相关信息，示例配置文件如下：
```
ldap:
    addr: "10.255.8.254:389"
    base_dn: "OU=company,DC=example,DC=net"
    filter:  "(sAMAccountName=%s)"
    attributes:
      - "sAMAccountName"
      - "name"
      - "mail"
      - "telephoneNumber"
      - "memberOf"
      - "logonCount"
      - "userAccountControl"
    bind_dn: "CN=administrator,CN=Users,DC=example,DC=net"
    bind_passwd: "example&2018"

http:
    listen: "0.0.0.0:9066"
```

将配置文件名传参到可执行文件便可
```
./app -config ./conf/config.yml
```

## 调用
程序启动后，内置http服务器便会侦听在配置文件指定端口，提供2个查询服务

### 验证账号密码
- 入口
    /api/auth
- 方法
    POST
- 请求参数
请求参数以json表示，携带 username 和 password 成员

- 调用方法示例
```
curl -X POST 127.0.0.1:9066/api/auth -d '{"username":"meeting","password":"abcd.1234"}'
```

- 返回结果示例
```
{
  "err_code": 0,
  "err_msg": "ok",
  "data": {
    "DN": "CN=会议,OU=服务依赖,OU=技术线,OU=company,DC=example,DC=net",
    "Attributes": {
      "logonCount": [
        "719"
      ],
      "name": [
        "会议"
      ],
      "sAMAccountName": [
        "meeting"
      ],
      "userAccountControl": [
        "66048"
      ]
    }
  }
}
```

### 查询账号列表
- 入口
    /api/list
- 方法
    POST
- 请求参数
请求参数以json表示，设置空的json对象即可，即：{}

- 调用方法示例
```
curl -X POST 127.0.0.1:9066/api/list -d '{}'
```

- 返回结果示例
```
{
  "err_code": 0,
  "err_msg": "ok",
  "count": 1,
  "data": [{
    "DN": "CN=会议,OU=服务依赖,OU=技术线,OU=company,DC=example,DC=net",
    "Attributes": {
      "logonCount": [
        "719"
      ],
      "name": [
        "会议"
      ],
      "sAMAccountName": [
        "meeting"
      ],
      "userAccountControl": [
        "66048"
      ]
    }
  },{
      <next record>
  },{
      <next record>
  }]
```
