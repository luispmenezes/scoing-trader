package trader

import (
	"github.com/shopspring/decimal"
	"log"
	"math/rand"
	"scoing-trader/trader/model/predictor"
	"scoing-trader/trader/model/trader"
	"scoing-trader/trader/model/trader/strategies"
	"sort"
	"sync"
	"time"
)

type Evolution struct {
	Predictions    []predictor.Prediction
	InitialBalance decimal.Decimal
	Fee            decimal.Decimal
	Uncertainty    float64
	GenerationSize int
	NumGenerations int
	MutationRate   float64
	StartingPoint  []float64
	StrategyName   string
}

type Specimen struct {
	Fitness decimal.Decimal
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
				Fitness: decimal.Zero,
				Config:  &config,
			})
	}

	for i := 0; i < evo.NumGenerations; i++ {
		testedSpecimens := evo.simulateGeneration(specimenPool)
		newCandidates := evo.selectCandidates(testedSpecimens, 2)

		log.Printf("Generation %d Fitness: %s", i, newCandidates[0].Fitness)
		log.Println(newCandidates[0].Config.ToSlice())

		if len(candidates) == 0 || newCandidates[0].Fitness.GreaterThan(candidates[0].Fitness) {
			candidates = newCandidates
		} else {
			log.Printf("New Generation worse than overall best (%s)", candidates[0].Fitness)
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
				Fitness: decimal.Zero,
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
	sim := NewSimulation(predictions, strategy, specimen.Config, evo.InitialBalance, evo.Fee, evo.Uncertainty, false, false)
	sim.Run()
	specimen.Fitness = sim.Trader.Accountant.NetWorth()
	out <- specimen
}

func (evo *Evolution) selectCandidates(specimens []Specimen, numCandidates int) []Specimen {
	sort.Slice(specimens, func(i, j int) bool {
		return specimens[i].Fitness.GreaterThan(specimens[j].Fitness)
	})
	return specimens[0:numCandidates]
}
