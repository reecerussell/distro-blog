import { NgModule } from "@angular/core";
import { Routes, RouterModule } from "@angular/router";
import { DashboardComponent } from "./dashboard/dashboard.component";
import { CreateUserComponent } from "./create-user/create-user.component";
import { UserListComponent } from "./user-list/user-list.component";
import { EditUserComponent } from "./edit-user/edit-user.component";
import { UserInfoComponent } from "./user-info/user-info.component";
import { PageListComponent } from "./page-list/page-list.component";
import { CreatePageComponent } from "./create-page/create-page.component";
import { EditPageComponent } from "./edit-page/edit-page.component";
import { PageInfoComponent } from "./page-info/page-info.component";
import { BlogListComponent } from "./blog-list/blog-list.component";
import { CreateBlogComponent } from "./create-blog/create-blog.component";
import { SettingsListComponent } from "./settings-list/settings-list.component";
import { EditSettingComponent } from "./edit-setting/edit-setting.component";

const routes: Routes = [
    {
        path: "",
        component: DashboardComponent,
    },
    {
        path: "users",
        component: UserListComponent,
    },
    {
        path: "users/create",
        component: CreateUserComponent,
    },
    {
        path: "users/:id",
        component: EditUserComponent,
    },
    {
        path: "users/:id/info",
        component: UserInfoComponent,
    },
    {
        path: "pages",
        component: PageListComponent,
    },
    {
        path: "pages/create",
        component: CreatePageComponent,
    },
    {
        path: "pages/:id",
        component: EditPageComponent,
    },
    {
        path: "pages/:id/info",
        component: PageInfoComponent,
    },
    {
        path: "blogs",
        component: BlogListComponent,
    },
    {
        path: "blogs/create",
        component: CreateBlogComponent,
    },
    {
        path: "blogs/:id",
        component: EditPageComponent,
    },
    {
        path: "blogs/:id/info",
        component: PageInfoComponent,
    },
    {
        path: "settings",
        component: SettingsListComponent,
    },
    {
        path: "settings/:key",
        component: EditSettingComponent,
    },
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule],
})
export class AppRoutingModule {}
