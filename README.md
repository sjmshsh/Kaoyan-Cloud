# 项目架构图

![项目架构图.excalidraw](.\images\项目架构图.excalidraw.png)

# 技术栈介绍

## 开发

- Go
- gin
- gorm
- grpc
- etcd
- MySQL
- Redis
- Rabbitmq

## 部署(部署之后出现了一些问题，目前还没有解决，但是本地环境下所有接口均已测试成功)

- 使用Dockerfile制作镜像
- 使用Kubernetes对容器进行编排

# 数据库制表

kaoyanyun.counter的fuzzy表：

![福昕截屏20230129145444662](.\images\福昕截屏20230129145444662.PNG)

这个是模糊计数表，用于计数服务的模糊计数，例如浏览数，点赞数等等。

kaoyanyun.counter的precise表：

![福昕截屏20230129145444662](.\images\福昕截屏20230129145444662.PNG)

这两个表是一模一样的，为什么要分开呢？

计数服务每隔一段时间会从数据库中同步数据，这会造成很大的性能抖动，我们要想办法解决

这里想到的方案是：把模糊计数和精准计数的表分开。

把模糊计数和精准计数的表分开，因为我们只有模糊计数才需要去进行扫描，批量更新，这样的话可以减少扫描波及的行数，尽可能的规避性能抖动。

kaoyanyun.redpack的redpack表：

这个表用于存放红包数据。

![福昕截屏20230129145813903](.\images\福昕截屏20230129145813903.PNG)

kaoyanyun.user的attention表，用户关注表：

![福昕截屏20230129145915286](.\images\福昕截屏20230129145915286.PNG)

kaoyanyun.user的follower表，用户追随表：

![福昕截屏20230129150116663](.\images\福昕截屏20230129150116663.PNG)

这个表和attention是差不多的，attention有一个flg字段，flg代表大V的标识，这个是为了项目未来针对大V专门进行一些设计，暂时没有用到。而为什么要分成attention和follower呢？其实一张表也可以做。

原因是：

这里有一个问题。table_relation根据follower进行拆分，查询某个用户关注的人很容易，因为相同的followerId的数据一定分布在相同的分片上（我关注了谁）。但是一旦需要查询谁关注了某个用户（谁关注了我），这样查询需要路由到所有分片上进行，因为相同的attentionId的数据分散在不同的分片上，查询效率低。由于查询follower和attention的访问量是相似的，所以无论根据followerId还是attentionId进行拆分，总会有一般的场景查询效率低下。

所以针对上述问题，进行垂直拆分，分为follower表和attention表，分表记录某个用户的关注者和被关注者，接下来在对follower表和attention表分别基于userId进行水平拆分。

kaoyanyun.user的blog表：

![福昕截屏20230129150435191](.\images\福昕截屏20230129150435191.PNG)

贴子表。

kaoyanyun.user的file表：

![福昕截屏20230129150507230](.\images\福昕截屏20230129150507230.PNG)

这个表用于存放文件的元信息。

kaoyanyun.user的comment表：

![福昕截屏20230129150609990](.\images\福昕截屏20230129150609990.PNG)

目前还没有考虑评论底下还可以继续评论的问题。type字段是为了增加项目的可拓展性，形成一个业务中台。

kaoyanyun.user的user表：

![福昕截屏20230129150746551](.\images\福昕截屏20230129150746551.PNG)

密码存放的是盐值加密后的密码。

# 缓存体系

- 用户缓存，保存用户的元信息
- 关注列表缓存，里面存放关注列表的ID
- 模糊计数缓存，里面存放的是模糊计数的数量
- 粉丝列表缓存，里面存放粉丝列表的ID
- 点赞列表缓存，里面存放点赞列表的ID
- 精准计数缓存，里面存放的是精准计数的数量
- Feed缓存，里面存放的是用户Feed流的文章ID
- 文章内容缓存，里面记录文章的内容
- 是否存在缓存，里面记录是否点赞过，关注过
- 签到缓存
- GEO缓存
- hyperloglog缓存
- 用户登录信息缓存
- 文件相关信息缓存，用于实现秒传等功能
- 文件下载进度缓存

