from websocket import create_connection
import random
import math
import message_pb2 as messagepb
import numpy as np

OBJECTIVE_NOT_COLLIDE = "not_collide"
OBJECTIVE_NOT_HIT = "not_hit"

LOCAL_ENV = 'local'
MAX_STEPS = 200
PLAYER_SIZE = 32
BLOCK_SIZE = 64
EPS = 5
WIN_REWARD = 30 # = 1 / distance to the nearest point = 1/ 2
GOAL = {"x": BLOCK_SIZE * 11, "y": BLOCK_SIZE * 5}

dx = [-1, 0, 1, 0]
dy = [0, -1, 0, 1]

class ObservationSpace:
    def __init__(self, shape):
        self.shape = shape

class ActionSpace:
    def __init__(self, shape):
        self.n = 4
        self.sample = random.randint(0, self.n)

class CS2DEnvironment:
    def __init__(self, env, name):
        if env == LOCAL_ENV:
            wshost = 'ws://localhost:8080/game/'
        else:
            wshost = 'ws://35.198.245.53/game/'

        self.name = name
        self.wshost = wshost
        self.init()
        self.observation_space = ObservationSpace(self.get_obs_size())
        self.action_space = ActionSpace(4)
        self.speed = 1
        self.objective = OBJECTIVE_NOT_COLLIDE

    def step(self, d):
        dx, dy = d[0], d[1]
        obs, reward, done = self.move_position(dx, dy)
        return obs, reward, done, None

    def init_player(self, name, client_id):
        self.current_input_number += 1
        message = messagepb.ClientGameMessage()

        init_player = messagepb.InitPlayer()
        init_player.name = name
        init_player.client_id = client_id
        message.input_sequence_number = self.current_input_number
        message.init_player_payload.CopyFrom(init_player)

        self.ws.send(message.SerializeToString())

    def set_position(self, x, y):
        message = messagepb.ClientGameMessage()

        set_position = messagepb.SetPosition()
        set_position.id = self.player.id
        set_position.x = x
        set_position.y = y
        message.set_position_payload.CopyFrom(set_position)
        message.input_sequence_number = self.current_input_number

        self.ws.send(message.SerializeToString())

        self.player.x = x
        self.player.y = y

    def move_position(self, dx, dy):
        self.num_steps += 1
        self.current_input_number += 1
        message = messagepb.ClientGameMessage()

        # construct message
        message.time_elapsed = 0.1
        move_position = messagepb.MovePosition()
        move_position.id = self.player.id
        move_position.dx = self.speed * dx
        move_position.dy = self.speed * dy
        message.move_position_payload.CopyFrom(move_position)
        message.input_sequence_number = self.current_input_number

        self.ws.send(message.SerializeToString())
        while True:
            # Set timeout
            binary_res = self.ws.recv()
            # Received client accepted from server with client_id
            serverMsg = messagepb.ServerGameMessage()
            serverMsg.ParseFromString(binary_res)

            if serverMsg.HasField("remove_player_payload") == True:
                remove_player_msg = serverMsg.remove_player_payload
                if self.player.id == remove_player_msg.id:
                    obs = self.__get_obs()
                    return obs, 0, True

            if serverMsg.HasField("update_player_payload") == True:
                update_player_msg = serverMsg.update_player_payload
                last_process_input = update_player_msg.current_input_number

                if self.player.id == update_player_msg.id \
                        and last_process_input == self.current_input_number \
                        and last_process_input > self.current_server_number:
                    # updated_reward = update_player_msg.ai_reward

                    # is_done = False
                    # if self.objective == OBJECTIVE_NOT_COLLIDE:
                        # is_done = update_player_msg.x == self.player.x and update_player_msg.y == self.player.y
                    # # if is_done == True:
                        # # updated_reward = 0 

                    # self.current_server_number = last_process_input
                    # self.player = update_player_msg
                    # obs = self.__get_obs()

                    # return obs, updated_reward, is_done
                    dist = self.dist(update_player_msg.x, update_player_msg.y, GOAL['x'], GOAL['y'])
                    ai_reward = 1 / dist
                    updated_reward = ai_reward

                    is_done = False

                    # if self.num_steps >= MAX_STEPS:
                        # updated_reward = ai_reward
                        # is_done = True
                    if self.num_steps >= MAX_STEPS or abs(update_player_msg.x - self.player.x) + abs(update_player_msg.y - self.player.y) <= EPS:
                        updated_reward = ai_reward
                        is_done = True
                    if dist <= WIN_REWARD:
                        updated_reward = 1
                        is_done = True
                        print('reward : ', updated_reward)

                    self.current_server_number = last_process_input
                    self.player = update_player_msg
                    obs = self.__get_obs()

                    return obs, updated_reward, is_done

    def dist(self, x1, y1, x2, y2):
        return math.sqrt((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2))

    def normalize(self, dx, dy):
        l = math.sqrt(dx * dx + dy * dy)
        return dx / l, dy / l

    def shoot(self, dx, dy):
        message = messagepb.ClientGameMessage()

        shoot = messagepb.Shoot()
        shoot.dx, shoot.dy = self.normalize(dx, dy)
        # shoot.dx, shoot.dy = dx, dy
        shoot.x = self.player.x
        shoot.y = self.player.y
        shoot.player_id = self.player.id
        message.shoot_payload.CopyFrom(shoot)

        self.ws.send(message.SerializeToString())

    def init(self):
        self.player = None
        print('create connection')
        self.ws = create_connection(self.wshost)
        print('create connection done')

        # rem_steps define number of steps left
        self.rem_steps = MAX_STEPS
        self.current_input_number = 0
        self.current_server_number = -1

        while True:
            # Set timeout
            binary_res = self.ws.recv()

            # Received client accepted from server with client_id
            serverMsg = messagepb.ServerGameMessage()
            serverMsg.ParseFromString(binary_res)

            if serverMsg.HasField("register_client_id_payload") == True:
                self.client_id = serverMsg.register_client_id_payload.client_id
                self.init_player(self.name, self.client_id)
            if serverMsg.HasField("init_player_payload") == True:
                if serverMsg.init_player_payload.is_main:
                    self.player = serverMsg.init_player_payload
            if serverMsg.HasField("init_all_payload") == True:
                self.map = serverMsg.init_all_payload.init_map.block
                self.map_ncols = serverMsg.init_all_payload.init_map.num_cols
                self.map_nrows = serverMsg.init_all_payload.init_map.num_rows
                self.map_bwidth = serverMsg.init_all_payload.init_map.block_width
                self.map_bheight = serverMsg.init_all_payload.init_map.block_height

            if self.player != None and self.map != None:
                break

        return self.player, self.map



        # self.client_id = res.client_id
        # self.init_player(name, self.client_id)

        # while True:
            # binary_res = self.ws.recv()
            # res = messagepb.ServerGameMessage()
            # res.ParseFromString(binary_res)
            # res = res.init_all_payload
            # for player in res.init_player:
                # if player.is_main:
                    # self.player = player
            # print("Myself is : ", self.player)


    def reset(self):
        self.num_steps = 0
        self.ws.close()
        self.init()
        self.set_position(100, 100)
        return self.__get_obs()
        # while True:
            # # Set timeout
            # binary_res = self.ws.recv()
            # # Received client accepted from server with client_id
            # serverMsg = messagepb.ServerGameMessage()
            # serverMsg.ParseFromString(binary_res)

            # if serverMsg.HasField("update_player_payload") == True:
                # update_player_msg = serverMsg.update_player_payload

                # if self.player.id == update_player_msg.id:
                    # return self.__get_obs(update_player_msg)

    def __get_nearest_dist(self, x, y):
        block_size = self.map_bwidth

        i = int(y / block_size)
        j = int(x / block_size)
        left = x
        right = self.map_ncols * block_size - x
        up = y
        down = self.map_nrows * block_size - y

        for jj in range(j, -1, -1):
            if self.map[i * self.map_ncols + jj] == 1:
                left = max(0, x - (jj + 1) * block_size)
                break

        for jj in range(j, self.map_ncols):
            if self.map[i * self.map_ncols + jj] == 1:
                right = max(0, jj * block_size - x)
                break

        for ii in range(i, -1, -1):
            if self.map[ii * self.map_ncols + j] == 1:
                up = max(0, y - (ii + 1) * block_size)
                break

        for ii in range(i, self.map_nrows):
            if self.map[ii * self.map_ncols + j] == 1:
                down = max(0, ii * block_size - y)
                break

        return (left, right, up, down)

    def __get_player_position(self, x, y):
        block_size = self.map_bwidth

        i = int(y / block_size)
        j = int(x / block_size)
        p = i * self.map_ncols + j
        arr = [0] * (self.map_ncols * self.map_nrows)
        arr[p] = 1

        return arr


    def __get_obs(self):
        # from player and map, we build observation
        # The first observation is player
        obs = []
        obs.extend([self.player.x, self.player.y])
        # Add game constant
        obs.extend([PLAYER_SIZE, self.map_ncols, self.map_nrows, self.map_bwidth, self.map_bheight])
        # Add distance to nearest block
        (left, right, up, down) = self.__get_nearest_dist(self.player.x, self.player.y)
        obs.extend([left, right, up, down])
        # Add player position
        player_pos_mat = self.__get_player_position(self.player.x, self.player.y)
        obs.extend(player_pos_mat)
        # The last observation is map
        obs.extend(self.map)
        obs = [obs]
        return np.array(obs)
        # The next observation is map

    def get_obs_size(self):
        return self.__get_obs().shape

    # def obs():
        # observation, reward, done, info = env.step(action)
        # return observation

    # def reset():

