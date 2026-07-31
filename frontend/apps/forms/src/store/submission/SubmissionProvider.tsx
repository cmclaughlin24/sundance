import { createContext, useRef } from "react";
import type { ISubmissionValue } from "@/types/submission";
import {
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

  return (
    <SubmissionStoreContext value={storeRef.current}>
      {children}
    </SubmissionStoreContext>
  );
};
