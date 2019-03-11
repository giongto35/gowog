import gym
import torch
import random
import collections
from torch.autograd import Variable

import numpy as np

from collections import namedtuple, deque

from agent import BaseAgent

# one single experience step
Experience = namedtuple('Experience', ['state', 'action', 'reward', 'done'])


class ExperienceSource:
    """
    Simple n-step experience source using single or multiple environments
    Every experience contains n list of Experience entries
    """
    def __init__(self, env, agent, steps_count=2, steps_delta=1):
        """
        Create simple experience source
        :param env: environment or list of environments to be used
        :param agent: callable to convert batch of states into actions to take
        :param steps_count: count of steps to track for every experience chain
        :param steps_delta: how many steps to do between experience items
        """
        # assert isinstance(env, (gym.Env, list, tuple))
        # assert isinstance(agent, BaseAgent)
        # assert isinstance(steps_count, int)
        assert steps_count >= 1
        self.agent = agent
        self.steps_count = steps_count
        self.steps_delta = steps_delta
        self.total_rewards = []
        self.total_steps = []
        self.env = env
        self.agent_states = agent.initial_state()

    def __iter__(self):
        states, histories, cur_rewards, cur_steps = [], [], 0.0, 0
        env = self.env
        states.append(env.reset()[0])
        histories.append(deque(maxlen=self.steps_count))

        iter_idx = 0
        while True:
            actions, self.agent_states = self.agent(states, self.agent_states)

            state = states[0]
            action = actions[0]
            history = histories[0]
            next_state, r, is_done, _ = env.step(action)
            cur_rewards += r
            cur_steps += 1
            history.append(Experience(state=state, action=action, reward=r, done=is_done))
            if len(history) == self.steps_count and iter_idx % self.steps_delta == 0:
                yield tuple(history)
            states[0] = next_state[0]
            if is_done:
                # generate tail of history
                while len(history) >= 1:
                    yield tuple(history)
                    history.popleft()
                self.total_rewards.append(cur_rewards)
                self.total_steps.append(cur_steps)
                cur_rewards = 0.0
                cur_steps = 0
                states[0] = env.reset()[0]
                self.agent_states[0] = self.agent.initial_state()
                history.clear()

            iter_idx += 1

    def pop_total_rewards(self):
        r = self.total_rewards
        if r:
            self.total_rewards = []
            self.total_steps = []
        return r

    def pop_rewards_steps(self):
        res = list(zip(self.total_rewards, self.total_steps))
        if res:
            self.total_rewards, self.total_steps = [], []
        return res


# those entries are emitted from ExperienceSourceFirstLast. Reward is discounted over the trajectory piece
ExperienceFirstLast = collections.namedtuple('ExperienceFirstLast', ('state', 'action', 'reward', 'last_state'))


class ExperienceSourceFirstLast(ExperienceSource):
    """
    This is a wrapper around ExperienceSource to prevent storing full trajectory in replay buffer when we need
    only first and last states. For every trajectory piece it calculates discounted reward and emits only first
    and last states and action taken in the first state.
    If we have partial trajectory at the end of episode, last_state will be None
    """
    def __init__(self, env, agent, gamma, steps_count=1, steps_delta=1):
        assert isinstance(gamma, float)
        super(ExperienceSourceFirstLast, self).__init__(env, agent, steps_count+1, steps_delta)
        self.gamma = gamma
        self.steps = steps_count

    def __iter__(self):
        for exp in super(ExperienceSourceFirstLast, self).__iter__():
            if exp[-1].done and len(exp) <= self.steps:
                last_state = None
                elems = exp
            else:
                last_state = exp[-1].state
                elems = exp[:-1]
            total_reward = 0.0
            for e in reversed(elems):
                total_reward *= self.gamma
                total_reward += e.reward
            yield ExperienceFirstLast(state=exp[0].state, action=exp[0].action,
                                      reward=total_reward, last_state=last_state)
