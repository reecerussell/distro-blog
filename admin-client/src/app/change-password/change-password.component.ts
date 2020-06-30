import { Component, OnInit } from "@angular/core";
import { ApiService } from "../api.service";

@Component({
    selector: "app-change-password",
    templateUrl: "./change-password.component.html",
    styleUrls: ["./change-password.component.scss"],
})
export class ChangePasswordComponent implements OnInit {
    error: string = null;
    success: string = null;
    loading: boolean = false;
    isDirty: boolean = false;

    currentPassword: string = "";
    newPassword: string = "";
    confirmPassword: string = "";

    currentPasswordError: string = null;
    newPasswordError: string = null;
    confirmPasswordError: string = null;

    constructor(private api: ApiService) {}

    ngOnInit(): void {}

    validateCurrentPassword(): any {
        this.isDirty = true;

        const password = this.currentPassword;
        if (password.length < 1) {
            this.currentPasswordError = "Please enter your current password.";
            return false;
        } else if (password.length > 256) {
            this.currentPasswordError =
                "Password cannot be longer than 256 characters.";
            return false;
        } else {
            this.currentPasswordError = null;
        }

        return true;
    }

    validateNewPassword(): any {
        this.isDirty = true;

        const password = this.newPassword;
        if (password.length < 1) {
            this.newPasswordError = "Please enter a new password.";
            return false;
        } else if (password.length > 256) {
            this.newPasswordError =
                "Your new password cannot be longer than 256 characters.";
            return false;
        } else {
            this.newPasswordError = null;
        }

        return true;
    }

    validateConfirmPassword(): any {
        this.isDirty = true;

        const password = this.confirmPassword;
        if (password.length < 1) {
            this.confirmPasswordError = "Please confirm your new password.";
            return false;
        } else if (password.length > 256) {
            this.confirmPasswordError =
                "Your new password cannot be longer than 256 characters.";
            return false;
        } else if (password !== this.newPassword) {
            this.confirmPasswordError = "Your passwords do not match!";
        } else {
            this.confirmPasswordError = null;
        }

        return true;
    }

    validate() {
        let valid = true;

        if (!this.validateCurrentPassword()) {
            valid = false;
        }

        if (!this.validateNewPassword()) {
            valid = false;
        }

        if (!this.validateConfirmPassword()) {
            valid = false;
        }

        return valid;
    }

    async onSubmit() {
        if (this.loading || !this.validate()) {
            return;
        }

        this.loading = true;

        const res = await this.api.Users.ChangePassword(
            this.currentPassword,
            this.newPassword
        );

        if (res.ok) {
            this.error = null;
            this.success = "You password has been changed successfully!";

            setTimeout(() => (this.success = null), 4000);
        } else {
            this.error = res.error;
        }

        this.loading = false;
    }
}
