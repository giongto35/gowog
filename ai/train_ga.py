"""
The code is inspired from "Deep Reinforcement Learning Hands-On" - Maxim Lapan
https://github.com/PacktPublishing/Deep-Reinforcement-Learning-Hands-On/blob/master/Chapter16/04_cheetah_ga.py
"""
#!/usr/bin/env python3
import os
import sys
import gym
import argparse
import collections
import copy
import time
import pickle
import numpy as np
import cs2denv_ga as cs2denv

import torch
import torch.nn as nn
import torch.multiprocessing as mp

# Noise to mutate the network
NOISE_STD = 0.05
# Size of population to perform mutation
POPULATION_SIZE = 200
# Number of parents for mutation
PARENTS_COUNT = 10
# Number of workers to be run parallely
WORKERS_COUNT = 6
# SEEDS
SEEDS_PER_WORKER = POPULATION_SIZE // WORKERS_COUNT
MAX_SEED = 2**32 - 1


class Net(nn.Module):
    """
    Return Policy network contains 3 layer
       obs_size * hidden_size
       hidden_size * hidden_size
       hidden_size * action_size

    The last layer is (dx, dy) which is the direction of agent should follow
    """
    def __init__(self, obs_size, act_size, hid_size=64):
        super(Net, self).__init__()

        self.mu = nn.Sequential(
            nn.Linear(obs_size, hid_size),
            nn.Tanh(),
            nn.Linear(hid_size, hid_size),
            nn.Tanh(),
            nn.Linear(hid_size, act_size),
            nn.Tanh(),
        )

    def forward(self, x):
        return self.mu(x)


def evaluate(env, net):
    """
    Return reward if agents follow current policy network
    """
    obs = env.reset()
    reward = -10000000.0
    steps = 0
    while True:
        obs_v = torch.FloatTensor(obs)
        action_v = net(obs_v)
        obs, r, done, _ = env.step(action_v.data.numpy()[0])
        reward = r
        steps += 1
        if done:
            break
    print("reward ", reward)
    return reward , steps


def mutate_net(net, seed, copy_net=True):
    """
    Return a new network from one seed. The network is mutated by adding NOISE_STD to all params

    Params:
      seed - each seed represents the set of mutated weights of the network
             because random will returns the same result for the same seed.
    """
    new_net = copy.deepcopy(net) if copy_net else net
    np.random.seed(seed)
    for p in new_net.parameters():
        noise_t = torch.tensor(np.random.normal(size=p.data.size()).astype(np.float32))
        p.data += NOISE_STD * noise_t
    return new_net


def build_net(env, seeds):
    """
    Return a new network from the list of seeds. The network is mutated from the seeds list.
    """
    torch.manual_seed(seeds[0])
    net = Net(env.observation_space.shape[1], env.action_space.n)
    for seed in seeds[1:]:
        net = mutate_net(net, seed, copy_net=False)
    return net


OutputItem = collections.namedtuple('OutputItem', field_names=['seeds', 'reward', 'steps'])

def worker_func(input_queue, output_queue):
    """
    Worker represents an agent. Workers receive the seeds from input_queue
    """
    env = cs2denv.CS2DEnvironment('local', 'Train Bot')
    cache = {}

    while True:
        parents = input_queue.get()
        if parents is None:
            break
        new_cache = {} # Cache list of seeds -> network
        for net_seeds in parents:
            if len(net_seeds) > 1:
                # Get if there is cache for previous
                net = cache.get(net_seeds[:-1])
                if net is not None:
                    # If there is, continue, just mutate the last one
                    net = mutate_net(net, net_seeds[-1])
                else:
                    # If not, build again from start
                    net = build_net(env, net_seeds)
            else:
                net = build_net(env, net_seeds)
            new_cache[net_seeds] = net
            # Evaluate the seeds
            reward, steps = evaluate(env, net)
            output_queue.put(OutputItem(seeds=net_seeds, reward=reward, steps=steps))
        cache = new_cache


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-n", "--name", required=False, help="Name of the file we will save the current progress")
    args = parser.parse_args()
    # the program accepts name, that is 
    save_path = os.path.join("save", args.name)
    fname = os.path.join(save_path, 'latest.dat')
    elitefname = os.path.join(save_path, 'elite.dat')

    mp.set_start_method('spawn')

    # generate population
    population = []
    # Continue from old training model
    if os.path.isdir(save_path):
        # Get the latest file
        f = open(fname,"rb")
        population = pickle.load(f)
        f.close()
    os.makedirs(save_path, exist_ok=True)

    # Spawn list of workers
    # Input is the list of input to workers
    input_queues = []
    # The result of an episode will be produced from workers
    output_queue = mp.Queue(maxsize=WORKERS_COUNT)
    workers = []
    for _ in range(WORKERS_COUNT):
        input_queue = mp.Queue(maxsize=1)
        input_queues.append(input_queue)
        w = mp.Process(target=worker_func, args=(input_queue, output_queue))
        w.start()
        # if population is not loaded (not continue from previous run)
        if len(population) == 0:
            seeds = [(np.random.randint(MAX_SEED),) for _ in range(SEEDS_PER_WORKER)]
            input_queue.put(seeds)

    gen_idx = 0 # monitor: number of generation
    # elite is the best run in the batch
    elite = None
    while True:
        t_start = time.time() # monitor: calculate running time

        batch_steps = 0
        while len(population) < SEEDS_PER_WORKER * WORKERS_COUNT:
            out_item = output_queue.get()
            population.append((out_item.seeds, out_item.reward))
            batch_steps += out_item.steps

        if elite is not None:
            population.append(elite)

        population.sort(key=lambda p: p[1], reverse=True)
        # Save population to the predefined file from flags
        record_latest_file = open(fname,"wb")
        record_elite_file = open(elitefname,"wb")
        pickle.dump(population, record_latest_file)
        pickle.dump(population[0], record_elite_file)
        record_latest_file.close()
        record_elite_file.close()

        # monitor: calculate rewards from batch
        rewards = [p[1] for p in population[:PARENTS_COUNT]]
        reward_mean = np.mean(rewards)
        reward_max = np.max(rewards)
        reward_std = np.std(rewards)
        speed = batch_steps / (time.time() - t_start)
        print("%d: reward_mean=%.2f, reward_max=%.2f, reward_std=%.2f, speed=%.2f f/s" % (
            gen_idx, reward_mean, reward_max, reward_std, speed))

        # Elite the best generation from population
        elite = population[0]
        for input_queue in input_queues:
            seeds = []
            for _ in range(SEEDS_PER_WORKER):
                parent = np.random.randint(PARENTS_COUNT)
                next_seed = np.random.randint(MAX_SEED)
                # each worker get a list of seeds
                # list of seeds = previous seeds + next seed
                seeds.append(tuple(list(population[parent][0]) + [next_seed]))
            input_queue .put(seeds)
        gen_idx += 1
        population = []

    pass
