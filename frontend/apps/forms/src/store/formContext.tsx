import { createContext, type PropsWithChildren } from "react";

export const FormContext = createContext({});

export const FormContextProvider: React.FC<PropsWithChildren<{}>> = ({
  children,
}) => {
  return <FormContext value={{}}>{children}</FormContext>;
};
