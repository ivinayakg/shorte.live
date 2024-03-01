import { createContext, useContext, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useToast } from "@/components/ui/use-toast";
import fetch from "@/utils/axios";

export type UserState = {
  email: string;
  name: string;
  picture: string;
  _id: string;
  created_at: string;
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
  login: false,
};

const initialState: MainProviderContextType = {
  userState: initialUserState,
  setUserState: () => {},
};

function resetState(stateSetter: any) {
  stateSetter({ ...initialState });
}

const MainProviderContext = createContext(initialState);

export function MainProvider({ children, ...props }: MainProviderProps) {
  const [userState, setUserState] = useState(initialUserState);
  const navigate = useNavigate();
  const { toast } = useToast();

  useEffect(() => {
    (async () => {
      try {
        const response = await fetch.get("/user/self", {
          withCredentials: true,
        });
        if (response?.status !== 200) {
          resetState(setUserState);
          navigate("/");
        } else {
          setUserState({ ...response.data, login: true });
          toast({
            title: "Logged in Successfully",
            duration: 2000,
          });
        }
      } catch (e: any) {
        if (typeof e === "string") toast({ title: e, duration: 2000 });
        if (e?.name === "AxiosError") {
          if (e?.response?.status === 503) {
            navigate("/maintenance");
            toast({
              title: "Server is down, try again later",
              duration: 2000,
            });
          }
        }
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
