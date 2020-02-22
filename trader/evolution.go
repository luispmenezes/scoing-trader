package trader

import (
	"log"
	"math/rand"
	"scoing-trader/trader/model/market/model"
	"scoing-trader/trader/model/predictor"
	"scoing-trader/trader/model/trader"
	"scoing-trader/trader/model/trader/strategies"
	"sort"
	"sync"
	"time"
)

type Evolution struct {
	Predictions    []predictor.Prediction
	InitialBalance float64
	Fee            float64
	Uncertainty    float64
	GenerationSize int
	NumGenerations int
	MutationRate   float64
	StartingPoint  []float64
	StrategyName   string
}

type Specimen struct {
	Fitness int64
	Config  trader.StrategyConfig
}

func (evo *Evolution) Run() Specimen {
	rand.Seed(time.Now().UnixNano())

	var specimenPool []Specimen
	var candidates []Specimen

	for i := 0; i < evo.GenerationSize; i++ {
		config := strategies.BasicWithMemoryConfig{}
		if evo.StartingPoint != nil {
			config.FromSlice(evo.StartingPoint)
			for j := 0; j < i; j++ {
				config.RandomizeParam()
			}
		} else {
			config.RandomFromSlices(config.ParamRanges())
		}

		specimenPool = append(specimenPool,
			Specimen{
				Fitness: 0,
				Config:  &config,
			})
	}

	for i := 0; i < evo.NumGenerations; i++ {
		testedSpecimens := evo.simulateGeneration(specimenPool)
		newCandidates := evo.selectCandidates(testedSpecimens, 2)

		log.Printf("Generation %d Fitness: %s", i, model.IntToString(newCandidates[0].Fitness))
		log.Println(newCandidates[0].Config.ToSlice())

		if len(candidates) == 0 || newCandidates[0].Fitness > candidates[0].Fitness {
			candidates = newCandidates
		} else {
			log.Printf("New Generation worse than overall best (%s)", model.IntToString(candidates[0].Fitness))
		}

		specimenPool = evo.breed(candidates)
	}

	return candidates[0]
}

func (evo *Evolution) breed(candidates []Specimen) []Specimen {
	var newGeneration []Specimen

	//rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	//TODO: Permitir mais de dois candidatos
	for i := 0; i < evo.GenerationSize; i++ {
		child := &strategies.BasicWithMemoryConfig{}
		child.RandomFromSlices(candidates[0].Config.ToSlice(), candidates[1].Config.ToSlice())

		if rand.Float64() <= evo.MutationRate {
			for i := 0; i < (child.NumParams() / 3); i++ {
				child.RandomizeParam()
			}
		}

		newGeneration = append(newGeneration,
			Specimen{
				Fitness: 0,
				Config:  child,
			})
	}

	return newGeneration
}

func (evo *Evolution) simulateGeneration(untestedSpecimens []Specimen) []Specimen {
	var testedSpecimens []Specimen

	resultChan := make(chan Specimen, evo.GenerationSize)
	var wg sync.WaitGroup

	for _, specimen := range untestedSpecimens {
		wg.Add(1)
		go evo.runSingleSimulation(specimen, &evo.Predictions, resultChan, &wg)
	}

	for i := 0; i < evo.GenerationSize; i++ {
		testedSpecimens = append(testedSpecimens, <-resultChan)
	}

	return testedSpecimens
}

func (evo *Evolution) runSingleSimulation(specimen Specimen, predictions *[]predictor.Prediction, out chan<- Specimen, wg *sync.WaitGroup) {
	defer wg.Done()
	strategy := strategies.NewBasicWithMemoryStrategy(specimen.Config.ToSlice(), 10)
	sim := NewSimulation(predictions, strategy, specimen.Config, evo.InitialBalance, evo.Fee, evo.Uncertainty, false)
	sim.Run()
	specimen.Fitness = sim.Trader.Accountant.NetWorth()
	out <- specimen
}

func (evo *Evolution) selectCandidates(specimens []Specimen, numCandidates int) []Specimen {
	sort.Slice(specimens, func(i, j int) bool {
		return specimens[i].Fitness > specimens[j].Fitness
	})
	return specimens[0:numCandidates]
}
