import { createContext, useEffect, useRef } from "react";
import type { ISubmissionSection } from "@/types/submission";
import {
  createSubmissionStore,
  type SubmissionStoreApi,
} from "./submissionStore";
import { hydrateSectionState } from "./utils";

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
  const hydrateFrom = useRef(rawSubmission);

  if (!storeRef.current) {
    storeRef.current = createSubmissionStore(rawSubmission);
  }

  useEffect(() => {
    if (hydrateFrom.current === rawSubmission) {
      return;
    }

    hydrateFrom.current = rawSubmission;
    storeRef.current!.setState({
      sections: hydrateSectionState(rawSubmission ?? []),
    });
  }, [rawSubmission]);

  return (
    <SubmissionStoreContext value={storeRef.current}>
      {children}
    </SubmissionStoreContext>
  );
};
