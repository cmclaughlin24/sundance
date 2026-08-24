export interface AppConfig {
  formsUrl: string;
  tenantsUrl: string;
  usersUrl: string;
}

export const CONFIG: Readonly<AppConfig> = {
  formsUrl: "/forms-api",
  tenantsUrl: "/tenants-api",
  usersUrl: "/users-api",
};
