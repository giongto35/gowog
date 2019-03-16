# GOWOG AI

## Summary

This folder contains CS2D environment following openAI gym format and a training script. The trainig script uses NeuroEvolution to find the shortest path to the destination in the maze.

## Installation

Set a python3 Environment with virtualenv. Access to virtualenv

run gowog environment. i.e `./run_local.sh`

run training script. `python train_ga.py -n save_file_name`

Note: 
save_fie_name is where the weights is saved. The next time, if we specify the existing file, it will continue training from weights in the last run of that file.

## Credits

The code is based on "Deep Reinforcement Learning Hands-On" - Maxim Lapan
https://github.com/PacktPublishing/Deep-Reinforcement-Learning-Hands-On/blob/master/Chapter16/04_cheetah_ga.py
