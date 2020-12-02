package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// добавляет перед каждым этапом конвеера обработку канала прерывания.
func addInterruptor(in In, done In, stageToRun Stage) Out {
	stageIn := make(Bi)

	go func() {
		defer close(stageIn)
		for {
			select {
			case <-done:
				return
			case n, good := <-in:
				if !good {
					return
				}
				stageIn <- n
			}
		}
	}()

	return stageToRun(stageIn)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	nextStageInput := in

	for _, stage := range stages {
		if stage == nil {
			continue
		}

		nextStageInput = addInterruptor(nextStageInput, done, stage)
	}

	return nextStageInput
}
