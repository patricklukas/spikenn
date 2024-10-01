package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Flatten a 2D coordinate to a 1D index
func flattenIdx(row, col, numCols int) int {
	return row*numCols + col
}

// Expand a 1D index into 2D coordinates
func expandIdx(index, numCols int) (int, int) {
	row := index / numCols
	col := index % numCols
	return row, col
}

type Synapse struct {
	PostSynapticNeuron uint64 // Postsynaptic neuron
	Weight             int16  // W
}

type Neuron struct {
	Potential       uint16
	Threshold       uint16
	ResetPotential  uint16
	PotentialDecay  uint16
	Synapses        []*Synapse
	Position        [2]float64 // Position on the screen
	IsFiring        bool
	RefractoryTimer int
}

type Network struct {
	Neurons []*Neuron
}

func MakeNeuron(networkSize int) *Neuron {
	n := &Neuron{}
	n.InitNeuron(networkSize)
	return n
}

func (n *Neuron) InitNeuron(networkSize int) {
	n.Potential = 0
	n.Threshold = uint16(rand.Intn(100) + 50)      // Random threshold between 50 and 150
	n.ResetPotential = uint16(rand.Intn(50))       // Random reset potential between 0 and 50
	n.PotentialDecay = uint16(rand.Intn(5) + 1)    // Random decay between 1 and 5
	n.Synapses = MakeSynapses(networkSize)
	n.IsFiring = false
	n.RefractoryTimer = 0
}

func (n *Network) InitNetwork(networkSize int, xSize, ySize int) {
	for i := 0; i < networkSize; i++ {
		neuron := MakeNeuron(networkSize)
		n.Neurons = append(n.Neurons, neuron)
	}
	// Set positions for neurons
	for i, neuron := range n.Neurons {
		row, col := expandIdx(i, xSize)
		neuron.Position = [2]float64{
			float64(col)*(800/float64(xSize)) + 50,
			float64(row)*(600/float64(ySize)) + 50,
		}
	}
}

func MakeSynapses(networkSize int) []*Synapse {
	numSynapses := rand.Intn(3) + 1 // 1 to 3 synapses
	synapses := make([]*Synapse, numSynapses)
	for i := 0; i < numSynapses; i++ {
		synapses[i] = &Synapse{
			PostSynapticNeuron: uint64(rand.Intn(networkSize)),
			Weight:             int16(rand.Intn(50) + 1), // Weight between 1 and 50
		}
	}
	return synapses
}

type Game struct {
	Network   *Network
	TickCount int
}

func (g *Game) Update() error {
	g.TickCount++

	// Update neuron potentials and check for firing
	for _, neuron := range g.Network.Neurons {
		if neuron.RefractoryTimer > 0 {
			neuron.RefractoryTimer--
			continue
		}
		neuron.Potential += uint16(rand.Intn(3)) // Random input between 0 and 2
		neuron.Potential = uint16(int(neuron.Potential) - int(neuron.PotentialDecay))
		if neuron.Potential >= neuron.Threshold {
			neuron.IsFiring = true
			neuron.Potential = neuron.ResetPotential
			neuron.RefractoryTimer = 20 // Refractory period
			// Propagate spikes to postsynaptic neurons
			for _, synapse := range neuron.Synapses {
				postNeuron := g.Network.Neurons[synapse.PostSynapticNeuron]
				postNeuron.Potential += uint16(synapse.Weight)
			}
		} else {
			neuron.IsFiring = false
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Draw synapses
	for i, neuron := range g.Network.Neurons {
		for _, synapse := range neuron.Synapses {
			prePos := neuron.Position
			postNeuron := g.Network.Neurons[synapse.PostSynapticNeuron]
			postPos := postNeuron.Position
			color := color.RGBA{100, 100, 100, 255}
			ebitenutil.DrawLine(screen, prePos[0], prePos[1], postPos[0], postPos[1], color)
		}
	}

	// Draw neurons
	for _, neuron := range g.Network.Neurons {
		var clr color.Color
		if neuron.IsFiring {
			clr = color.RGBA{255, 0, 0, 255} // Red if firing
		} else {
			clr = color.RGBA{0, 255, 0, 255} // Green if not firing
		}
		x, y := neuron.Position[0], neuron.Position[1]
		radius := 5.0
		ebitenutil.DrawCircle(screen, x, y, radius, clr)
	}

	// Optionally, display debug info
	msg := fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 900, 700
}

func main() {
	rand.Seed(time.Now().UnixNano())

	x, y := 16, 12
	numNeurons := x * y

	network := &Network{}
	network.InitNetwork(numNeurons, x, y)

	game := &Game{
		Network: network,
	}

	ebiten.SetWindowSize(900, 700)
	ebiten.SetWindowTitle("Spiking Neural Network Visualization")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
