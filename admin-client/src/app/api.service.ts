import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { fromFetch } from "rxjs/fetch";
import { Observable } from "rxjs";
import Validate from "./models/validation/user";
import { CreateUser } from "./models/user";

@Injectable({
    providedIn: "root",
})
export class ApiService {
    baseUrl = "https://kw6lirghub.execute-api.eu-west-2.amazonaws.com/dev/";

    constructor(private http: HttpClient) {}

    CreateUser(data: CreateUser): Promise<Response> {
        const error = Validate(data);
        if (error) {
            throw new Error(error);
        }

        return fetch(this.baseUrl + "users", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        });
    }

    Users() {
        return this.http.get(this.baseUrl + "users");
    }
}
