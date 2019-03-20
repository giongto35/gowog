#!/usr/bin/env python3
import cs2denv
import random

agent = cs2denv.CS2DEnvironment('prod', 'Train Bot')
obs = agent.reset()

while True:
    agent.move_position( random.randint(-1, 1), random.randint(-1, 1))
    agent.shoot(random.uniform(-2, 2), random.uniform(-2, 2))

