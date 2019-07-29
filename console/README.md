## goCron: A Golang Job Scheduling Package.

goCron是一个Golang作业调度包，它允许您使用简单，人性化的语法以预定间隔定期运行Go函数。

回到这个包，您可以使用下面这个简单的API来运行cron调度程序。
``` go
package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
)

func task() {
	fmt.Println("I am runnning task.")
}

func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}

func main() {
	// Do jobs with params
	gocron.Every(1).Second().Do(taskWithParams, 1, "hello")

	// Do jobs without params
	gocron.Every(1).Second().Do(task)
	gocron.Every(2).Seconds().Do(task)
	gocron.Every(1).Minute().Do(task)
	gocron.Every(2).Minutes().Do(task)
	gocron.Every(1).Hour().Do(task)
	gocron.Every(2).Hours().Do(task)
	gocron.Every(1).Day().Do(task)
	gocron.Every(2).Days().Do(task)

	// Do jobs on specific weekday
	gocron.Every(1).Monday().Do(task)
	gocron.Every(1).Thursday().Do(task)

	// function At() take a string like 'hour:min'
	gocron.Every(1).Day().At("10:30").Do(task)
	gocron.Every(1).Monday().At("18:30").Do(task)

	// remove, clear and next_run
	_, time := gocron.NextRun()
	fmt.Println(time)

	gocron.Remove(task)
	gocron.Clear()

	// function Start start all the pending jobs
	<- gocron.Start()

	// also , you can create a your new scheduler,
	// to run two scheduler concurrently
	s := gocron.NewScheduler()
	s.Every(3).Seconds().Do(task)
	<- s.Start()

}
```

完整的测试用例和[文档](http://godoc.org/github.com/jasonlvhit/gocron)即将推出。

再一次，感谢Ruby clockwork和Python schedule package的伟大作品。使用BSD许可证，请参阅文件许可证以获取详细信息。

玩得开心！