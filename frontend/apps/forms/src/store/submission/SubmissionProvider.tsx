import { createContext, useEffect, useRef } from "react";
import type { ISubmissionSection } from "@/types/submission";
import {
  hydrateSubmissionValues,
  createSubmissionStore,
  type SubmissionStoreApi,
} from "./submissionStore";

export type SubmissionProviderProps = React.PropsWithChildren<{
  rawSubmission: ISubmissionSection[] | undefined;
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
      sections: hydrateSubmissionValues(rawSubmission),
    });
  }, [rawSubmission]);

  return (
    <SubmissionStoreContext value={storeRef.current}>
      {children}
    </SubmissionStoreContext>
  );
};
