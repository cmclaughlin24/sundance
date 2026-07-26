export type LookupValue = string | number;

export interface ILookup {
  label: string;
  value: LookupValue;
}

export interface IBindingSource {
  type: "field" | "static";
  key: string;
  value: any;
}

export interface IDataSourceRef {
  dataSourceId: string;
  bindings: Record<string, IBindingSource>;
}

export interface HasDataSourceRef {
  dataSourceRef?: IDataSourceRef;
}
