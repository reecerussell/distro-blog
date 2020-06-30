import { Injectable } from "@angular/core";

@Injectable({
    providedIn: "root",
})
export class UserService {
    AccessTokenKey = "distro_blog:access_token";
    ExpiresKey = "distro_blog:token_expiry";

    listeners: Map<string, any>;

    constructor() {
        this.listeners = new Map<string, any>();
    }

    Listen(key: string, callback: any): void {
        this.listeners.set(key, callback);
    }

    Unlisten(key: string): void {
        this.listeners.delete(key);
    }

    Login(token, expires): void {
        expires *= 1000;
        const utcDate = new Date(expires).toUTCString();
        const localDate = new Date(utcDate).getTime();

        localStorage.setItem(this.AccessTokenKey, token);
        localStorage.setItem(this.ExpiresKey, localDate.toString());

        trigger(this.listeners);
    }

    Logout(): void {
        localStorage.removeItem(this.AccessTokenKey);
        localStorage.removeItem(this.ExpiresKey);

        trigger(this.listeners);
    }

    IsAuthenticated(): boolean {
        const token = localStorage.getItem(this.AccessTokenKey);
        if (!token) {
            return false;
        }

        const time = localStorage.getItem(this.ExpiresKey);
        const expiry = new Date(parseInt(time));

        return !(expiry < new Date());
    }

    GetId(): string {
        return getPayload(localStorage.getItem(this.AccessTokenKey))["uid"];
    }

    IsInScope(scopeName: string): boolean {
        const token = localStorage.getItem(this.AccessTokenKey);
        const scopes = getPayload(token)["scps"];

        for (let i = 0; i < scopes.length; i++) {
            if (scopes[i] === scopeName) {
                return true;
            }
        }

        return false;
    }
}

const trigger = (actions: Map<string, any>) => {
    actions.forEach((callback) => callback());
};

const getPayload = (token: string) => {
    if (!token) {
        return null;
    }

    const parts = token.split(".");
    if (parts.length < 2) {
        return null;
    }

    const payloadData = atob(parts[1]);

    return JSON.parse(payloadData);
};
