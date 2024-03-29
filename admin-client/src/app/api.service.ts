import { Injectable } from "@angular/core";
import { CreateUser, UpdateUser } from "./models/user";
import { CreatePage, UpdatePage } from "./models/page";
import Setting from "./models/setting";
import Constants from "./common/constants";

interface ApiResponse {
    ok: boolean;
    error: string;
    data: any;
}

const getAuthHeader = () =>
    "Bearer " + localStorage.getItem("distro_blog:access_token");

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
        return await parseResponse(res);
    }

    async Get(id: string, ...expand): Promise<ApiResponse> {
        let queryString = "?";
        for (let i = 0; i < expand.length; i++) {
            if (i !== 0) {
                queryString += "&";
            }

            queryString += "expand=" + expand[i];
        }

        const res = await fetch(this.baseUrl + `/${id}` + queryString, {
            method: "GET",
            headers: {
                Authorization:
                    "Bearer " +
                    localStorage.getItem("distro_blog:access_token"),
            },
        });
        return await parseResponse(res);
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

        return await parseResponse(res);
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

        return await parseResponse(res);
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
        return await parseResponse(res);
    }

    async ChangePassword(
        currentPassword: string,
        newPassword: string
    ): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl + "/password", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization:
                    "Bearer " +
                    localStorage.getItem("distro_blog:access_token"),
            },
            body: JSON.stringify({
                currentPassword,
                newPassword,
            }),
        });

        return await parseResponse(res);
    }

    async ResetPassword(id: string): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl + `/password/reset/${id}`, {
            method: "POST",
            headers: {
                Authorization:
                    "Bearer " +
                    localStorage.getItem("distro_blog:access_token"),
            },
        });

        return await parseResponse(res);
    }
}

class PageService {
    constructor(private baseUrl: string) {
        this.baseUrl = baseUrl + "pages";
    }

    async List(): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl, {
            method: "GET",
            headers: {
                Authorization: getAuthHeader(),
            },
        });
        return await parseResponse(res);
    }

    async Create(data: CreatePage): Promise<ApiResponse> {
        try {
            const res = await fetch(this.baseUrl, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: getAuthHeader(),
                },
                body: JSON.stringify(data),
            });

            return await parseResponse(res);
        } catch (err) {
            return new Promise(() => ({
                ok: false,
                data: null,
                error: err,
            }));
        }
    }

    async Get(id: string, ...expand): Promise<ApiResponse> {
        let queryString = "?";
        for (let i = 0; i < expand.length; i++) {
            if (i !== 0) {
                queryString += "&";
            }

            queryString += "expand=" + expand[i];
        }

        const res = await fetch(this.baseUrl + `/${id}` + queryString, {
            method: "GET",
            headers: {
                Authorization: getAuthHeader(),
            },
        });
        return await parseResponse(res);
    }

    async Update(data: UpdatePage, image: Blob = null): Promise<ApiResponse> {
        let res: Response;

        if (image === null) {
            res = await fetch(this.baseUrl, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: getAuthHeader(),
                },
                body: JSON.stringify(data),
            });
        } else {
            const formData = new FormData();
            Object.keys(data).map((k) => {
                if (k === "seo") {
                    const { title, description, index, follow } = data[k];
                    formData.set("seoTitle", title);
                    formData.set("seoDescription", description);
                    formData.set("seoIndex", index.toString());
                    formData.set("seoFollow", follow.toString());

                    return;
                }

                formData.set(k, data[k]);
            });
            formData.set("image", image, "uploaded-image");

            res = await fetch(this.baseUrl, {
                method: "PUT",
                headers: {
                    Authorization: getAuthHeader(),
                },
                body: formData,
            });
        }

        return await parseResponse(res);
    }

    async Activate(id: string): Promise<ApiResponse> {
        const res = await fetch(Constants.APIBase + `pages/${id}/activate`, {
            method: "POST",
            headers: {
                Authorization: getAuthHeader(),
            },
        });
        return await parseResponse(res);
    }

    async Deactivate(id: string): Promise<ApiResponse> {
        const res = await fetch(Constants.APIBase + `pages/${id}/deactivate`, {
            method: "POST",
            headers: {
                Authorization: getAuthHeader(),
            },
        });
        return await parseResponse(res);
    }

    async DropdownOptions(): Promise<ApiResponse> {
        const res = await fetch(Constants.APIBase + `pages/dropdown`, {
            method: "GET",
            headers: {
                Authorization: getAuthHeader(),
            },
        });
        return await parseResponse(res);
    }
}

class BlogService {
    constructor(private baseUrl: string) {
        this.baseUrl = baseUrl + "blogs";
    }

    async List(): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl, {
            method: "GET",
            headers: {
                Authorization: getAuthHeader(),
            },
        });
        return await parseResponse(res);
    }

    async Create(data: CreatePage): Promise<ApiResponse> {
        try {
            const res = await fetch(this.baseUrl, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: getAuthHeader(),
                },
                body: JSON.stringify(data),
            });

            return await parseResponse(res);
        } catch (err) {
            return new Promise(() => ({
                ok: false,
                data: null,
                error: err,
            }));
        }
    }
}

class SettingsService {
    constructor(private baseUrl: string) {
        this.baseUrl = baseUrl + "settings";
    }

    async List(): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl, {
            method: "GET",
            headers: {
                Authorization: getAuthHeader(),
            },
        });
        return await parseResponse(res);
    }

    async Get(key: string): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl + `/${key}`, {
            method: "GET",
            headers: {
                Authorization: getAuthHeader(),
            },
        });
        return await parseResponse(res);
    }

    async Update(data: Setting): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl, {
            method: "PUT",
            headers: {
                Authorization: getAuthHeader(),
            },
            body: JSON.stringify(data),
        });
        return await parseResponse(res);
    }
}

const parseResponse = async (res: Response): Promise<ApiResponse> => {
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
};

@Injectable({
    providedIn: "root",
})
export class ApiService {
    baseUrl: string;

    constructor() {
        this.baseUrl = Constants.APIBase;

        this.Users = new UserService(this.baseUrl);
        this.Pages = new PageService(this.baseUrl);
        this.Blogs = new BlogService(this.baseUrl);
        this.Settings = new SettingsService(this.baseUrl);
    }

    public Users: UserService;
    public Pages: PageService;
    public Blogs: BlogService;
    public Settings: SettingsService;

    async Login(email: string, password: string): Promise<ApiResponse> {
        const res = await fetch(this.baseUrl + "token", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ email, password }),
        });

        return parseResponse(res);
    }
}
