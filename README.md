# Nil
mc-nil后端代码

名词解释：mc即Minecraft，nil即为空（站名没有想好）

项目目标：搭建一个开放式的MC交流论坛

后端api文档：https://apifox.com/apidoc/shared-2b10cb9d-7ce5-41ff-8f03-267d95d76079



### 实现功能：

###### 	1.关于用户

​	1.1基于邮件验证码激活的用户注册

​	1.2用户密码采用MD5单向加密存入

​	1.3记录用户的关注用户列表，发表帖子列表，评论列表，浏览记录(基于LRU缓存算法)列表

​	1.4用户点赞



###### 	2.关于板块

​	2.1板块的创建以及信息获取(详(单个)/略(全部)信息拆分，各一个接口)



###### 	3.关于帖子

​	3.1帖子发布

​	3.2根据时间或者得分（点赞数）获取帖子列表

​	3.3基于板块以及时间/分数获取帖子列表



###### 	4.关于评论

​	4.1帖子评论的实现

​	4.2二级评论的实现



###### 	5.其他

​	5.1部分接口采用JWT的Bearer用户令牌进行身份校验。其中包含uid，username两个信息

​	5.2基于立牌算法进行限流处理

​	5.3基于雪花算法进行uid，chunk_id，post_id,comment_id的生成

​	5.4基于jordan-wright/email以及smtp技术使用网易邮箱进行激活码邮件的发送

​	5.5基于阿里云oss云存储技术实现图片，视频等文件上传

​	5.6基于docker容器化进行服务器部署服务



**实现的接口功能目录如下：**

![image-20240427181558429](https://rosyrain.oss-cn-hangzhou.aliyuncs.com/img2/202404271815566.png)



除上传文件的接口均已部署到 74.48.160.188:5000，可自行使用

