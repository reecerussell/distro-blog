interface CreateUser {
    firstname: string;
    lastname: string;
    email: string;
    password: string;
}

interface UserListItem {
    id: string;
    name: string;
    email: string;
}

interface UpdateUser {
    id: string;
    firstname: string;
    lastname: string;
    email: string;
}

export { CreateUser, UserListItem, UpdateUser };
