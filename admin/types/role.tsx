import { PermissionsType } from "./permissions";

export interface RoleType {
    ID: string;
    Name: string;
    Type: "super_admin" | "content_admin" | "editor" | "viewer"
    Description: string;
    Permissions: PermissionsType[];
};