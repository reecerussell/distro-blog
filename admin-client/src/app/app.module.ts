import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";
import { HttpClientModule } from "@angular/common/http";
import { FormsModule } from "@angular/forms";
import { CKEditorModule } from "@ckeditor/ckeditor5-angular";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { DashboardComponent } from "./dashboard/dashboard.component";
import { CreateUserComponent } from "./create-user/create-user.component";
import { UserListComponent } from "./user-list/user-list.component";
import { EditUserComponent } from "./edit-user/edit-user.component";
import { LoginComponent } from "./login/login.component";
import { UserInfoComponent } from "./user-info/user-info.component";
import { ScopedComponent } from "./scoped/scoped.component";
import { ChangePasswordComponent } from "./change-password/change-password.component";
import { ResetPasswordComponent } from "./reset-password/reset-password.component";
import { PageListComponent } from "./page-list/page-list.component";
import { CreatePageComponent } from "./create-page/create-page.component";
import { EditPageComponent } from './edit-page/edit-page.component';
import { PageInfoComponent } from './page-info/page-info.component';
import { BlogListComponent } from './blog-list/blog-list.component';
import { CreateBlogComponent } from './create-blog/create-blog.component';
import { SettingsListComponent } from './settings-list/settings-list.component';

@NgModule({
    declarations: [
        AppComponent,
        DashboardComponent,
        CreateUserComponent,
        UserListComponent,
        EditUserComponent,
        LoginComponent,
        UserInfoComponent,
        ScopedComponent,
        ChangePasswordComponent,
        ResetPasswordComponent,
        PageListComponent,
        CreatePageComponent,
        EditPageComponent,
        PageInfoComponent,
        BlogListComponent,
        CreateBlogComponent,
        SettingsListComponent,
    ],
    imports: [
        BrowserModule,
        AppRoutingModule,
        HttpClientModule,
        FormsModule,
        CKEditorModule,
    ],
    providers: [],
    bootstrap: [AppComponent],
})
export class AppModule {}
