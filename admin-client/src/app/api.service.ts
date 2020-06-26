import { Injectable } from "@angular/core";
import { CreateUser, UpdateUser } from "./models/user";

interface ApiResponse {
    ok: boolean;
    error: string;
    data: any;
}

class UserService {
    constructor(private baseUrl: string) {
        this.baseUrl = baseUrl + "users";
    }

    async List(): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl, {
            method: "GET",
            headers: {
                Authorization:
                    "Bearer " +
                    localStorage.getItem("distro_blog:access_token"),
            },
        });
        return await this.parseResponse(res);
    }

    async Get(id: string): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl + `/${id}`, {
            method: "GET",
            headers: {
                Authorization:
                    "Bearer " +
                    localStorage.getItem("distro_blog:access_token"),
            },
        });
        return await this.parseResponse(res);
    }

    async Create(data: CreateUser): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization:
                    "Bearer " +
                    localStorage.getItem("distro_blog:access_token"),
            },
            body: JSON.stringify(data),
        });

        return await this.parseResponse(res);
    }

    async Update(data: UpdateUser): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                Authorization:
                    "Bearer " +
                    localStorage.getItem("distro_blog:access_token"),
            },
            body: JSON.stringify(data),
        });

        return await this.parseResponse(res);
    }

    async Delete(id: string): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl + `/${id}`, {
            method: "DELETE",
            headers: {
                Authorization:
                    "Bearer " +
                    localStorage.getItem("distro_blog:access_token"),
            },
        });
        return await this.parseResponse(res);
    }

    async parseResponse(res: Response): Promise<ApiResponse> {
        if (res.status === 200) {
            const { data } = await res.json();
            return {
                ok: true,
                error: null,
                data: data,
            };
        }

        const { error } = await res.json();
        return {
            ok: false,
            error: error,
            data: null,
        };
    }
}

@Injectable({
    providedIn: "root",
})
export class ApiService {
    baseUrl = "https://kw6lirghub.execute-api.eu-west-2.amazonaws.com/dev/";

    constructor() {
        this.Users = new UserService(this.baseUrl);
    }

    public Users: UserService;

    async Login(email: string, password: string): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl + "token", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
        });

        if (res.status === 200) {
            const { data } = await res.json();
            return {
                ok: true,
                error: null,
                data: data,
            };
        }

        const { error } = await res.json();
        return {
            ok: false,
            error: error,
            data: null,
        };
    }
}
