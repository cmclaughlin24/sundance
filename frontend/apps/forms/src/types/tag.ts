import type { HasVersion } from "@/utils/version";

export const TAG_COLLECTION_SEGMENT = "[*]";

export type TagNodeType = "primitive" | "object";

export interface ITag {
  id: string;
  tenantId: string;
  keyPath: string;
  displayName: string;
  nodeType: TagNodeType;
  createdAt: Date;
  updatedAt: Date;
  versions?: ITagVersion[];
}

export type TagVersionStatus = "draft" | "active" | "deprecated" | "retired";

export interface ITagVersion extends HasVersion {
  id: string;
  tagId: string;
  status: TagVersionStatus;
  createdAt: Date;
  deprecatedAt: Date;
  publishedAt: Date;
  RetiredAt: Date;
}
