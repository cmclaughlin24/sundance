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
}

export type TagVersionStatus = "draft" | "active" | "deprecated" | "retired";

export interface ITagVersion {
  id: string;
  tagId: string;
  version: number;
  createdAt: Date;
  deprecatedAt: Date;
  publishedAt: Date;
  RetiredAt: Date;
}
