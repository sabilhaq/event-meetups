# Common Errors

This document specify common errors that may occurs throughout all specified REST API endpoints in the system.

- Bad Request (`400`)

    ```json
    HTTP/1.1 400 Bad Request
    Content-Type: application/json

    {
        "ok": false,
        "err": "ERR_BAD_REQUEST",
        "msg": "invalid value of `type`",
        "ts": 1676989698
    }
    ```

    Client will receive this error when there is an issue with client request. Check `msg` for the cause details.

- Invalid Access Token (`401`)

    ```json
    HTTP/1.1 401 Unauthorized
    Content-Type: application/json

    {
        "ok": false,
        "err": "ERR_INVALID_ACCESS_TOKEN",
        "msg": "invalid access token",
        "ts": 1676989698
    }
    ```

    Client will receive this error when the submitted access token is no longer valid (expired) or the token is literally invalid (wrong access token).

- Forbidden Access (`403`)

    ```json
    HTTP/1.1 403 Forbidden
    Content-Type: application/json

    {
        "ok": false,
        "err": "ERR_FORBIDDEN_ACCESS",
        "msg": "user doesn't have enough authorization",
        "ts": 1676989698
    }
    ```

    Client will receive this error when it tried to access resource that unauthorized for user.

- Not Found (`404`)

    ```json
    HTTP/1.1 404 Not Found
    Content-Type: application/json

    {
        "ok": false,
        "err": "ERR_NOT_FOUND",
        "msg": "resource is not found",
        "ts": 1676989698
    }
    ```

    Client will receive this error when the resource it tried to access is not found.

- Internal Server Error (`500`)

    ```json
    {
        "status": 500,
        "err": "ERR_INTERNAL_ERROR",
        "msg": "unable to connect into database",
        "ts": 1676989698
    }
    ```

    Client will receive this error when there is some issue in the server. Check `msg` for details.

[Back to Top](#common-errors)

---