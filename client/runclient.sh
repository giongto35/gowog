#!/bin/bash
npm run deploy -- --env.HOST_IP=$HOST_IP
http-server build -p 3000