# 项目核心亮点介绍

## 登录系统设计

使用JWT+Redis进行设计。

使用JWT的核心原因：

- 用cookie在分布式的应用上会限制负载均衡的能力，如果想要起作用需要每一个服务都复制一份session，浪费过多存储空间。
- 会有被CSRF攻击的风险。
- 方便获取用户的相关信息，例如用户ID

JWT其实已经可以了为什么还要用Redis？

- JWT虽然是无状态的，但是一旦被泄漏，就无能为力了，安全性仍然存在一定问题。
- 服务端把权限控制在手是更好的一种选择，既然要控制在手，而且又要让所有的服务都可以看到，那么Redis是一个很好的选择。
- 服务端把权限控制在手可以控制用户的行为，用户用同一种设备登录时可以踢除等等。

### 风险控制

#### 多设备登录校验

- 第一次登录的token：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzZXJuYW1lIjoidG9tIiwiZXhwIjoxNjU3OTgzNzAxLCJpc3MiOiJnaW4tSU0iLCJuYmYiOjE2NTc4OTczMDF9.ipiIDgAdTwrv8EX45y0UD6wy0fOOdzhIDysyB8kJais
- 第二次登录的token：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzZXJuYW1lIjoidG9tIiwiZXhwIjoxNjU3OTgzODAxLCJpc3MiOiJnaW4tSU0iLCJuYmYiOjE2NTc4OTc0MDF9.3ZDrBr0FaFKpcicJpNkvEVCd8UdEQp079mg4fr2jBcc

通过测试可以发现，两个token都是有效的。只有到了指定的日期后，token才会失效。这显然是不合理的。所以我在这里的处理是：**引入Redis**。

具体逻辑：使用redis中的hash结构。`key`为user的唯一标识uid；`filed`为该user的User-Agent，表示是哪一个设备（同一个设备只能有1个token有效）；`value`存储该user的唯一有效token。

结构如下：

![登录测试1](.\images\登录测试1.PNG)

- 在`middleware/jwt.go`中增加判定逻辑。通过`uid`和`User-Agent`（从解析token中包含的相关信息claims中获取uid），查出redis中的token。判定携带的token是否和redis中的token一样。如果不一样说明是旧的token，直接`c.Abort()`然后`return`。

> 这只是我自己的一个想法，如果以后发现更好的解决方案，会继续更新的。

更新后：

- 第一次token：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzZXJuYW1lIjoidG9tIiwiZXhwIjoxNjU4MDU0NzUxLCJpc3MiOiJnaW4tSU0iLCJuYmYiOjE2NTc5NjgzNTF9.-FvhHHpJokeigiSJOUkTWaQ4ytsYDZcxaTklPLzJGR4
- 第二次token：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzZXJuYW1lIjoidG9tIiwiZXhwIjoxNjU4MDU0NzgwLCJpc3MiOiJnaW4tSU0iLCJuYmYiOjE2NTc5NjgzODB9.uBkmCpbTfEbr3fBiMQ26XrxOQc-hl6H5jvS_3BfW-2o

可以发现使用第一次token去请求会403：

![登录测试2](.\images\登录测试2.PNG)

据说微信就是这样做的（跟群友讨论的）：

- 这不就是提掉线吗？
- 登录后 将以前此用户的token删除掉即可
- 如果想多设备登录 就加入设备就可以 当前token和用户id，设备绑定
- 微信就是这样做的

跟大家讨论，感觉基本都是基于redis缓存token的，踢掉用户也是这么干的。

不过感觉这样就跟jwt的无状态背道而驰了，回到了session。如果以后有更优雅更好的方式，会再记录的。

#### 验证码接口防刷

验证码我们需要防刷的，后台也要做！！

有的人可能会问，为什么后台也要做呢？前端这边做一下限制不就可以了吗。我发送一次验证码值用户在一定的时间内就不可以继续发送了，这样不就实现了防刷了吗？但是实际上并不是这样的。因为前端虽然可以防住小白用户，但是有些很懂的人，可以抓包或者发送验证码的请求，然后用JMeter爆刷！！这样前端的限制就没有用了。

