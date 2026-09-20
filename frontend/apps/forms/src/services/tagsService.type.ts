import type { TagNodeType } from "@/types/tag";

export interface CreateTagRequest {
  displayName: string;
  keyPath: string;
  nodeType: TagNodeType;
  primitiveType?: string;
  isCollection: boolean;
}

export type UpdateTagRequest = Pick<CreateTagRequest, 'displayName'>

export interface CreateTagVersionRequest {}
