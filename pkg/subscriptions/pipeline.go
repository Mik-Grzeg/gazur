package subscriptions

type Executor[In any, Out any] func(chan In) (chan Out)

type TaskSupervisor[In any, Out any] struct {
	errC chan error
	tasks []task[In, Out]
}

type TaskManager[In any, Out any] interface {
	AddTask(executor Executor[In, Out]) task[In, Out]
	Merge() <-chan interface{}
}

type task[In any, Out any] struct {
	inC chan In
	outC chan Out
	executor Executor[In, Out]
}

func (s *TaskSupervisor[In, Out]) AddTask(executor Executor[In, Out]) task[In, Out] {
	inC := make(chan In)
}

type pipeline[In any, Out any] struct {
	dataC     chan In
	errC      chan error
	executors []Executor
}

func (p *pipeline[In, Out]) Pipe(executor Executor) Pipeline {
	p.executors = append(p.executors, executor)
	return p
}

func (p *pipeline[In, Out]) Merge() <-chan interface{} {
	for i := 0; i < len(p.executors); i++ {
		p.dataC, p.errC = run(p.dataC, p.executors[i])
	}
	return p.dataC
}

func New[In any, Out any](f func(chan In)) Pipeline {
	inC := make(chan In)
	go f(inC)
	return &pipeline[In, Out]{
		dataC:     inC,
		errC:      make(chan error),
		executors: []Executor{},
	}
}

func run(
	inC <-chan interface{},
	f Executor,
) (chan interface{}, chan error) {
	outC := make(chan interface{})
	errC := make(chan error)

	go func() {
		defer close(outC)
		for v := range inC {
			res, err := f(v)
			if err != nil {
				errC <- err
				continue
			}
			outC <- res
		}
	}()
	return outC, errC
}
