import { Injectable } from "@angular/core";
import { CreateUser, UpdateUser } from "./models/user";

@Injectable({
    providedIn: "root",
})
export class ApiService {
    baseUrl = "https://kw6lirghub.execute-api.eu-west-2.amazonaws.com/dev/";

    constructor() {}

    CreateUser(data: CreateUser): Promise<Response> {
        return fetch(this.baseUrl + "users", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        });
    }

    Users(): Promise<Response> {
        return fetch(this.baseUrl + "users");
    }

    GetUser(id: string): Promise<Response> {
        return fetch(this.baseUrl + "users/" + id);
    }

    UpdateUser(data: UpdateUser): Promise<Response> {
        return fetch(this.baseUrl + "users", {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        });
    }
}
