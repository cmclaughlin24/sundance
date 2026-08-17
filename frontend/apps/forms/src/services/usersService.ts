import type { IUserLookup } from "@/types/userLookup";
import { BaseHttpService, type DefaultRequestOptions } from "./baseHttpService";

export class UsersService extends BaseHttpService {
  static readonly serviceKey = "UsersService";

  /**
   * Gets a list of user lookups.
   * @param searchTerm
   * @param returns A promise that resolves to a list of users.
   */
  async getUserLookups(
    _searchTerm: string,
    _options: DefaultRequestOptions,
  ): Promise<IUserLookup[]> {
    return Promise.resolve<IUserLookup[]>([]);
  }
}
