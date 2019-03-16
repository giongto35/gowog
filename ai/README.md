# GOWOG AI

## Summary

This folder contains CS2D environment following openAI gym format and a training script. The trainig script uses NeuroEvolution to find the shortest path to the destination in the maze.

## Installation

Set a python3 Environment with virtualenv. Access to virtualenv

run gowog environment. i.e `./run_local.sh`

run training script. `python train_ga.py -n save_file_name`

save_fie_name is where the weights is saved. The next time, if we specify the existing file, it will continue training from weights in the last run of that file.

## Genetic Algorithm

### Agent

`cs2denv_ga.py`

CS2D Agent is built on CS2D following openAI agent. It supports the fundamental methods of OpenAI Agent including : reset(), step()

The ObservationSpace is 1D array constructed from the update_player message from server
1. Player position, player size, number of columns, number of rows, block width, block height
2. The nearest block to the left, right, up, down
3. The player position in the binary block map. The map is a 2D array of 0, 1 when 0 is empty and 1 is block.
4. The block map in 0, 1

The Reward is the 1 / distance to the goal. If the agent is close to the goal by 100, the reward is 1 and the episode finishes.

### NeuroEvolution

NeuroEvolution is an AI that uses evolutionary algorithms to generate artificial neural network. For each iteration (generation), the program will generate a new set of neural network weights based on the best settings in the previous iteration. The process of generating a NN from previous NN called *Mutate*, which added a random noise to each params in the NN.

One special enhancement is that instead of storing all weights of a generation, we store only list of the noise seeds applied to the neural network. Because under the same seed, all the randomization is the same, so a seed can represent a mutation operator of a network. Instead of keeping all the weights of a generation, we can just store a set of seeds from begin to the current generation, and re-construct the weights from that set to get the weights of all the neural networks.

## Credits

The code is based on "Deep Reinforcement Learning Hands-On" - Maxim Lapan
https://github.com/PacktPublishing/Deep-Reinforcement-Learning-Hands-On/blob/master/Chapter16/04_cheetah_ga.py
