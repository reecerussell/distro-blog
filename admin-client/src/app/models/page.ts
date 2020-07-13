interface CreatePage {
    title: string;
    description: string;
    content: string;
}

interface UpdatePage {
    id: string;
    title: string;
    description: string;
    content: string;
}

interface Page {
    id: string;
    title: string;
    description: string;
    content: string;
    isBlog: boolean;
    isActive: boolean;
    imageId?: string;
    audit?: PageAudit[];
}

interface PageAudit {
    UserFullname: string;
    UserId: string;
    Date: Date;
    Message: string;
}

export { CreatePage, UpdatePage, Page, PageAudit };