但是这个时候存在两个问题：

- 我要的前台和后台都要去防止对面去恶意刷验证码。因为我们的请求的时候，发送验证码的请求已经被暴露了。对方懂行的人可以去用JMeter刷验证码的
- 虽然我的前端限制验证码是过60秒之后才能继续刷的，但是我们刷新一次页面之后就会发现，这个限制没有了，因为我们没有后台记录所剩下的时间，我只需要在游览器上面刷新游览器的缓存，那么保存在游览器里面的时间就被刷新了，等于说这个东西是无意义的

我们接下来要解决这两个问题。

首先验证码可以存到redis里面，并且设置过期时间。

**解决页面刷新之后验证码就可以马上发送的情况：**

#### 密码多次错误惩戒

代码和上面的差不多，不做过多解释。

## 抢红包系统是怎么设计的？

![抢红包业务分析.excalidraw](.\images\抢红包业务分析.excalidraw.png)

可以明显的看到打开了红包不一定可以抢到。这样做的好处是：

- 符合现实生活逻辑，有仪式感
- 防止误领，发现不对劲可以马上退出
- 流程拆的长一些，平摊高并发下的压力

![红包.excalidraw](.\images\红包.excalidraw.png)

预拆包：我在发红包的时候，就已经把所有的东西都计算好了放在redis里面了。

实时计算：在抢红包的时候进行现场计算。

我选择的时把红包提前拆好，虽然提前拆好需要占用Redis的一些存储，但是这样可以让抢红包的速度到达最快。

### 发红包

发红包的时候需要一些算法，这里采用的是**二倍均值法**

二倍均值法可以在正态分布的前提下做到抢红包相对的公平。计算方法如下：

每次抢到的金额 = 随机区间(0, M/N*2)。

![image-20230102140655356](.\images\image-20230102140655356.png)

发红包的时候需要指定红包的金额和个数。然后通过二倍均值法计算后放入Redis。这里我们采用的数据结构是List类型，因为List抢走一个红包和红包库存呢减一是一个原子性操作。而且List类型的数据结构就决定了获取库存的时间复杂度是O(1)，是可行的。到时候按照顺序弹出即可。

![image-20221231160638576](.\images\image-20221231160638576.png)

### 拆红包

用户点击红包的时候查看库存，如果库存为0的话直接显示红包已经抢完了，并且把详细信息显示。如果还有库存的话，就可以进入抢红包的环节了。

### 抢红包

用户点击抢红包，此时再次判断库存是否充足，如果充足的话判断用户是否抢过红包，如果两个条件都满足的话，就抢红包成功，从Redis里面弹出一个红包给用户，更新对应的Redis(红包总金额要扣除了)，然后使用消息队列异步的调用服务，将用户的金额进行相应的加上。

![image-20221231161312676](.\images\image-20221231161312676.png)

与此同时，需要记录哪些用户已经抢过红包了。

![image-20221231161635993](.\images\image-20221231161635993.png)

### 退款(目前还没有实现)

使用延迟队列，24小时后进行消息的消费，如果发现数据库里面还有相应的缓存，那么就需要把剩下的金额退还给用户。具体过程是查看Redis的某个key是否还存在。如果还存在的话那么就把List里面的红包全部弹出，然后把金额归还给用户。然后把红包key删除掉。

### 超卖问题

因为查看**是否还有库存**和**扣减库(抢红包)存**两个操作不是原子性的，因此我们可以使用Redis+Lua实现原子性操作。

Lua实现抢红包的流程：

- 查询用户是否抢过红包
- 查询是否还有红包
- 有的话就扣减红包Redis弹出一个
- 没有的话就返回

使用 EVAL 命令每次请求都需要传输 Lua 脚本 ，若 Lua 脚本过长，不仅会消耗网络带宽，而且也会对 Redis 的性能造成一定的影响。

