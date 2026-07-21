import { BaseHttpService } from "./baseHttpService";

export class FormsService extends BaseHttpService {
  static readonly serviceKey = "FormsService";

  constructor(baseURL: string) {
    super(baseURL);
  }
}
