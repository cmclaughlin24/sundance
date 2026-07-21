import { BaseHttpService } from "./baseHttpService";

export class SubmissionsService extends BaseHttpService {
  static readonly serviceKey = "SubmissionsService";

  constructor(baseURL: string) {
    super(baseURL);
  }
}
