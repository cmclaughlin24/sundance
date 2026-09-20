import { CONFIG } from "@/constants/config";
import { BaseHttpService, type DefaultRequestOptions } from "./baseHttpService";
import type { ITag, ITagVersion } from "@/types/tag";
import type {
  CreateTagRequest,
  CreateTagVersionRequest,
  UpdateTagRequest,
} from "./tagsService.type";

export class TagsService extends BaseHttpService {
  static readonly serviceKey = "TagsService";

  constructor() {
    super(CONFIG.formsUrl);
  }

  getTags(
    options: DefaultRequestOptions,
    params: { include: "versions" | "" } = { include: "" },
  ): Promise<ITag[]> {
    const search = new URLSearchParams();
    params.include && search.set("include", params.include);
    return this._get<ITag[]>(`/api/v1/tags`, options, search);
  }

  getTag(id: string, options: DefaultRequestOptions): Promise<ITag> {
    return this._get<ITag>(`/api/v1/tags/${id}`, options);
  }

  async createTag(
    tag: CreateTagRequest,
    options: DefaultRequestOptions,
  ): Promise<ITag> {
    const resp = await this._post<CreateTagRequest, ITag>(
      `/api/v1/tags`,
      tag,
      options,
    );
    return resp.data;
  }

  async updateTag(
    id: string,
    tag: UpdateTagRequest,
    options: DefaultRequestOptions,
  ): Promise<ITag> {
    const resp = await this._put<UpdateTagRequest, ITag>(
      `/api/v1/tags/${id}`,
      tag,
      options,
    );
    return resp.data;
  }

  async getTagVersions(
    tagId: string,
    options: DefaultRequestOptions,
  ): Promise<ITagVersion[]> {
    return this._get<ITagVersion[]>(`/api/v1/tags/${tagId}/versions`, options);
  }

  async getTagVersion(
    tagId: string,
    versionId: string,
    options: DefaultRequestOptions,
  ): Promise<ITagVersion> {
    return this._get<ITagVersion>(
      `/api/v1/tags/${tagId}/versions/${versionId}`,
      options,
    );
  }

  async createTagVersion(
    tagId: string,
    version: CreateTagVersionRequest,
    options: DefaultRequestOptions,
  ): Promise<ITagVersion> {
    const resp = await this._post<CreateTagVersionRequest, ITagVersion>(
      `/api/v1/tags/${tagId}/versions`,
      version,
      options,
    );
    return resp.data;
  }

  async publishTagVersion(
    tagId: string,
    versionId: string,
    options: DefaultRequestOptions,
  ): Promise<ITagVersion> {
    const resp = await this._post<{}, ITagVersion>(
      `/api/v1/tags/${tagId}/versions/${versionId}/publish`,
      {},
      options,
    );
    return resp.data;
  }

  async deprecateTagVersion(
    tagId: string,
    versionId: string,
    options: DefaultRequestOptions,
  ): Promise<ITagVersion> {
    const resp = await this._post<{}, ITagVersion>(
      `/api/v1/tags/${tagId}/versions/${versionId}/deprecate`,
      {},
      options,
    );
    return resp.data;
  }

  async retireTagVersion(
    tagId: string,
    versionId: string,
    options: DefaultRequestOptions,
  ): Promise<ITagVersion> {
    const resp = await this._post<{}, ITagVersion>(
      `/api/v1/tags/${tagId}/versions/${versionId}/retire`,
      {},
      options,
    );
    return resp.data;
  }
}
