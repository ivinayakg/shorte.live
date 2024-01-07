import auth from "@/utils/auth";
import { getFromLocalStorage } from "@/utils/localstorage";
import { createContext, useContext, useEffect, useState } from "react";

export type UserState = {
  email: string;
  name: string;
  picture: string;
  _id: string;
  created_at: string;
  token: string;
  login: boolean;
};

type MainProviderProps = {
  children: React.ReactNode;
};

type MainProviderContextType = {
  userState: UserState;
  setUserState: React.Dispatch<React.SetStateAction<UserState>>;
};

const initialUserState: UserState = {
  email: "",
  name: "",
  picture: "",
  _id: "",
  created_at: "",
  token: "",
  login: false,
};

const initialState: MainProviderContextType = {
  userState: initialUserState,
  setUserState: () => {},
};

const MainProviderContext = createContext(initialState);

export function MainProvider({ children, ...props }: MainProviderProps) {
  const [userState, setUserState] = useState(initialUserState);

  useEffect(() => {
    (async () => {
      const token = getFromLocalStorage("userToken");
      if (token) {
        const userData = await auth(token);
        setUserState({ ...userData, login: true });
      }
    })();
  }, []);

  const value = {
    userState,
    setUserState,
  };

  return (
    <MainProviderContext.Provider {...props} value={value}>
      {children}
    </MainProviderContext.Provider>
  );
}

export const useMain = () => {
  const context = useContext(MainProviderContext);

  if (context === undefined)
    throw new Error("useMain must be used within a MainProvider");

  return context;
};
