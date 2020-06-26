import { Component, OnInit, ViewEncapsulation } from "@angular/core";
import { ApiService } from "../api.service";
import { UserService } from "../user.service";

@Component({
    selector: "app-login",
    templateUrl: "./login.component.html",
    styleUrls: ["./login.component.scss"],
    encapsulation: ViewEncapsulation.None,
    host: {
        "[class.login-background]": "true",
    },
})
export class LoginComponent implements OnInit {
    constructor(private api: ApiService, private user: UserService) {}

    loading: boolean = false;
    error: string = null;
    isDirty: boolean = false;

    email: string = "";
    emailError: string = null;
    password: string = "";
    passwordError: string = null;

    ngOnInit(): void {}

    validateEmail(isSubmit: boolean): any {
        this.isDirty = true;

        const email = this.email;
        if (email.length < 1) {
            this.emailError = "This field is required.";
            return false;
        } else if (email.length > 100) {
            this.emailError = "Email cannot be longer than 100 characters";
            return false;
        } else if (
            !email.match("[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,6}") &&
            isSubmit
        ) {
            this.emailError = "Please enter a valid email address.";
            return false;
        } else {
            this.emailError = null;
        }

        return true;
    }

    validatePassword(): any {
        this.isDirty = true;

        if (this.password.length < 1) {
            this.passwordError = "Please enter your password.";
            return false;
        } else {
            this.passwordError = null;
        }

        return true;
    }

    validate() {
        let valid = true;

        if (!this.validateEmail(true)) {
            valid = false;
        }

        if (!this.validatePassword()) {
            valid = false;
        }

        return valid;
    }

    async onSubmit() {
        if (this.loading || !this.validate()) {
            return;
        }

        this.loading = true;

        const res = await this.api.Login(this.email, this.password);
        if (res.ok) {
            const { token, expires } = res.data;
            this.user.Login(token, expires);
        } else {
            this.error = res.error;
        }

        this.loading = false;
    }
}
