package snowflake

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

var node *sf.Node

//初始化方法：初始化一个全局node节点
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	//指定时间因子,这个时间就从startTime开始往后能用69年,sonyflake可以用174年
	st, err = time.Parse("", startTime)
	if err != nil {
		return
	}
	//初始化开始时间，毫秒值
	sf.Epoch = st.UnixNano() / 1000000
	//指定机器ID，machineID是机器的唯一标识
	node, err = sf.NewNode(machineID)
	return
}

func GenID() int64 {
	//拿到node节点后就可以生成ID值了
	return node.Generate().Int64()
}

//仅用于测试执行，正常代码中用不到
//func main() {
//	//传入startTime的时间和machineID
//	if err := Init("2021-07-05", 1);err != nil{
//		fmt.Printf("init failed, err:%v", err)
//		return
//	}
//	id := GenID()
//	fmt.Println(id)
//}
