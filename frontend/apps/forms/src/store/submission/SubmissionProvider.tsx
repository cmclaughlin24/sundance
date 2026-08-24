import { createContext, useEffect, useRef } from "react";
import type { ISubmissionValue } from "@/types/submission";
import {
  createFormValues,
  createSubmissionStore,
  type SubmissionStoreApi,
} from "./submissionStore";

export type SubmissionProviderProps = React.PropsWithChildren<{
  rawSubmission: ISubmissionValue[] | undefined;
}>;

export const SubmissionStoreContext = createContext<SubmissionStoreApi | null>(
  null,
);

export const SubmissionProvider: React.FC<SubmissionProviderProps> = ({
  children,
  rawSubmission,
}) => {
  const storeRef = useRef<SubmissionStoreApi | null>(null);

  if (!storeRef.current) {
    storeRef.current = createSubmissionStore(rawSubmission);
  }

  useEffect(() => {
    storeRef.current!.setState({
      values: createFormValues(rawSubmission),
      errors: {},
    });
  }, [rawSubmission]);

  return (
    <SubmissionStoreContext value={storeRef.current}>
      {children}
    </SubmissionStoreContext>
  );
};
