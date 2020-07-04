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

export { CreatePage, UpdatePage };
