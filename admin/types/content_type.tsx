export interface ContentTypeType {
    ID: string;
    slug: string;
    name: string;
    type: string;
    createdAt: string;
    updatedAt: string;
    fields: ContentFieldType[];
}

export interface ContentFieldType {
    ID: string;
    ContentTypeID: string;
    Label: string;
    Type: string;
    Required: boolean;
    CreatedAt: string;
    UpdatedAt: string;
}

export interface ContentFieldRequestType {
    label: string;
    type: string;
    required: boolean;
}