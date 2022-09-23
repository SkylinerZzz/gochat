# gochat
A simple chat room with gin + gorm + websocket
### structure
```
    gochat
        |-controller        /* background logic */
            |-user.go       // user register and login
            |-home.go       // list all rooms and create new room
            |-room.go       // enter the chatroom
            |-server.go     // upgrade http to websocket, handle messages from connection and broadcast them
        |-model             /* database and cache */
            |-mysql.go      // connect to mysql
            |-user.go       // handle user info
            |-room.go       // handle room info
            |-message.go    // handle message info
            |-redis.go      //connect to redis
            |-cache.go      // user cache and room cache
            |-cache_test.go // test cache
        |-router            /* restful router */
            |-router.go     
        |-session           /* session */
            |-session.go    
        |-sql //DDL
        |-static            /* icon and js, craete websocket connection */
        |-view              /* foreground */
        |-main.go           // run this project!
```
### package level variables
+ controller/server.go
    ```
    type Client struct {
        Conn     *websocket.Conn
        Username string
        RoomId   string
    }
    type Message struct {
        MsgType int         `json:"msgType"`
        Data    interface{} `json:"data"`
    }

    var (
        mutex   = sync.Mutex{}                     
        once    = sync.RWMutex{}
        rooms   = make(map[string][]Client)
        users   = make(map[string]map[string]bool)
        enter   = make(chan Client, 10)
        leave   = make(chan Client, 10)
        message = make(chan Message, 100)
    )
    ```
    Client是客户端信息，包含一个websocket连接和用户名以及房间号；Message是客户端的消息体，包含消息类型以及消息正文。  
    rooms的键是房间号，值是一个Client切片。代表有房间里有哪些客户端接入。  
    users是一个嵌套map，代表指定房间（号）里有哪些用户（名）。  
    enter是一个通道，如果websocket上接收的消息表明用户进入房间，则构造对应的客户端连接发送至此通道，并构造对应的Message发送至message通道。  
    leave是一个通道，如果websocket上接收的消息表明用户退出房间，则构造对应的客户端连接发送至此通道，并构造对应的Message发送至message通道。  
    message是一个通道，存放等待处理的Message。  
    用户加入房间后服务器开启两个线程read和write。read负责从connection中读取消息并存储聊天记录。write负责从通道中接收消息，维护users用户映射和rooms客户端连接以及处理各类消息。  
+ model/mysql.go  
    ```
    var ChatDB *gorm.DB
    ```
    mysql句柄  
+ model/user.go
    ```
    type User struct {
        gorm.Model
        Username string
        Password string
    }
    ```
    用户对象
+ model/room.go
    ```
    type Room struct {
        gorm.Model
        UserId   uint
        RoomName string
    }
    ```
    房间对象
+ model/message.go
    ```
    type Message struct {
        gorm.Model
        UserId  uint
        RoomId  uint
        Content string
    }
    ```
    聊天记录对象
+ model/redis.go
    ```
    var ChatCache redis.Conn
    ```
    redis句柄