思路是先将 Lua 脚本先缓存起来 , 返回给客户端 Lua 脚本的 sha1 摘要。 客户端存储脚本的 sha1 摘要 ，每次请求执行 EVALSHA 命令即可。

![福昕截屏20230129144224941](.\images\福昕截屏20230129144224941.PNG)

## 大文件

### 上传

- 分片
- 断点
- 秒传（判断文件哈希值）

前端不断的发送请求，如果用户暂停上传的话，那么就是前端停止发送请求就可以了。我分片了，而且记录了分片的相关信息，所以实现了断点功能。

前端把文件进行分片，然后用**多次**请求发给后端，请求的内容里面包含了文件的很多相关信息。

我们约定常量key，也就是redis里面的key：

```go
/**
 * 常量表
 */
const (
	// FileMd5Key 保存文件所在的路径 eg:FILE_MD5:468s4df6s4a
	FileMd5Key = "FILE_MD5:"
	// FileUploadStatus 保存上传文件的状态
	FileUploadStatus = "FILE_UPLOAD_STATUS"
)
```

约定枚举类：

```go
/**
 * 1 开头为判断文件在系统的状态
 */
const (
	// IsHave 文件以及存在了
	IsHave = 160
	// NoHave 该文件没有上传过
	NoHave = 161
	// IngHave 该文件上传了一部分
	IngHave = 162
)
```

FILE_UPLOAD_STATUS里面存放的值是false或者true，如果是false的话就说明文件上传没有完成，如果是true的话就说明文件的上传已经完成了，我们的系统里面保存有这个文件。

判断文件的上传状态以及获取文件还有哪些分片还有被完成。

这里是基于bitmap去做的，如果位是1，那么就代表这个分片以及上传完成了，如果位是0表示这个分片还没有上传完成。

但是我们也可以创建一个.conf文件，把相关信息使用.conf文件记录下来也可以。当比特位的值是Byte.MAX_VALUE的时候就代表着这个位置的分片是已经上传完成了。

### 下载

实际上需要客户端和服务器：我的用户发送一个请求，说要下载，然后请求是发送给客户端的，客户端先去问服务器支不支持断点下载，如果支持的话就不断的发送请求给服务器实现断点下载。因此我的项目需要实现的有客户但，以及服务器。

![福昕截屏20230129144451926](.\images\福昕截屏20230129144451926.PNG)

一个比较常见的场景,就是断点续传/下载,在网络情况不好的时候,可以在断开连接以后,仅继续获取部分内容. 例如在网上下载软件,已经下载了 95% 了,此时网络断了,如果不支持范围请求,那就只有被迫重头开始下载.但是如果有范围请求的加持,就只需要下载最后 5% 的资源,避免重新下载.

另一个场景就是多线程下载,对大型文件,开启多个线程, 每个线程下载其中的某一段,最后下载完成之后, 在本地拼接成一个完整的文件,可以更有效的利用资源.

![福昕截屏20230129144525582](.\images\福昕截屏20230129144525582.PNG)

#### Range & Content-Range

HTTP1.1 协议（RFC2616）开始支持获取文件的部分内容,这为并行下载以及断点续传提供了技术支持. 它通过在 Header 里两个参数实现的,客户端发请求时对应的是 Range ,服务器端响应时对应的是 Content-Range.

```
$ curl --location --head 'https://download.jetbrains.com/go/goland-2020.2.2.exe'
date: Sat, 15 Aug 2020 02:44:09 GMT
content-type: text/html
content-length: 138
location: https://download-cf.jetbrains.com/go/goland-2020.2.2.exe
server: nginx
strict-transport-security: max-age=31536000; includeSubdomains;
x-frame-options: DENY
x-content-type-options: nosniff
x-xss-protection: 1; mode=block;
x-geocountry: United States
x-geocode: US

HTTP/1.1 200 OK
Content-Type: binary/octet-stream
Content-Length: 338589968
Connection: keep-alive
x-amz-replication-status: COMPLETED
Last-Modified: Wed, 12 Aug 2020 13:01:03 GMT
x-amz-version-id: p7a4LsL6K1MJ7UioW7HIz_..LaZptIUP
Accept-Ranges: bytes
Server: AmazonS3
Date: Fri, 14 Aug 2020 21:27:08 GMT
ETag: "1312fd0956b8cd529df1100d5e01837f-41"
X-Cache: Hit from cloudfront
Via: 1.1 8de6b68254cf659df39a819631940126.cloudfront.net (CloudFront)
X-Amz-Cf-Pop: PHX50-C1
X-Amz-Cf-Id: LF_ZIrTnDKrYwXHxaOrWQbbaL58uW9Y5n993ewQpMZih0zmYi9JdIQ==
Age: 19023
```

