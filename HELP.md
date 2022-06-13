## Supported endpoints

- "/" : NO-OP  
- "/v1/users": usersV1(w, r)
- "/v1/jwt": jwtHandler(w, r)
- "/v1/auth": authHandler(w, r)
- "/v1/accounts": usersV1Handler(w, r)
- "/v2/accounts": usersV2V3Handler(w, r)
- "/v3/accounts": usersV2V3Handler(w, r)
- "/api/entitlements/v1/services": entitlements(w, r)


## usersV1:

 - Only admits POST requests
 - Reads the POST body request, sets filter object with it. 
 - Reads QS parameters
    - "admin_only"
    - "status"
    - "limit"
 - Calls "findUsersBy" with accountNo and OrgId as "", input as nil. , passes in admin_only, status, limit and filter object and retrieves list of users
 - prints list of users.

 ## findUsersBy(accountNo, orgId, adminOnly, status, limit, input, users)

- Calls "getUsers()"

## getUsers()

- Requests all users from Keycloak admin endpoint
- Iterates over the users list and:
 + filters out based on "adminOnly" 

