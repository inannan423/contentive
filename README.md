# _"Contentive"_ 📖 Headless CMS

## **Description**

Contentive is an open-source headless CMS built with Go and PostgreSQL. It allows you to create and manage content together with your flexible frontend.

- **Go**
- **PostgreSQL**

## **Installation**

To install Contentive, you will need to have Go and PostgreSQL installed on your system.

<!-- TODO -->

## API Documentation

### Content Types

#### Get Content Type

- **URL**: `/api/content-types/:identifier`
- **Method**: `GET`
- **Parameters**:
  - `identifier`: Content Type ID or slug

#### Create Content Type

- **URL**: `/api/content-types`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "blog_post",
    "type": "collection"
  }
  ```
- **Response**: `201 Created`
  ```json
  {
    "id": "uuid",
    "name": "blog_post",
    "type": "collection",
    "createdAt": "2024-03-21T10:00:00Z",
    "updatedAt": "2024-03-21T10:00:00Z",
    "fields": []
  }
  ```

#### Get All Content Types

- **URL**: `/api/content-types`
- **Method**: `GET`
- **Response**: `200 OK`
  ```json
  [
    {
      "id": "uuid",
      "name": "blog_post",
      "type": "collection",
      "createdAt": "2024-03-21T10:00:00Z",
      "updatedAt": "2024-03-21T10:00:00Z",
      "fields": [...]
    }
  ]
  ```

### Fields

#### Add Field to Content Type

- **URL**: `/api/content-types/:contentTypeId/fields`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "label": "title",
    "type": "text"
  }
  ```
- **Supported Field Types**:
  - `text`: Plain text
  - `rich_text`: Rich text editor
  - `number`: Numeric value
  - `date`: ISO 8601 date
  - `boolean`: True/false value
  - `enum`: Enumerated value
  - `relation`: Reference to another content type
- **Response**: `201 Created`
  ```json
  {
    "id": "uuid",
    "contentTypeId": "uuid",
    "label": "title",
    "type": "text",
    "createdAt": "2024-03-21T10:00:00Z",
    "updatedAt": "2024-03-21T10:00:00Z"
  }
  ```

#### Update Field

- **URL**: `/api/content-types/:contentTypeId/fields/:id`
- **Method**: `PUT`
- **Request Body**: Same as Add Field
- **Response**: `200 OK`

#### Delete Field

- **URL**: `/api/content-types/:contentTypeId/fields/:id`
- **Method**: `DELETE`
- **Response**: `204 No Content`

### Content Entries

#### Create Content Entry

- **URL**: `/api/content-types/:contentTypeId/entries`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "title": "My First Post",
    "content": "Hello World!",
    "published": true,
    "publishDate": "2024-03-21T10:00:00Z"
  }
  ```
  Note: The fields in the request body must match the fields defined in the content type.
- **Response**: `201 Created`
  ```json
  {
    "id": "uuid",
    "contentTypeId": "uuid",
    "data": {
      "title": "My First Post",
      "content": "Hello World!",
      "published": true,
      "publishDate": "2024-03-21T10:00:00Z"
    },
    "createdAt": "2024-03-21T10:00:00Z",
    "updatedAt": "2024-03-21T10:00:00Z"
  }
  ```

### Error Responses

All endpoints may return the following error responses:

- `400 Bad Request`: Invalid input data
- `404 Not Found`: Requested resource not found
- `500 Internal Server Error`: Server-side error

Error response body format:

```json
{
  "error": "Error message description"
}
```
