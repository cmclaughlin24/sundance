import type { IUserLookup } from "@/types/userLookup";
import { BaseHttpService, type DefaultRequestOptions } from "./baseHttpService";
import { LegacyCache } from "@/decorators/cached.decorator";

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

  /**
   * Gets a list of active accounts for a user.
   * @param userId - The unique identifier for a user.
   * @param returns A promise that resolves to a list of active accounts for the user.
   */
  @LegacyCache()
  async getUserAccounts(
    _userId: string,
    _options: DefaultRequestOptions,
  ): Promise<IUserLookup[]> {
    return Promise.resolve<IUserLookup[]>([]);
  }
}
