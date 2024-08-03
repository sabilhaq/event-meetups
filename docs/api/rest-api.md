# REST API

This document describes specification for REST API endpoint used by the system.

List of available endpoints:

- [REST API](#rest-api)
  - [Login](#login)
  - [List Events](#list-events)
  - [List Venues](#list-venues)
  - [Get Venue](#get-venue)
  - [Create Meetup](#create-meetup)
  - [List Meetups](#list-meetups)
  - [Get Meetup Info](#get-meetup-info)
  - [Update Meetup](#update-meetup)
  - [Cancel Meetup](#cancel-meetup)
  - [Join Meetup](#join-meetup)
  - [Leave Meetup](#leave-meetup)
  - [List Incoming Meetups](#list-incoming-meetups)

Beside these endpoints, client need to handle properly [Common Errors](./common-errors.md) that may occurs throughout all specified REST API endpoints in the system.

> **Note:**
>
> To simplify development, we will use predefined list of users, events, and venues which can be found in [this directory](./data/).

## Login

POST: `/session`

This endpoint is used to login to the system. It returns a JWT token that can be used to access other endpoints named `access_token`.

The decoded access token will contains at least the following information:

```json
{
  "id": 1,
  "username": "marion",
  "email": "marion@eveners.com"
}
```

The JWT secret will be `event-maker`, and the access token will be valid for 1 year.

For the sake of simplicity, we already have a predefined list of users that can be used for login in [here](./data/users.json).

**Example Request:**

```json
POST /session
Content-Type: application/json

{
  "username": "marion",
  "password": "123456"
}
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "ok": true,
  "data": {
    "id": 1,
    "username": "marion",
    "email": "marion@eveners.com",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJtYXJpb24iLCJpYXQiOjE3MzY1NjgwMDB9.8LAqGhaNwKx2JMr2piPCLTkkg1J0wbDgO4AEsExIw2M",
  },
  "ts": 1704954526
}
```

**Error Response:**

- Invalid Credentials

  ```json
  HTTP/1.1 401 Unauthorized
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_INVALID_CREDS",
    "msg": "Invalid username or password",
    "ts": 1704954526
  }
  ```

[Back to Top](#rest-api)

---

## List Events

GET: `/events`

This endpoint is used to list all events supported by the system.

The list of supported events can be found in [here](./data/events.json).

**Headers:**

- `Authorization` => The value is `Bearer {access_token}`.

**Example Request:**

```bash
GET /events
Authorization: Bearer {access_token}
```

**Success Response:**

  ```json
  HTTP/1.1 200 OK
  Content-Type: application/json

  {
    "ok": true,
    "data": [
      {
        "id": 1,
        "name": "Wedding"
      },
      {
        "id": 2,
        "name": "Exhibition"
      },
      {
        "id": 3,
        "name": "Bazaar"
      },
      {
        "id": 4,
        "name": "Workshop"
      }
    ],
    "ts": 1704954526
  }
  ```

**Error Response:**

No specific error response.

[Back to Top](#rest-api)

---

## List Venues

GET: `/venues`

This endpoint is used to list all venues available in the system.

For simplicity, we will have pre-defined list of venues, and the list cannot be changed. The list can be found in [here](./data/venues.json).

Field `open_days` is an array of integers, where each integer represents a day of the week. The value of the integer is as follows:

- `0` => Sunday
- `1` => Monday
- `2` => Tuesday
- `3` => Wednesday
- `4` => Thursday
- `5` => Friday
- `6` => Saturday

Field `open_at` and `closed_at` is string that represent time of the day in its timezone when the venue can start to accomodate meetups until it closed. The value is between `"00:00"` and `"23:59"`. For this project, we do not consider venues that open past midnight.

**Headers:**

- `Authorization` => The value is `Bearer {access_token}`.

**Query Params:**

- `event_id` => Filter venues that support the specified event id.
- `meetup_start_ts` => Filter venues that can accomodate meetups that start at the specified timestamp. The value is unix timestamp in seconds.
- `meetup_end_ts` => Filter venues that can accomodate meetups that end at the specified timestamp. The value is unix timestamp in seconds.

**Example Request:**

```bash
GET /venues?event_id=1&meetup_start_ts=1704938400&meetup_end_ts=1704945600
Authorization: Bearer {access_token}
```

**Success Response:**

  ```json
  HTTP/1.1 200 OK
  Content-Type: application/json

  {
    "ok": true,
    "data": [
      {
        "id": 1,
        "name": "Si Jalak Harupat",
        "open_days": [0, 1, 2, 3, 4, 5, 6],
        "open_at": "8:00",
        "closed_at": "23:59",
        "timezone": "Asia/Jakarta",
        "supported_events": [
          {
            "id": 1,
            "name": "Wedding",
            "meetups_capacity": 2
          },
          {
            "id": 2,
            "name": "Exhibition",
            "meetups_capacity": 2
          }
        ]
      },
      {
        "id": 3,
        "name": "Ice BSD",
        "open_days": [0, 1, 2, 3, 4, 5, 6],
        "open_at": "05:00",
        "closed_at": "23:59",
        "timezone": "Asia/Jakarta",
        "supported_events": [
          {
            "id": 1,
            "name": "Wedding",
            "meetups_capacity": 2
          },
          {
            "id": 2,
            "name": "Basketball",
            "meetups_capacity": 1
          },
          {
            "id": 3,
            "name": "Bazaar",
            "meetups_capacity": 4
          }
          {
            "id": 4,
            "name": "Tennis",
            "meetups_capacity": 2
          }
        ]
      }
    ],
    "ts": 1704954526
  }
  ```

**Error Response:**

No specific error response.

[Back to Top](#rest-api)

---

## Get Venue

GET: `/venues/{venue_id}`

This endpoint returns information about a single venue.

**Headers:**

- `Authorization` => The value is `Bearer {access_token}`.

**Example Request:**

```bash
GET /venues/1
Authorization: Bearer {access_token}
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "ok": true,
  "data": {
    "id": 1,
    "name": "Si Jalak Harupat",
    "open_days": [0, 1, 2, 3, 4, 5, 6],
    "open_at": "8:00",
    "closed_at": "23:59",
    "timezone": "Asia/Jakarta",
    "supported_events": [
      {
        "id": 1,
        "name": "Wedding",
        "meetups_capacity": 2
      },
      {
        "id": 2,
        "name": "Basketball",
        "meetups_capacity": 2
      }
    ]
  },
  "ts": 1704954526
}
```

**Error Response:**

No specific error response.

[Back to Top](#rest-api)

---

## Create Meetup

POST: `/meetups`

This endpoint is used to add a new meetup to the system. Meetup is a gathering event in a venue in a specific range of time.

Meetup can only be created in a venue that supports the event, within the operating hours of the venue, and not exceeding the capacity of the venue in that time.

User that create a meetup NOT automatically joined to the meetup. He/she needs to join the meetup by using [Join Meetup](#join-meetup) endpoint.

**Headers:**

- `Authorization` => The value is `Bearer {access_token}`.

**Example Request:**

```json
POST /meetups
Authorization: Bearer {access_token}

{
  "name": "Wedding Fulan",
  "venue_id": 1,
  "event_id": 1,
  "start_ts": 1704938400,
  "end_ts": 1704945600,
  "max_persons": 12
}
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "ok": true,
  "data": {
    "id": 1,
    "name": "Wedding Fulan",
    "venue": {
      "id": 1,
      "name": "Si Jalak Harupat"
    },
    "event": {
      "id": 1,
      "name": "Wedding"
    },
    "start_ts": 1704938400,
    "end_ts": 1704945600,
    "max_persons": 12,
    "organizer": {
      "id": 1,
      "username": "marion",
      "email": "marion@eveners.com"
    },
    "joined_persons": [],
    "joined_persons_count": 0,
    "is_joined": false,
    "status": "open"
  },
  "ts": 1704954526
}
```

**Error Response:**

- Invalid event

  ```json
  HTTP/1.1 400 Bad Request
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_INVALID_EVENT",
    "msg": "Event is not supported by the venue",
    "ts": 1704954526
  }
  ```

- Exceed venue capacity

  ```json
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_EXCEED_VENUE_CAPACITY",
    "msg": "Venue capacity is full on the designated meetup time",
    "ts": 1704954526
  }
  ```

- Not within venue operating hours

  ```json
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_VENUE_IS_CLOSED",
    "msg": "Venue is closed on the designated meetup time",
    "ts": 1704954526
  }
  ```

[Back to Top](#rest-api)

---

## List Meetups

GET: `/meetups`

This endpoint is used to list future meetups that is still open. The result is sorted from the nearest meetup time to the furthest.

This endpoint cannot be used to check whether a user has joined a meetup or not. To check it, use [Get Meetup Info](#get-meetup-info) or [List Incoming Meetups](#list-incoming-meetups).

**Headers:**

- `Authorization` => The value is `Bearer {access_token}`.

**Query Params:**

- `limit`, _OPTIONAL_ => Limit the maximum number of meetups to be returned. Default is `100` & maximum is `1000`.
- `event_id`, _OPTIONAL_ => Filter meetups that only support the specified event id.

**Example Request:**

```bash
GET /meetups?limit=5&event_id=1
Authorization: Bearer {access_token}
```

**Success Response:**

  ```json
  HTTP/1.1 200 OK
  Content-Type: application/json

  {
    "ok": true,
    "data": [
      {
        "id": 1,
        "name": "Wedding Fulan",
        "venue": {
          "id": 1,
          "name": "Si Jalak Harupat"
        },
        "event": {
          "id": 1,
          "name": "Wedding"
        },
        "start_ts": 1704938400,
        "end_ts": 1704945600,
        "max_persons": 12,
        "organizer": {
          "id": 1,
          "username": "marion",
          "email": "marion@eveners.com"
        },
        "joined_persons_count": 6,
        "status": "open"
      },
      {
        "id": 2,
        "name": "BBW",
        "venue": {
          "id": 2,
          "name": "Parahyangan Convention"
        },
        "event": {
          "id": 3,
          "name": "Bazaar"
        },
        "start_ts": 1704938400,
        "end_ts": 1704945600,
        "max_persons": 4,
        "organizer": {
          "id": 1,
          "username": "marion",
          "email": "marion@eveners.com"
        },
        "joined_persons_count": 4,
        "status": "open"
      }
    ],
    "ts": 1704954526
  }
  ```

**Error Response:**

No specific error response.

[Back to Top](#rest-api)

---

## Get Meetup Info

GET: `/meetups/{meetup_id}`

This endpoint is used to get an information about a single meetup.

Field `joined_persons`, `cancelled_reason` and `cancelled_at` only appears if the user is the organizer or joined person of the meetup.

If the `meetup_id` is a cancelled meetup, the `cancelled_reason` field will be filled with the reason why the meetup is cancelled.

Status of the meetup can be one of the following:

- `open`: meetup is still open for joining
- `closed`: meetup is closed for joining because it has reached the maximum number of persons
- `cancelled`: meetup is cancelled by the organizer
- `finished`: meetup is past its end time without being cancelled

**Headers:**

- `Authorization` => The value is `Bearer {access_token}`.

**Example Request:**

```bash
GET /meetups/1
Authorization: Bearer {access_token}
```

**Success Response:**

- Response to organizer or participant for non-cancelled meetup

  ```json
  HTTP/1.1 200 OK
  Content-Type: application/json

  {
    "ok": true,
    "data": {
      "id": 1,
      "name": "Wedding Fulan",
      "venue": {
        "id": 1,
        "name": "Si Jalak Harupat"
      },
      "event": {
        "id": 1,
        "name": "Wedding"
      },
      "start_ts": 1704938400,
      "end_ts": 1704945600,
      "max_persons": 12,
      "organizer": {
        "id": 1,
        "username": "marion",
        "email": "marion@eveners.com"
      },
      "joined_persons": [
        {
          "id": 1,
          "username": "marion",
          "email": "marion@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 2,
          "username": "todd",
          "email": "todd@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 3,
          "username": "anthony",
          "email": "anthony@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 4,
          "username": "kevin",
          "email": "kevin@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 5,
          "username": "eric",
          "email": "eric@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 6,
          "username": "elnora",
          "email": "elnora@eveners.com",
          "joined_at": 1704954526
        }
      ],
      "joined_persons_count": 6,
      "is_joined": true,
      "status": "open"
    },
    "ts": 1704954526
  }
  ```

- Response to organizer or participant for cancelled meetup

  ```json
  HTTP/1.1 200 OK
  Content-Type: application/json

  {
    "ok": true,
    "data": {
      "id": 1,
      "name": "Exhibition Gadget",
      "venue": {
        "id": 1,
        "name": "Si Jalak Harupat"
      },
      "event": {
        "id": 2,
        "name": "Exhibition"
      },
      "start_ts": 1704938400,
      "end_ts": 1704945600,
      "max_persons": 12,
      "organizer": {
        "id": 1,
        "username": "marion",
        "email": "marion@eveners.com"
      },
      "joined_persons": [
        {
          "id": 1,
          "username": "marion",
          "email": "marion@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 2,
          "username": "todd",
          "email": "todd@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 3,
          "username": "anthony",
          "email": "anthony@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 4,
          "username": "kevin",
          "email": "kevin@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 5,
          "username": "eric",
          "email": "eric@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 6,
          "username": "elnora",
          "email": "elnora@eveners.com",
          "joined_at": 1704954526
        }
      ],
      "joined_persons_count": 6,
      "is_joined": true,
      "status": "cancelled",
      "cancelled_reason": "Not enough sponsors to cover the cost",
      "cancelled_at": 1704954530
    },
    "ts": 1704954526
  }
  ```

- Response to non-organizer or non-participant

  ```json
  HTTP/1.1 200 OK
  Content-Type: application/json

  {
    "ok": true,
    "data": {
      "id": 1,
      "name": "Wedding Fulan",
      "venue": {
        "id": 1,
        "name": "Si Jalak Harupat"
      },
      "event": {
        "id": 1,
        "name": "Wedding"
      },
      "start_ts": 1704938400,
      "end_ts": 1704945600,
      "max_persons": 12,
      "organizer": {
        "id": 1,
        "username": "marion",
        "email": "marion@eveners.com"
      },
      "joined_persons_count": 6,
      "is_joined": false,
      "status": "open"
    },
    "ts": 1704954526
  }
  ```

**Error Response:**

No specific error response.

[Back to Top](#rest-api)

---

## Update Meetup

PUT: `/meetups/{meetup_id}`

This endpoint is used to update a meetup. Only the organizer of the meetup can update the meetup. Update action is limited to:

- Change the name of the meetup
- Change the start and end time of the meetup, but it still need to follows the rules of [Create Meetup](#create-meetup) endpoint
- Update maximum number of persons that can join the meetup

This endpoint can be used to force close a meetup by setting the `max_persons` to the number of persons that already joined the meetup.

**Header:**

- `Authorization` => The value is `Bearer {access_token}`.

**Example Request:**

```json
PUT /meetups/1
Content-Type: application/json
Authorization: Bearer {access_token}

{
  "name": "Wedding Fulan",
  "start_ts": 1704938400,
  "end_ts": 1704945600,
  "max_persons": 12
}
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json
Authorization: Bearer {access_token}

{
  "ok": true,
  "data": {
    "id": 1,
    "name": "Wedding Fulan",
    "venue": {
      "id": 1,
      "name": "Si Jalak Harupat"
    },
    "event": {
      "id": 1,
      "name": "Wedding"
    },
    "start_ts": 1704938400,
    "end_ts": 1704945600,
    "max_persons": 12,
    "organizer": {
      "id": 1,
      "username": "marion",
      "email": "marion@eveners.com"
    },
    "joined_persons_count": 6,
    "status": "open"
  },
  "ts": 1704954526
}
```

**Error Response:**

- User is not the organizer of the meetup

  ```json
  HTTP/1.1 403 Forbidden
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_FORBIDDEN",
    "msg": "User is not authorized to access this resource",
    "ts": 1704954526
  }
  ```

- Max persons is less than number of joined persons

  ```json
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_MAX_PERSONS_LESS_THAN_JOINED_PERSONS",
    "msg": "Max persons is less than number of joined persons",
    "ts": 1704954526
  }
  ```

- Exceed venue capacity

  ```json
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_EXCEED_VENUE_CAPACITY",
    "msg": "Venue capacity is full on the designated meetup time",
    "ts": 1704954526
  }
  ```

- Not within venue operating hours

  ```json
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_VENUE_IS_CLOSED",
    "msg": "Venue is closed on the designated meetup time",
    "ts": 1704954526
  }
  ```

[Back to Top](#rest-api)

---

## Cancel Meetup

DELETE: `/meetups/{meetup_id}`

This endpoint is used to cancel a meetup. Only the organizer of the meetup can cancel the meetup.
Meetup can only be cancelled if it isn't started yet.

Everytime a meetup is cancelled, the system needs to send notification to the email of all joined persons of the meetup. This notification will contains information about the reason of the meetup cancellation.

**Headers:**

- `Authorization` => The value is `Bearer {access_token}`.

**Query Params:**

- `cancelled_reason` => Reason why the meetup is cancelled.

**Example Request:**

```bash
DELETE /meetups/1?cancelled_reason=Not%20enough%20sponsors%20to%20cover%20the%20cost
Authorization: Bearer {access_token}
```

**Success Response:**

  ```json
  HTTP/1.1 200 OK
  Content-Type: application/json

  {
    "ok": true,
    "data": {
      "id": 1,
      "name": "Wedding Fulan",
      "venue": {
        "id": 1,
        "name": "Si Jalak Harupat"
      },
      "event": {
        "id": 1,
        "name": "Wedding"
      },
      "start_ts": 1704938400,
      "end_ts": 1704945600,
      "max_persons": 12,
      "organizer": {
        "id": 1,
        "username": "marion",
        "email": "marion@eveners.com"
      },
      "status": "cancelled",
      "cancelled_reason": "Not enough sponsors to cover the cost",
      "cancelled_at": 1704954530
    },
    "ts": 1704954526
  }
  ```

**Error Response:**

- Meetup is already started

  ```json
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_MEETUP_STARTED",
    "msg": "Meetup is started",
    "ts": 1704954526
  }
  ```

- Cancelled reason is required

  ```json
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_CANCELLED_REASON_REQUIRED",
    "msg": "Cancelled reason is required",
    "ts": 1704954526
  }
  ```

[Back to Top](#rest-api)

---

## Join Meetup

PUT: `/incoming-meetups/{meetup_id}`

This endpoint is used to join a meetup. User can only join a meetup if the meetup is still `open` which means the meetup hasn't reached the maximum number of persons, not cancelled, and not finished yet.

Also user cannot join different meetup that overlaps with start and end time of other meetups that he/she already joined.

Everytime a user joins a meetup, the system needs to send notification to the meetup organizer. This notification is simply contains the information about the user that just joined the meetup and the latest total number of joined persons in the meetup.

**Headers:**

- `Authorization` => The value is `Bearer {access_token}`.

**Example Request:**

```bash
PUT /incoming-meetups/1
Authorization: Bearer {access_token}
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "ok": true,
  "data": {
    "id": 1,
    "name": "Wedding Fulan",
    "venue": {
      "id": 1,
      "name": "Si Jalak Harupat"
    },
    "event": {
      "id": 1,
      "name": "Wedding"
    },
    "start_ts": 1704938400,
    "end_ts": 1704945600,
    "max_persons": 12,
    "organizer": {
      "id": 1,
      "username": "marion",
      "email": "marion@eveners.com"
    },
    "joined_persons": [
      {
        "id": 1,
        "username": "marion",
        "email": "marion@eveners.com",
        "joined_at": 1704954526
      },
      {
        "id": 2,
        "username": "todd",
        "email": "todd@eveners.com",
        "joined_at": 1704954526
      },
      {
        "id": 3,
        "username": "anthony",
        "email": "anthony@eveners.com",
        "joined_at": 1704954526
      },
      {
        "id": 4,
        "username": "kevin",
        "email": "kevin@eveners.com",
        "joined_at": 1704954526
      },
      {
        "id": 5,
        "username": "eric",
        "email": "eric@eveners.com",
        "joined_at": 1704954526
      },
      {
        "id": 6,
        "username": "elnora",
        "email": "elnora@eveners.com",
        "joined_at": 1704954526
      }
    ],
    "joined_persons_count": 6,
    "is_joined": true,
    "status": "open"
  },
  "ts": 1704954526
}
```

**Error Response:**

- Meetup is finished

  ```json
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_MEETUP_FINISHED",
    "msg": "Meetup is finished",
    "ts": 1704954526
  }
  ```

- Meetup is cancelled

  ```json
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_MEETUP_CANCELLED",
    "msg": "Meetup is cancelled",
    "ts": 1704954526
  }
  ```

- Meetup is closed

  ```json
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_MEETUP_CLOSED",
    "msg": "Meetup is closed",
    "ts": 1704954526
  }
  ```

- Meetup overlaps with other meetup that user already joined

  ```json
  HTTP/1.1 409 Conflict
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_MEETUP_OVERLAPS",
    "msg": "Meetup overlaps with other meetup that user already joined",
    "ts": 1704954526
  }
  ```

[Back to Top](#rest-api)

---

## Leave Meetup

DELETE: `/incoming-meetups/{meetup_id}`

This endpoint is used to leave a meetup. User can only leave a meetup if he/she already joined the meetup, also the meetup is not cancelled or finished yet.

After leaving a meetup, the user will be removed from the meetup' list of joined persons, and will opening a slot for other users to join the meetup.

Everytime a user leaves a meetup, the system needs to send notification to the meetup organizer. This notification is simply contains the information about the user that just left the meetup and the latest number of joined persons in the meetup.

**Headers:**

- `Authorization` => The value is `Bearer {access_token}`.

**Example Request:**

```json
DELETE /incoming-meetups/1
Authorization: Bearer {access_token}
```

**Success Response:**

  ```json
  HTTP/1.1 200 OK
  Content-Type: application/json

  {
    "ok": true,
    "ts": 1704954526
  }
  ```

**Error Response:**

- User is not a meetup participant

  ```json
  HTTP/1.1 403 Forbidden
  Content-Type: application/json

  {
    "ok": false,
    "err": "ERR_USER_NOT_PARTICIPANT",
    "msg": "User is not a participant",
    "ts": 1704954526
  }
  ```

[Back to Top](#rest-api)

---

## List Incoming Meetups

GET: `/incoming-meetups`

This endpoint is used to list future meetups that are joined by a user. The returned meetup statuses are either `open` or `cancelled`.

The reason why we still return <u>future</u> cancelled meetups is because we want to show the user that the future meetup is cancelled.

The list of meetups is sorted by the start time of the meetup, from the nearest to the furthest.

**Headers:**

- `Authorization` => The value is `Bearer {access_token}`.

**Query Params:**

- `status`, _OPTIONAL_ => Filter meetups that have the specified status. The value is one of the following: `open`, `cancelled`, `all`. By default it will set to `all` which return both `open` and `cancelled` meetups.
- `event_ids`, _OPTIONAL_ => Filter meetups that support the specified event ids. The value is comma separated list of event ids.
- `venue_ids`, _OPTIONAL_ => Filter meetups that are held in the specified venue ids. The value is comma separated list of venue ids.

**Example Request:**

```bash
GET /incoming-meetups/joined?event_ids=1,3&venue_ids=1,2&status=all
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "ok": true,
  "data": [
    {
      "id": 1,
      "name": "Wedding Fulan",
      "venue": {
        "id": 1,
        "name": "Si Jalak Harupat"
      },
      "event": {
        "id": 1,
        "name": "Wedding"
      },
      "start_ts": 1704938400,
      "end_ts": 1704945600,
      "max_persons": 12,
      "organizer": {
        "id": 1,
        "username": "marion",
        "email": "marion@eveners.com"
      },
      "joined_persons": [
        {
          "id": 1,
          "username": "marion",
          "email": "marion@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 2,
          "username": "todd",
          "email": "todd@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 3,
          "username": "anthony",
          "email": "anthony@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 4,
          "username": "kevin",
          "email": "kevin@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 5,
          "username": "eric",
          "email": "eric@eveners.com",
          "joined_at": 1704954526
        },
        {
          "id": 6,
          "username": "elnora",
          "email": "elnora@eveners.com",
          "joined_at": 1704954526
        }
      ],
      "joined_persons_count": 6,
      "is_joined": true,
      "status": "cancelled",
    },
  ],
  "ts": 1704954526
}
```

**Error Response:**

No specific error response.

[Back to Top](#rest-api)

---