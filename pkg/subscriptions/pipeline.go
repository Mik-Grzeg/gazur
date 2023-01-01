package subscriptions

type Executor func(interface{}) interface{}

type TaskSupervisor struct {
	errC  chan error
	tasks []task
}

type TaskManager interface {
	AddSource(executor Executor, inC <-chan interface{}) TaskManager
	AddTask(executor Executor) TaskManager
	Run() chan interface{}
	//Merge() <-chan interface{}
}

type task struct {
	inC      <-chan interface{}
	outC     chan interface{}
	errC     chan error
	executor Executor
}

func (t *task) run() {
	for {
		select {
		case elem := <-t.inC:
			t.outC <- t.executor(elem)
		case <-t.errC:
			return
		}
	}
}

func (s *TaskSupervisor) Run() chan interface{} {
	for _, task := range s.tasks {
		go task.run()
	}
	return s.tasks[len(s.tasks)-1].outC
}

func (s *TaskSupervisor) AddSource(executor Executor, inC <-chan interface{}) TaskManager {
	outC := make(chan interface{})

	task := task{
		inC,
		outC,
		s.errC,
		executor,
	}

	s.tasks = append(s.tasks, task)
	return s
}

func (s *TaskSupervisor) AddTask(executor Executor) TaskManager {
	inC := s.tasks[len(s.tasks)-1].outC
	outC := make(chan interface{})

	task := task{
		inC,
		outC,
		s.errC,
		executor,
	}

	s.tasks = append(s.tasks, task)
	return s
}

func NewTaskManager() TaskManager {
	return &TaskSupervisor{
		errC:  make(chan error),
		tasks: []task{},
	}
}

//func run(
//	inC <-chan interface{},
//	f Executor,
//) (chan interface{}, chan error) {
//	outC := make(chan interface{})
//	errC := make(chan error)
//
//	go func() {
//		defer close(outC)
//		for v := range inC {
//			res, err := f(v)
//			if err != nil {
//				errC <- err
//				continue
//			}
//			outC <- res
//		}
//	}()
//	return outC, errC
//}
