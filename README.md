# micro-mall-users-consumer

#### 介绍
用户系统事件消费者，针对用户相关事件消费处理

#### 软件架构
queue

#### 框架，库依赖
kelvins框架支持（gRPC，cron，queue，web支持）：https://gitee.com/kelvins-io/kelvins   
g2cache缓存库支持（两级缓存）：https://gitee.com/kelvins-io/g2cache   

#### 安装教程

1.仅构建  sh build.sh   
2 运行  sh build-run.sh    
3 停止 sh stop.sh

#### 使用说明
配置参考
```toml
[kelvins-server]
IsRecordCallResponse = true
Environment = "dev"

[kelvins-logger]
RootPath = "./logs"
Level = "debug"

[kelvins-mysql]
Host = "127.0.0.1:3306"
UserName = "root"
Password = "fasdfa"
DBName = "micro_mall_user"
Charset = "utf8"
PoolNum =  10
MaxIdleConns = 5
ConnMaxLifeSecond = 3600
MultiStatements = true
ParseTime = true

[kelvins-redis]
Host = "127.0.0.1:6379"
Password = "uytrutyu"
DB = 1
PoolNum = 10

[kelvins-queue-server]
CustomQueueList = "user_register_notice,user_state_notice"
WorkerConcurrency = 1

[kelvins-queue-amqp]
Broker = "amqp://micro-mall:szJ9aePR@localhost:5672/micro-mall"
DefaultQueue = "user_register_notice"
ResultBackend = "redis://urtyuryt@127.0.0.1:6379/8"
ResultsExpireIn = 3600
Exchange = "user_register_notice"
ExchangeType = "direct"
BindingKey = "user_register_notice"

[email-config]
User = "urtyu@qq.com"
Password = "urtyiirtirty"
Host = "smtp.qq.com"
Port = "465"
```

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request