#### Range

The Range 是一个请求首部,告知服务器返回文件的哪一部分. 在一个 Range 首部中,可以一次性请求多个部分,服务器会以 multipart 文件的形式将其返回. 如果服务器返回的是范围响应,需要使用 206 Partial Content 状态码. 假如所请求的范围不合法,那么服务器会返回 416 Range Not Satisfiable 状态码,表示客户端错误. 服务器允许忽略 Range 首部,从而返回整个文件,状态码用 200 .`Range:(unit=first byte pos)-[last byte pos]`

Range 头部的格式有以下几种情况：

```
Range: <unit>=<range-start>-
Range: <unit>=<range-start>-<range-end>
Range: <unit>=<range-start>-<range-end>, <range-start>-<range-end>
Range: <unit>=<range-start>-<range-end>, <range-start>-<range-end>, <range-start>-<range-end>
```

#### Content-Range

假如在响应中存在 Accept-Ranges 首部（并且它的值不为 “none”）,那么表示该服务器支持范围请求(支持断点续传). 例如,您可以使用 cURL 发送一个 `HEAD` 请求来进行检测.`curl -I http://i.imgur.com/z4d4kWk.jpg`

```
HTTP/1.1 200 OK
...
Accept-Ranges: bytes
Content-Length: 146515
```

在上面的响应中, `Accept-Ranges: bytes` 表示界定范围的单位是 bytes . 这里 `Content-Length` 也是有效信息,因为它提供了要检索的图片的完整大小.

如果站点未发送 Accept-Ranges 首部,那么它们有可能不支持范围请求.一些站点会明确将其值设置为 “none”,以此来表明不支持.在这种情况下,某些应用的下载管理器会将暂停按钮禁用.

```
Run go run main.go
2020/08/15 02:15:31 开始[9]下载from:376446150 to:418273495
2020/08/15 02:15:31 开始[0]下载from:0 to:41827349
2020/08/15 02:15:31 开始[1]下载from:41827350 to:83654699
2020/08/15 02:15:31 开始[5]下载from:209136750 to:250964099
2020/08/15 02:15:31 开始[6]下载from:250964100 to:292791449
2020/08/15 02:15:31 开始[7]下载from:292791450 to:334618799
2020/08/15 02:15:31 开始[2]下载from:83654700 to:125482049
2020/08/15 02:15:31 开始[8]下载from:334618800 to:376446149
2020/08/15 02:15:31 开始[4]下载from:167309400 to:209136749
2020/08/15 02:15:31 开始[3]下载from:125482050 to:167309399
2020/08/15 02:15:36 开始合并文件
2020/08/15 02:15:38 文件SHA-256校验成功

 文件下载完成耗时: 7.169149 second
```

## 敏感词系统是怎么设计的？

AC自动机即可。

### AC自动机

#### AC自动机是干嘛的？

我有一个敏感词数组，里面装的是所有的敏感词，还有一篇大文章，我要求出大文章里面所有的敏感词。

敏感词数组本身的组织是一颗前缀树。

AC自动机就是在前缀树的基础上做升级。

#### 流程

- 我们在前缀树的基础上给每一个节点加上`fail`指针并且做出规定：头节点的`fail`指针一定指向`null`

- 头节点的下级直接节点的`fail`指针一律指向头部

- 在**整颗**前缀树全部建立完毕之后，再去建立`fail`指针

