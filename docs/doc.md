# document
### structure
```
    gochat
        |-common            /* global data structure */
            |-common.go     /* 定义了WsClient、WsMessage以及PubSubMessage
                               定义了消息类别和用户状态
                               定义了全局的房间-用户映射
                               定义了消息队列名和redis键前缀 */
        |-config            /* configuration files */
            |-config.yaml   /* redis和mysql配置 */
        |-controller        /* backend logic */
            |-session       /* session */     
            |-maincontroller.go     /* HandlerFunc
                                       包括注册、登录、登出、搜索房间、私聊 */
            |-wscontroller.go       /* HandlerFunc
                                       负责建立websocket连接 */
        |-model             /* expired */
        |-modelv2           /* database operation */
            |-baseinfo.go   /* 定义用户和房间基本信息以及方法 */
            |-common.go     /* 初始化validator */
            |-message.go    /* 定义聊天信息以及方法 */
            |-privatechat.go /* 获取私聊房间号
                                如果是新创建的私聊就创建新的私聊房间号，这个过程需要获取分布式锁 */
        |-pkg               /* function module */
            |-adapter       /* code reuse based on decorator pattern */
                |-base.go           /* 定义了QueueTask接口
                                       定义了QueueTask状态和描述信息 
                                       定义了Handler接口 */
                |-logger.go         /* 实现了Handler接口，用于统计QueueTask的完成情况 */
                |-queuetask.go      /* 定义了QueueTaskAdapter，装饰一个实现QueueTask接口的实例 */
            |-lock          /* a simple distributed lock */
            |-queue         /* a simple mq base on redis */
                |-base.go           /* 定义了消息
                                       定义了Node接口，包含生产者和消费者接口 */
                |-queue.go  /* 装配器模式，封装实现了Node接口的实例的方法 */
                |-redis.go  /* 定义了RedisNode
                               基于redis的api实现了Node接口，包含发送消息、接收消息、订阅和发布 */
            |-service       /* backend service */
                |-base.go           /* 定义了Service接口，service下的模块都要实现这一接口 */
                |-entry.go          /* 处理用户加入房间、订阅房间号、更新房间-用户映射以及读websocket连接 */
                |-broadcaster.go    /* 广播消息 */
                |-subscriber.go     /* 订阅频道并接收消息 */
            |-task          /* QueueTask implementation */
                |-dispatcher        /* 根据聊天信息类型将消息分发给对应的handler */
                |-contenthandler    /* 处理文本类型聊天信息 */
        |-router            /* restful router */
            |-router.go      
        |-sql               /* DDL */
        |-util              /* utility */
            |-initserver.go         /* 载入配置文件
                                       初始化redis
                                       初始化mysql
                                       初始化mq */
        |-view              /* frontend */
        |-main.go           /* run this project! */
```
### log in and sign up
![image](https://github.com/SkylinerZzz/gochat/blob/main/docs/img/img.png)
### create and search room
![image](https://github.com/SkylinerZzz/gochat/blob/main/docs/img/img_1.png)
### log in and sign up
![image](https://github.com/SkylinerZzz/gochat/blob/main/docs/img/img_2.png)
### group chat
![image](https://github.com/SkylinerZzz/gochat/blob/main/docs/img/img_3.png)
### private chat
![image](https://github.com/SkylinerZzz/gochat/blob/main/docs/img/img_4.png)
![image](https://github.com/SkylinerZzz/gochat/blob/main/docs/img/img_5.png)
