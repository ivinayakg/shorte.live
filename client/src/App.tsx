import { Button } from "@/components/ui/button";
import { ThemeProvider } from "@/components/theme-provider";
import "./App.css";
import { ModeToggle } from "./components/mode-toggle";
import { Input } from "./components/ui/input";
import {
  createBrowserRouter,
  Outlet,
  RouterProvider,
  useNavigate,
  useSearchParams,
  Link,
} from "react-router-dom";
import { MainProvider, useMain } from "./components/main-provider";
import { useEffect } from "react";
import { setInLocalStorage } from "./utils/localstorage";
import { AvatarImage, Avatar, AvatarFallback } from "./components/ui/avatar";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Main />,
    children: [
      {
        path: "",
        element: <CreateShort />,
      },
      // {
      //   path: "/my-urls",
      // },
    ],
  },
  {
    path: "/auth/",
    element: <Auth />,
  },
]);

function Auth() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    const token = searchParams.get("token");
    if (token) {
      setInLocalStorage("userToken", token);
    }
    navigate("/");
  });

  return null;
}

function Header() {
  const loginWithGoogleUrl = import.meta.env.VITE_GOOGLE_SIGN_IN;
  const { userState } = useMain();
  const navigate = useNavigate();

  const logout = () => {
    setInLocalStorage("userToken", null);
    navigate(0);
  };

  return (
    <div className="min-w-full flex justify-between items-center py-5">
      <h1 className="text-3xl">URL Shortner</h1>
      <div className="flex justify-center items-center gap-4">
        {userState.login ? (
          <>
            <Avatar>
              <AvatarImage src={userState.picture} alt="@shadcn" />
              <AvatarFallback>U</AvatarFallback>
            </Avatar>
            <Button className="" onClick={logout}>
              Logout
            </Button>
          </>
        ) : (
          <Link to={loginWithGoogleUrl}>
            <Button className="gap-2">
              Login With Google <img src="/google.svg" alt="" />
            </Button>
          </Link>
        )}
        <ModeToggle />
      </div>
    </div>
  );
}

function CreateShort() {
  const { userState } = useMain();
  const loginWithGoogleUrl = import.meta.env.VITE_GOOGLE_SIGN_IN;

  return (
    <div className="w-full flex flex-col justify-center items-center gap-10 createshort">
      <h2 className="text-7xl font-bold">Generate shorten link</h2>
      <div className="flex w-full max-w-5xl items-center space-x-2">
        <Input className="" type="text" placeholder="URL" />
        <Input className="w-2/5" type="text" placeholder="shorten" />
        {userState.login ? (
          <Button className="" type="submit">
            Generate
          </Button>
        ) : (
          <>
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button>Generate</Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Login First?</AlertDialogTitle>
                  <AlertDialogDescription>
                    Our service is free, but only to verified users!
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction>
                    <Link to={loginWithGoogleUrl}>
                      <Button className="gap-2">
                        Login With Google <img src="/google.svg" alt="" />
                      </Button>
                    </Link>
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </>
        )}
      </div>
    </div>
  );
}

function Main() {
  return (
    <>
      <Header />
      <Outlet />
    </>
  );
}

function App() {
  return (
    <>
      <MainProvider>
        <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
          <RouterProvider router={router} />
        </ThemeProvider>
      </MainProvider>
    </>
  );
}

export default App;