- 以下是其他节点的定义规则，看图说话。

  ![AC自动机.excalidraw](.\images\AC自动机.excalidraw.png)

  - 假设有一个节点X，X的父亲节点P，P的`fail`指针指的是谁，看图可以知道指向的是头节点。

  - X的父亲节点到X的路径存放的是`b`，因此我询问X的父亲节点的`fail`指针指向的节点，也就是图中的头节点，你的路径中有没有`b`，可以看到头节点有的路径是`a`, `c`, `e`，没有`b`。

  - 于是继续往跳，头部节点的`fail`指针指向的节点是`null`,`null`当然不会有`b`路径，因此X的`fail`指针直接指向头部

  - 再看S节点。

    ![AC自动机1.excalidraw](.\images\AC自动机1.excalidraw.png)

  - S的父亲节点X的`fail`指针指向的节点是头节点，它们之间的路径是c，而头节点有c这个路径，所以X的`fail`指针指向的是头节点的以c为路径的孩子节点，如图。

    ![AC自动机2.excalidraw](.\images\AC自动机2.excalidraw.png)

#### AC自动机的fail指针的作用

我们再来画一个图：

![AC自动机3.excalidraw](.\images\AC自动机3.excalidraw.png)

`fail`指针的含义比较抽象，但是我们还是尝试去概括一下：

> 当字符串无法匹配时，我们有最后一个字符，我们命名为`last`，当必须以`last`结尾时，与字符串拥有同一后缀的最长的字符串，`fail`指针的作用就是方便的找到这样一个字符串。

我们看到上图：

有节点X，假设字符串`abc`就是我们无法匹配成功的字符串。`fail`指针指向的节点和头节点连接而成的路径是c，那么这个字符串`c`实际上就是与`abc`拥有同一后缀并且最长的字符串。

![image-20221225230728118](.\images\image-20221225230728118.png)

有节点Y，假设字符串`abcd`就是我们无法匹配成功的字符串。Y的`fail`指针指向的节点与头节点连成的字符串是`cd`，那么`cd`就是与`abcd`拥有相同最长后缀的字符串，与`abcd`拥有相同后缀的字符串还有`c`，但是`c`没有`cd`长，所以`fail`指针没有指向另一头的节点。

#### 大文章敏感词匹配

![AC自动机4.excalidraw](.\images\AC自动机4.excalidraw.png)

- 我们有大文章`abcdex`，我们对着这个AC自动机从0位置开始进行匹配，发现只能匹配到字符串`abcde`，因此得出结论，从0位置开始匹配，是无法匹配出敏感词的。
- 我们匹配失败了，只能匹配到字符串`abcde`，此时的节点是X，这时候，我们就跳往X的`fail`指针指向的位置的节点Y。
- 然后我们从头节点到Y的字符串是`cde`，因此，我们得到了最长的前缀保留，看，是不是跟`KMP`算法非常类似，AC自动机不过是在前缀树上的`KMP`算法。

## 计数服务是怎么设计的？

我一开始没有准备这个技术服务，但是后来又加上了，为什么呢？因为我发现有太多计数了，点赞数，评论数，转发数，收藏数，浏览数，而这些计数的代码逻辑都是相似的，因此我想设计一个计数服务把这些业务逻辑进行解耦。

而计数分为两种，第一种是**模糊计数**，第二种是**精准计数**。模糊计数的话我们可以读写都是Redis，然后启动一个后台线程去异步的对齐更新。而精准计数为了精准，我们采用写入的时候直接删除缓存的策略。

其次就计数服务的定时器应该定时睡眠，否则对CPU的占用太严重了

普通用户的操作和大V用户不应该一样，部分大V用户不应该设置缓存过期时间，或者是不采用插入就删除缓存的策略。

## Redis数据结构的应用

引用位图，hyper log log，GEO实现了一些功能。

## 解决缓存的常见问题

### 缓存击穿

放了空值解决。

### 缓存穿透，雪崩

互斥锁+随机缓存过期时间。

### 一致性问题

先删数据库再删缓存，引入消息队列，删除失败就持续的进行删除，因为没有搭建集群，所以暂无延迟双删。

## 部署

正在加工中。。。。




















