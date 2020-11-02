
### GPM模型

![avatar](http://www.topgoer.com/static/7.1/gmp/12.jpg)

G：go-routine  
P：processor 本地队列 全局队列 
M：thread 

####调度流程

![avatar](http://www.topgoer.com/static/7.1/gmp/13.jpg)

- 创建G
- 入本地队列 | 入全局队列 
- M从P获取G | 从全局P获取 | 从其他P偷取G 
- 调度
- 执行  | 阻塞 | 重新获取M | 接管P 
- 销毁G
- 返回