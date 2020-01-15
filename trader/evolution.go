package trader

import (
	"log"
	"math/rand"
	"sort"
	"super-trader/trader/model/predictor"
	"super-trader/trader/model/trader"
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
}

type Specimen struct {
	Fitness float64
	Config  trader.TraderConfig
}

func (evo *Evolution) Run() Specimen {
	rand.Seed(time.Now().UnixNano())

	var specimenPool []Specimen
	var candidates []Specimen

	for i := 0; i < evo.GenerationSize; i++ {
		specimenPool = append(specimenPool,
			Specimen{
				Fitness: 0,
				Config:  *trader.RandomConfig(-1, 1),
			})
	}

	for i := 0; i < evo.NumGenerations; i++ {
		testedSpecimens := evo.simulateGeneration(specimenPool)
		newCandidates := evo.selectCandidates(testedSpecimens, 2)

		log.Printf("Generation %d Fitness: %f", i, newCandidates[0].Fitness)

		if len(candidates) == 0 || newCandidates[0].Fitness > candidates[0].Fitness {
			candidates = newCandidates
		} else {
			log.Printf("New Generation worse than overall best (%f)", candidates[0].Fitness)
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
		child := *trader.RandomBetweenTwo(candidates[0].Config, candidates[1].Config)

		if rand.Float64() <= evo.MutationRate {
			for i := 0; i < (candidates[0].Config.NumParams() / 3); i++ {
				child.RandomizeParam(-1, 1)
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
	resultChan := make(chan Specimen, evo.GenerationSize)
	var testedSpecimens []Specimen
	var wg sync.WaitGroup

	for _, specimen := range untestedSpecimens {
		wg.Add(1)
		go evo.runSingleSimulation(specimen, evo.Predictions, resultChan, &wg)
	}

	for i := 0; i < evo.GenerationSize; i++ {
		testedSpecimens = append(testedSpecimens, <-resultChan)
	}

	return testedSpecimens
}

func (evo *Evolution) runSingleSimulation(specimen Specimen, predictions []predictor.Prediction, out chan<- Specimen, wg *sync.WaitGroup) {
	defer wg.Done()
	sim := NewSimulation(predictions, specimen.Config, evo.InitialBalance, evo.Fee, evo.Uncertainty, false)
	sim.Run()
	specimen.Fitness = sim.Trader.Wallet.NetWorth()
	out <- specimen
}

func (evo *Evolution) selectCandidates(specimens []Specimen, numCandidates int) []Specimen {
	sort.Slice(specimens, func(i, j int) bool {
		return specimens[i].Fitness > specimens[j].Fitness
	})
	return specimens[0:numCandidates]
}
