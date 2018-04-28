# nsq传输设置

* topic：设置为"deploy"
* 传输数据：

	```
	{
		//命令类型，（构建：0， 部署：1，回滚：2）
		OrderType 	int 	`json:"order_type"`
		//对应操作的数据库ID（构建：构建记录表ID，部署：部署记录表ID，回滚：部署记录表ID）
		DataId 	  	int64 	`json:"data_id"`
	}
	```


