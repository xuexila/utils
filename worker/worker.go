package worker

type Job struct {
	// 工作函数
	Func func(interface{})
	// 工作参数
	Params interface{}
}

type worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
}

func (w worker) start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel
			job := <-w.JobChannel
			// we have received a work request.
			job.Func(job.Params)
		}
	}()
}

// 定义开始工作结构体
type StartWorker struct {
	// 最大运行数
	MaxSize    int `ini:"max_size"`
	WorkerPool chan chan Job
}

// 初始化工作池
func (s StartWorker) Init() {
	for i := 0; i < s.MaxSize; i++ {
		worker := worker{
			WorkerPool: s.WorkerPool,
			JobChannel: make(chan Job),
		}
		worker.start()
	}
}

// 运行
func (s StartWorker) Run(j *Job) {
	jobChannel := <-s.WorkerPool
	jobChannel <- *j
}
