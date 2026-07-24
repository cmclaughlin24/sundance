import { useAsyncData } from "./useAsyncData";
import { useTenantId } from "@/store/useFormContext";
import { useEvalContext } from "@/store/evalContext";
import type { IBindingSource, IDataSourceRef, ILookup } from "@/types/data";
import { useDataSourcesService } from "./useHttpService";
import { useMemo } from "react";
import type { EvalContext } from "@/utils/evaluate";

export function useDataSource(dataSourceRef: IDataSourceRef) {
  const tenantId = useTenantId();
  const evalCtx = useEvalContext();
  const dataSourcesService = useDataSourcesService();
  const accessToken = "placeholder";

  const filters = useMemo(
    () => resolveBindings(dataSourceRef.bindings, evalCtx),
    [dataSourceRef, evalCtx],
  );
  // NOTE: To protect against unnecessary API calls, the filters need to be serialized. Otherwise,
  // the filters reference value will change each time the evalCtx changes. (we only need to re-query
  // the data source if it's specific params have change)
  const serialized = serializeFilters(filters);

  const result = useAsyncData<ILookup[]>(async () => {
    if (!tenantId) {
      return [];
    }

    return await dataSourcesService.getLookups(
      dataSourceRef.dataSourceId,
      filters,
      {
        tenantId,
        token: accessToken,
      },
    );
  }, [tenantId, serialized]);

  return result;
}

function serializeFilters(filters: Record<string, any>): string {
  return Object.keys(filters)
    .sort()
    .map((key) => `${key}=${JSON.stringify(filters[key])}`)
    .join("&");
}

function resolveBindings(
  bindings: Record<string, IBindingSource>,
  evalCtx: EvalContext,
): Record<string, any> {
  const filters: Record<string, any> = {};

  for (const [key, binding] of Object.entries(bindings)) {
    if (binding.type === "static") {
      filters[key] = binding.value;
    } else if (binding.type === "field") {
      filters[key] = evalCtx[binding.key];
    }
  }

  return filters;
}
