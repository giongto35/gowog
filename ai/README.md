# GOWOG AI

## Summary

This repo contains CS2D environment following openAI gym format and a training script. The training script uses NeuroEvolution to find the shortest path to the destination in the maze.

## Installation

run gowog environment as instruction from . i.e install Docker and `./run_local.sh`

set up a python3 virtual environment with virtualenv. Run `requirements.txt`

run training script. `python train_ga.py -n save_file_name`

save_fie_name is where the weights are saved. In the next time, if we specify an existing file, it will continue training from weights in the last run of that file.

## Genetic Algorithm

### Agent
_Implementation at cs2denv_ga.py_

CS2D Agent is built on CS2D for machine learning purpose. It follows openAI gym by supporting fundamental methods of an agent including : reset(), step(), observation_space and action_space.

The ObservationSpace is a 1D array constructed from the update_player message from server  
1. Player position, player size, number of columns, number of rows, block width, block height
2. The distance to nearest block to the left, right, up, down. This input is to avoid collision
3. The player position in the binary block map. The map is a 2D array of 0, 1 when 0 is empty and 1 is block.
4. The binary block map.

The Reward is the 1 / distance to the goal. If the agent is close to the goal by 100, the reward is 1 and the episode finishes.

### NeuroEvolution
_Implementation at train_ga.py_

Neural Network is used to get the best action by passing the input (observation space) through the neural network.

NeuroEvolution is an AI that uses evolutionary algorithms to continuously improve artificial neural network. For each iteration (generation), the program will generate a new set of neural network weights based on best settings in the previous iteration. The process of generating a NN from previous NN called *Mutate*, which added random noise to each params in the NN.  

One special enhancement is that instead of storing all weights of a generation, we store only list of the noise seeds applied to the neural network. Because under the same seed, all the randomization is the same, so a seed can represent a mutation operator of a network. Instead of keeping all the weights of a generation, we can just store a set of seeds from begin to the current generation, and re-construct the weights from that set to get the weights of all the neural networks.

## Credits

The code is based on "Deep Reinforcement Learning Hands-On" - Maxim Lapan
https://github.com/PacktPublishing/Deep-Reinforcement-Learning-Hands-On/blob/master/Chapter16/04_cheetah_ga.py
