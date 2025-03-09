# _"Contentive"_ 📖 Headless CMS

## **Description**

Contentive is an open-source headless CMS built with Go and PostgreSQL. It allows you to create and manage content together with your flexible frontend.

- **Go**
- **PostgreSQL**

## **Installation**

To install Contentive, you will need to have Go and PostgreSQL installed on your system.

## **API Documentation**

### Authentication

#### Login

```typescript
POST / admin / auth / login;

// Request Body
interface LoginRequest {
  email: string;
  password: string;
}

// Response
interface LoginResponse {
  token: string;
}
```

#### Validate Token

```typescript
GET / admin / auth / validate;

// Headers
Authorization: Bearer<token>;

// Response
interface ValidateResponse {
  valid: boolean;
  user?: {
    id: string;
    email: string;
    name: string;
    role: Role;
  };
}
```

### Users

#### Get All Users

```typescript
GET / admin / users;

// Headers
Authorization: Bearer<token>;

// Required Permission: ManageUsers

// Response
interface User {
  id: string;
  email: string;
  name: string;
  role: Role;
  createdAt: string;
  updatedAt: string;
}

type GetUsersResponse = User[];
```

#### Create User

```typescript
POST / admin / users;

// Headers
Authorization: Bearer<token>;

// Required Permission: ManageUsers

// Request Body
interface CreateUserRequest {
  email: string;
  password: string;
  name: string;
  roleId: string;
}

// Response
interface CreateUserResponse extends User {}
```

#### Update User

```typescript
PUT /admin/users/:id

// Headers
Authorization: Bearer <token>

// Required Permission: ManageUsers
// Note: Super Admin required for updating other users

// Request Body
interface UpdateUserRequest {
  email?: string;
  password?: string;
  name?: string;
  roleId?: string;
}

// Response
interface UpdateUserResponse extends User {}
```

#### Delete User

```typescript
DELETE /admin/users/:id

// Headers
Authorization: Bearer <token>

// Required Permission: ManageUsers
// Note: Super Admin required for deleting other users

// Response
interface DeleteUserResponse {
  success: boolean;
  message: string;
}
```

### Roles

#### Get All Roles

```typescript
GET / admin / roles;

// Headers
Authorization: Bearer<token>;

// Required Permission: ManageRoles

// Response
interface Permission {
  id: string;
  name: string;
  description: string;
}

interface Role {
  id: string;
  name: string;
  permissions: Permission[];
  createdAt: string;
  updatedAt: string;
}

type GetRolesResponse = Role[];
```

#### Get Role

```typescript
GET /admin/roles/:id

// Headers
Authorization: Bearer <token>

// Required Permission: ManageRoles

// Response
interface GetRoleResponse extends Role {}
```

#### Get All Permissions

```typescript
GET / admin / permissions;

// Headers
Authorization: Bearer<token>;

// Required Permission: ManageRoles

// Response
type GetPermissionsResponse = Permission[];
```
