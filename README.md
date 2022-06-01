# _This repo is archived. See http://github.com/waggle-sensor/edge-scheduler_

# sage-ses
SAGE Edge Scheduler

# SAGE SES API
This is a mock API for the Sage Edge Scheduler (SES).

# Authentication
SAGE users authenticate via token they can get from the SAGE website.

example:
```bash
-H "Authorization: sage <sage_user_token>"
```

for test environment:
```bash
-H "Authorization: sage user:test"
```

# Start-up

```bash
docker-compose up
```
This starts a test environment without token verification.

# CURL to GET and POST

To test the API:

```bash
export SES_URL="localhost:8080"
export SAGE_USER_TOKEN=user:testuser
```

## /api/v1/metrics (to be completed)
To get stats of SES:

```bash
curl -X GET "${SES_URL}/api/v1/metrics"  -H "Authorization: sage ${SAGE_USER_TOKEN}"
```

## /api/v1/goals

To get a list of goals:

```bash
curl -X GET "${SES_URL}/api/v1/goals"  -H "Authorization: sage ${SAGE_USER_TOKEN}"
```

To take a goal to be added:

```bash
curl -X POST "${SES_URL}/api/v1/goals?name=mygoal"  -H "Authorization: sage ${SAGE_USER_TOKEN}"
```

example response:

```bash
{
  "id": "0000000001",
  "name": "mygoal",
  "owner": "testuser"
}
```

to do the following:
```bash
export GOAL_ID=0000000001
```

## /api/v1/goals/{id}/status

to get the current status of a goal:
```bash
curl -X GET "${SES_URL}/api/v1/goals/${GOAL_ID}/status"  -H "Authorization: sage ${SAGE_USER_TOKEN}"
```

to take a status of a goal to be added (e.g., SUBMITTED, SUSPENDED, SCHEDULED, ACTIVATED, DONE): 
```bash
curl -X POST "${SES_URL}/api/v1/goals/${GOAL_ID}/status?status=scheduled"  -H "Authorization: sage ${SAGE_USER_TOKEN}"
```

## /api/v1/goals/metrics

to get stats of goals:
```bash
curl -X GET "${SES_URL}/api/v1/goals/metrics"  -H "Authorization: sage ${SAGE_USER_TOKEN}"
```
