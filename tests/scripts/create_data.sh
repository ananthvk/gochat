#!/bin/bash
set -e
# Create dummy data that can be used to check if the frontend displays it correctly
# This script should be run after calling create_users.sh

# Create three groups

# User 1 & 2 belong to a group (1 is the owner)
# User 2 & 3 belong to a group (2 is the owner)
# User 3 & 1 belong to a group (3 is the owner)

# Apart from these, each user has 100 groups with only them as the member

API_BASE_URL="http://localhost:8000/api/v1"

echo "Getting auth tokens"
AUTH_USER_1=$(http POST "$API_BASE_URL/auth/login" email='user1@example.com' password='password123' | jq -r '.authenticate.token')
AUTH_USER_2=$(http POST "$API_BASE_URL/auth/login" email='user2@example.com' password='password123' | jq -r '.authenticate.token')
AUTH_USER_3=$(http POST "$API_BASE_URL/auth/login" email='user3@example.com' password='password123' | jq -r '.authenticate.token')
echo "$AUTH_USER_1"
echo "$AUTH_USER_2"
echo "$AUTH_USER_3"

# Create the main three groups
echo "Creating groups"
GROUP_U12=$(http POST "$API_BASE_URL/group" name="U12 Group" description="Group with U1 & U2" Authorization:"Bearer $AUTH_USER_1" | jq -r '.id')
GROUP_U23=$(http POST "$API_BASE_URL/group" name="U23 Group" description="Group with U2 & U3" Authorization:"Bearer $AUTH_USER_2" | jq -r '.id')
GROUP_U31=$(http POST "$API_BASE_URL/group" name="U31 Group" description="Group with U3 & U1" Authorization:"Bearer $AUTH_USER_3" | jq -r '.id')
GROUP_U123=$(http POST "$API_BASE_URL/group" name="U123 Group" description="Group with U1, U2, and U3" Authorization:"Bearer $AUTH_USER_1" | jq -r '.id')
echo "U12: $GROUP_U12"
echo "U23: $GROUP_U23"
echo "U31: $GROUP_U31"
echo "U123: $GROUP_U123"

# Create the other 10 groups
echo "Creating other 10 groups"
for i in {1..10}; do
    X=$(http POST "$API_BASE_URL/group" name="U1-$i" description="U1's $i th group" Authorization:"Bearer $AUTH_USER_1" | jq '.id')
    X=$(http POST "$API_BASE_URL/group" name="U2-$i" description="U2's $i th group" Authorization:"Bearer $AUTH_USER_2" | jq '.id')
    X=$(http POST "$API_BASE_URL/group" name="U3-$i" description="U3's $i th group" Authorization:"Bearer $AUTH_USER_3" | jq '.id')
done

# Join the groups
http PUT "$API_BASE_URL/group/$GROUP_U12/member" Authorization:"Bearer $AUTH_USER_2"
http PUT "$API_BASE_URL/group/$GROUP_U23/member" Authorization:"Bearer $AUTH_USER_3"
http PUT "$API_BASE_URL/group/$GROUP_U31/member" Authorization:"Bearer $AUTH_USER_1"

http PUT "$API_BASE_URL/group/$GROUP_U123/member" Authorization:"Bearer $AUTH_USER_2"
http PUT "$API_BASE_URL/group/$GROUP_U123/member" Authorization:"Bearer $AUTH_USER_3"

echo "Creating messages"
# Create dummy messages in each group
counter=1
for i in {1..100}; do
    X=$(http POST "$API_BASE_URL/group/$GROUP_U123/message" type="text" content="User 1's important message $i ($counter)" Authorization:"Bearer $AUTH_USER_1")
    counter=$((counter+1))
    X=$(http POST "$API_BASE_URL/group/$GROUP_U123/message" type="text" content="User 2's important message $i ($counter)" Authorization:"Bearer $AUTH_USER_2")
    counter=$((counter+1))
    X=$(http POST "$API_BASE_URL/group/$GROUP_U123/message" type="text" content="User 3's important message $i ($counter)" Authorization:"Bearer $AUTH_USER_3")
    counter=$((counter+1))
done
echo "Created for U123"


counter=1
for i in {1..100}; do
    X=$(http POST "$API_BASE_URL/group/$GROUP_U12/message" type="text" content="U1 sent this in U12, $i ($counter)" Authorization:"Bearer $AUTH_USER_1")
    counter=$((counter+1))
    X=$(http POST "$API_BASE_URL/group/$GROUP_U12/message" type="text" content="U2 sent this in U12, $i ($counter)" Authorization:"Bearer $AUTH_USER_2")
    counter=$((counter+1))
done
echo "Created for U12"

counter=1
for i in {1..100}; do
    X=$(http POST "$API_BASE_URL/group/$GROUP_U23/message" type="text" content="U2 sent this in U23, $i ($counter)" Authorization:"Bearer $AUTH_USER_2")
    counter=$((counter+1))
    X=$(http POST "$API_BASE_URL/group/$GROUP_U23/message" type="text" content="U3 sent this in U23, $i ($counter)" Authorization:"Bearer $AUTH_USER_3")
    counter=$((counter+1))
done
echo "Created for U23"

counter=1
for i in {1..100}; do
    X=$(http POST "$API_BASE_URL/group/$GROUP_U12/message" type="text" content="U3 sent this in U31, $i ($counter)" Authorization:"Bearer $AUTH_USER_3")
    counter=$((counter+1))
    X=$(http POST "$API_BASE_URL/group/$GROUP_U12/message" type="text" content="U1 sent this in U31, $i ($counter)" Authorization:"Bearer $AUTH_USER_1")
    counter=$((counter+1))
done
echo "Created for U31"