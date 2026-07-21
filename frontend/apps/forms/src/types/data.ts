export interface ILookup {
  label: string;
  value: any;
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
