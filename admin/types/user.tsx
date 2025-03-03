import { RoleType } from "./role";

export interface UserType {
    active: boolean;
    created_at: string;
    email: string;
    id: string;
    last_login: string;
    role: RoleType;
    role_id: string;
    updated_at: string;
    username: string;
}

export interface CreateUserType {
    username: string;
    email: string;
    password: string;
    role_id: string;
    active: boolean;
}

export interface UpdateUserType {
    username?: string;
    email?: string;
    password?: string;
    active?: boolean;
    role_id?: string;
}

export interface AuthUserType {
    email: string;
    id: string;
    role: "super_admin" | "content_admin" | "editor" | "viewer"
    username: string;
}