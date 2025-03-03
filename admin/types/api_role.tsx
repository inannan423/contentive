// Get Api Roles:

// [
//     {
//         "ID": "6222a27a-38f2-4212-ae2a-7dfd2bc152f5",
//         "Name": "Blog Reader",
//         "Type": "custom",
//         "Description": "只能读取博客文章的API角色",
//         "APIKey": "l8--qzM-oUfGi-03pfpAGLNQw-oXJ-nT_XbjwZ8R2g4=",
//         "IsSystem": false,
//         "CreatedAt": "2025-03-03T17:17:53.974209+08:00",
//         "UpdatedAt": "2025-03-03T17:17:53.974209+08:00",
//         "Permissions": []
//     },
//     {
//         "ID": "469dfdcc-3167-4464-a8a2-95d2c40a93d3",
//         "Name": "Public User",
//         "Type": "public_user",
//         "Description": "Public Access",
//         "APIKey": "",
//         "IsSystem": true,
//         "CreatedAt": "2025-03-03T15:58:29.99005+08:00",
//         "UpdatedAt": "2025-03-03T17:21:14.783848+08:00",
//         "Permissions": [
//             {
//                 "ID": "cc497044-e59d-4502-a4a7-79d25b4e858c",
//                 "APIRoleID": "469dfdcc-3167-4464-a8a2-95d2c40a93d3",
//                 "ContentTypeID": "e3a29cff-6bb2-430a-a4ae-647245cba475",
//                 "ContentType": {
//                     "ID": "e3a29cff-6bb2-430a-a4ae-647245cba475",
//                     "slug": "blog-post",
//                     "name": "blog_post2212",
//                     "type": "collection",
//                     "createdAt": "2025-03-03T17:19:12.756523+08:00",
//                     "updatedAt": "2025-03-03T17:19:12.756523+08:00",
//                     "fields": null
//                 },
//                 "Operation": "read",
//                 "Enabled": true
//             }
//         ]
//     },
//     {
//         "ID": "5c32e0d1-edbd-4f0e-b99b-829f10d6e4f4",
//         "Name": "Authenticated User",
//         "Type": "authenticated_user",
//         "Description": "Authenticated Access",
//         "APIKey": "_APbQrO_HeP7j3ZO3U-zAO7X_goCHiC9jb5pxpzFgqs=",
//         "IsSystem": true,
//         "CreatedAt": "2025-03-03T15:58:29.991013+08:00",
//         "UpdatedAt": "2025-03-03T17:21:14.784314+08:00",
//         "Permissions": [
//             {
//                 "ID": "929ccc6a-9b8f-4984-ae7d-7293ff2feffa",
//                 "APIRoleID": "5c32e0d1-edbd-4f0e-b99b-829f10d6e4f4",
//                 "ContentTypeID": "e3a29cff-6bb2-430a-a4ae-647245cba475",
//                 "ContentType": {
//                     "ID": "e3a29cff-6bb2-430a-a4ae-647245cba475",
//                     "slug": "blog-post",
//                     "name": "blog_post2212",
//                     "type": "collection",
//                     "createdAt": "2025-03-03T17:19:12.756523+08:00",
//                     "updatedAt": "2025-03-03T17:19:12.756523+08:00",
//                     "fields": null
//                 },
//                 "Operation": "create",
//                 "Enabled": true
//             },
//             {
//                 "ID": "f351a584-405b-4f3f-bb88-b6c746c35440",
//                 "APIRoleID": "5c32e0d1-edbd-4f0e-b99b-829f10d6e4f4",
//                 "ContentTypeID": "e3a29cff-6bb2-430a-a4ae-647245cba475",
//                 "ContentType": {
//                     "ID": "e3a29cff-6bb2-430a-a4ae-647245cba475",
//                     "slug": "blog-post",
//                     "name": "blog_post2212",
//                     "type": "collection",
//                     "createdAt": "2025-03-03T17:19:12.756523+08:00",
//                     "updatedAt": "2025-03-03T17:19:12.756523+08:00",
//                     "fields": null
//                 },
//                 "Operation": "read",
//                 "Enabled": true
//             },
//             {
//                 "ID": "ec34c296-b01a-49ca-93c0-12a432eb3966",
//                 "APIRoleID": "5c32e0d1-edbd-4f0e-b99b-829f10d6e4f4",
//                 "ContentTypeID": "e3a29cff-6bb2-430a-a4ae-647245cba475",
//                 "ContentType": {
//                     "ID": "e3a29cff-6bb2-430a-a4ae-647245cba475",
//                     "slug": "blog-post",
//                     "name": "blog_post2212",
//                     "type": "collection",
//                     "createdAt": "2025-03-03T17:19:12.756523+08:00",
//                     "updatedAt": "2025-03-03T17:19:12.756523+08:00",
//                     "fields": null
//                 },
//                 "Operation": "update",
//                 "Enabled": true
//             },
//             {
//                 "ID": "c1cf98e9-2dcb-4889-8f19-c6ce3be255c2",
//                 "APIRoleID": "5c32e0d1-edbd-4f0e-b99b-829f10d6e4f4",
//                 "ContentTypeID": "e3a29cff-6bb2-430a-a4ae-647245cba475",
//                 "ContentType": {
//                     "ID": "e3a29cff-6bb2-430a-a4ae-647245cba475",
//                     "slug": "blog-post",
//                     "name": "blog_post2212",
//                     "type": "collection",
//                     "createdAt": "2025-03-03T17:19:12.756523+08:00",
//                     "updatedAt": "2025-03-03T17:19:12.756523+08:00",
//                     "fields": null
//                 },
//                 "Operation": "delete",
//                 "Enabled": true
//             }
//         ]
//     }
// ]



export type APIRoleType = {
  ID: string;
  Name: string;
  Type: "custom" | "public_user" | "authenticated_user"
  Description: string;
  APIKey: string;
  IsSystem: boolean;
  ExpiresAt: string | null;
  CreatedAt: string;
  UpdatedAt: string;
  Permissions: APIPermissionType[] | null;
};

export type APIPermissionType = {
  ID: string;
  APIRoleID: string;
  ContentTypeID: string;
  ContentType: {
    ID: string;
    slug: string;
    name: string;
    type: string;
    createdAt: string;
    updatedAt: string;
    fields: null;
  };
};

// {
//   "name": "Blog Reader",
//   "description": "只能读取博客文章的API角色"
// }
export interface CreateAPIRoleType {
  name: string;
  description: string;
  expires_at?: string | null;
}

export interface UpdateAPIRoleType {
  name?: string;
  description?: string;
  expires_at?: string | null;
}