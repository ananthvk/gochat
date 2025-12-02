#!/bin/bash

http POST http://localhost:8000/api/v1/auth/signup name="User 1" email="user1@example.com" username="user1" password="password123"
http POST http://localhost:8000/api/v1/auth/signup name="User 2" email="user2@example.com" username="user2" password="password123"
http POST http://localhost:8000/api/v1/auth/signup name="User 3" email="user3@example.com" username="user3" password="password123"